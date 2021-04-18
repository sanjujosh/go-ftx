package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/uscott/go-ftx/models"
)

const (
	wsUrl                 = "wss://ftx.com/ws/"
	websocketTimeout      = time.Second * 60
	pingPeriod            = (websocketTimeout * 9) / 10
	reconnectCount    int = 10
	reconnectInterval     = time.Second
)

type Stream struct {
	client                 *Client
	mu                     *sync.Mutex
	url                    string
	conn                   *websocket.Conn
	dialer                 *websocket.Dialer
	wsReconnectionCount    int
	wsReconnectionInterval time.Duration
	isLoggedIn             bool
	WsSub                  *WsSub
	tickersC               chan *models.TickerResponse
	marketsC               chan *models.Market
	tradesC                chan *models.TradeResponse
	booksC                 chan *models.OrderBookResponse
	fillsC                 chan *models.FillResponse
	ordersC                chan *models.OrdersResponse
}

type TrivialMap map[string]struct{}

type WsSub struct {
	ChannelTypes map[models.ChannelType]TrivialMap
	Requests     []models.WSRequest
}

func NewStream(client *Client) *Stream {
	return &Stream{
		client:                 client,
		mu:                     &sync.Mutex{},
		url:                    wsUrl,
		dialer:                 websocket.DefaultDialer,
		wsReconnectionCount:    reconnectCount,
		wsReconnectionInterval: reconnectInterval,
		WsSub:                  NewWsSub(),
		tickersC:               make(chan *models.TickerResponse),
		marketsC:               make(chan *models.Market),
		tradesC:                make(chan *models.TradeResponse),
		booksC:                 make(chan *models.OrderBookResponse),
		fillsC:                 make(chan *models.FillResponse),
		ordersC:                make(chan *models.OrdersResponse),
	}
}

func NewWsSub() *WsSub {
	return &WsSub{
		ChannelTypes: make(map[models.ChannelType]TrivialMap),
		Requests:     make([]models.WSRequest, 0, 64),
	}
}

func MakeRequests(
	chantype models.ChannelType, symbols TrivialMap) []models.WSRequest {

	if len(symbols) == 0 {
		return []models.WSRequest{
			{ChannelType: chantype, Op: models.Subscribe},
		}
	}

	requests := make([]models.WSRequest, len(symbols))

	i := 0
	for s := range symbols {
		requests[i] = models.WSRequest{
			ChannelType: chantype,
			Market:      s,
			Op:          models.Subscribe,
		}
		i++
	}

	return requests
}

func (s *Stream) Authorize() (err error) {

	if s.conn == nil {
		if err = s.CreateNewConnection(); err != nil {
			return
		}
	}

	if s.isLoggedIn {
		return nil
	}

	wsra, err := s.GetAuthRequest()
	if err != nil {
		return
	}

	if err = s.conn.WriteJSON(wsra); err != nil {
		return errors.WithStack(err)
	}

	s.isLoggedIn = true

	return
}

func (s *Stream) Connect(requests ...models.WSRequest) (err error) {

	if err = s.CreateNewConnection(); err != nil {
		return
	}

	s.client.Logger.Debugf("connected to %v", s.url)

	if err = s.Subscribe(); err != nil {
		return errors.WithStack(err)
	}

	lastPong := time.Now()
	s.conn.SetPongHandler(
		func(msg string) error {
			lastPong = time.Now()
			if time.Since(lastPong) > websocketTimeout {
				// TODO handle this case
				errmsg := "PONG response time has been exceeded"
				s.client.Logger.Debug(errmsg)
				return errors.New(errmsg) // Handled?
			}
			s.client.Logger.Debug("PONG")
			return nil
		})
	return nil
}

func (s *Stream) CreateNewConnection() (err error) {

	s.isLoggedIn = false

	s.conn, _, err = s.dialer.Dial(s.url, nil)
	if err != nil {
		return errors.WithStack(err)
	}

	return
}

func (s *Stream) GetAuthRequest() (*models.WSRequestAuthorize, error) {

	ms := time.Now().UTC().UnixNano() / int64(time.Millisecond)
	mac := hmac.New(sha256.New, []byte(s.client.secret))

	_, err := mac.Write([]byte(fmt.Sprintf("%dwebsocket_login", ms)))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	args := map[string]interface{}{
		"key":  s.client.apiKey,
		"sign": hex.EncodeToString(mac.Sum(nil)),
		"time": ms,
	}

	if s.client.SubAccount != nil {
		args["subaccount"] = *s.client.SubAccount
	}

	return &models.WSRequestAuthorize{
		Op:   "login",
		Args: args,
	}, nil
}

func (s *Stream) GetEventResponse(ctx context.Context, msg *models.WsResponse) (err error) {

	if msg == nil {
		return errors.New("Nil pointer")
	}

	if err = s.conn.ReadJSON(&msg); err != nil {

		s.client.Logger.Debugf("read msg: %v", err)

		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			return
		}

		if err = s.Reconnect(ctx); err != nil {
			s.client.Logger.Debugf("reconnect: %+v", err)
			return
		}

		return nil
	}

	if msg.ResponseType == models.Subscribed || msg.ResponseType == models.UnSubscribed {
		return
	}

	var response interface{}

	switch msg.ChannelType {
	case models.TickerChannel:
		response, err = msg.MapToTickerResponse()
	case models.TradesChannel:
		response, err = msg.MapToTradesResponse()
	case models.OrderBookChannel:
		response, err = msg.MapToOrderBookResponse()
	case models.MarketsChannel:
		response = msg.Data
	case models.FillsChannel:
		response, err = msg.MapToFillResponse()
	case models.OrdersChannel:
		response, err = msg.MapToOrdersResponse()
	}

	if err != nil {
		return
	}

	go s.SendToChannel(msg.ChannelType, response)

	return
}

func (s *Stream) IsLoggedIn() bool {
	return s.isLoggedIn
}

func (s *Stream) Reconnect(ctx context.Context) (err error) {

	for i := 0; i < s.wsReconnectionCount; i++ {
		if err = s.Connect(); err == nil {
			return nil
		}
		select {
		case <-time.After(s.wsReconnectionInterval):
			if err = s.Connect(); err != nil {
				continue
			}
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return errors.New("Reconnection failed")
}

func (s *Stream) SetDebugMode(isDebugMode bool) {
	s.mu.Lock()
	s.isDebugMode = isDebugMode
	s.mu.Unlock()
}

func (s *Stream) SetReconnectionCount(count int) {
	s.mu.Lock()
	s.wsReconnectionCount = count
	s.mu.Unlock()
}

func (s *Stream) SetReconnectionInterval(interval time.Duration) {
	s.mu.Lock()
	s.wsReconnectionInterval = interval
	s.mu.Unlock()
}

func (s *Stream) sub() (err error) {
	for _, r := range s.WsSub.Requests {
		if err = s.conn.WriteJSON(r); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (s *Stream) Subscribe() (err error) {

	if !s.isLoggedIn {
		for _, r := range s.WsSub.Requests {
			ct := r.ChannelType
			if ct == models.FillsChannel || ct == models.OrdersChannel {
				if err = s.Authorize(); err != nil {
					return
				}
				break
			}
		}
	}

	return s.sub()
}

func (s *Stream) SendToChannel(ct models.ChannelType, response interface{}) {

	switch ct {
	case models.TickerChannel:
		ticker, ok := response.(*models.TickerResponse)
		if ok && ticker != nil {
			s.tickersC <- ticker
		}
	case models.TradesChannel:
		trades, ok := response.(*models.TradesResponse)
		if ok && trades != nil {
			for _, t := range trades.Trades {
				s.tradesC <- &models.TradeResponse{
					Trade:        t,
					BaseResponse: trades.BaseResponse,
				}
			}
		}
	case models.OrderBookChannel:
		book, ok := response.(*models.OrderBookResponse)
		if ok && book != nil {
			s.booksC <- book
		}
	case models.MarketsChannel:
		markets, err := MapToMarketData(response)
		if err == nil {
			for _, m := range markets {
				if m != nil {
					s.marketsC <- m
				}
			}
		}
	case models.FillsChannel:
		fill, ok := response.(*models.FillResponse)
		if ok && fill != nil {
			s.fillsC <- fill
		}
	case models.OrdersChannel:
		order, ok := response.(*models.OrdersResponse)
		if ok && order != nil {
			s.ordersC <- order
		}
	}
}

func (s *Stream) Serve(ctx context.Context) (err error) {

	if err = s.Connect(); err != nil {
		return errors.WithStack(err)
	}

	msg := models.WsResponse{}

	go func() {

		go func() {
			for {
				s.client.mu.Lock()
				if err = s.GetEventResponse(ctx, &msg); err != nil {
					s.client.mu.Unlock()
					return
				}
				s.client.mu.Unlock()
			}
		}()

		for {

			select {

			case <-ctx.Done():

				s.client.mu.Lock()
				err = s.conn.WriteMessage(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

				if err != nil {
					s.client.Logger.Debugf("write close msg: %v", err)
					s.client.mu.Unlock()
					return
				}
				s.client.mu.Unlock()

				time.Sleep(time.Second)

				return

			case <-time.After(pingPeriod):

				s.client.Logger.Debug("PING")

				s.client.mu.Lock()
				err = s.conn.WriteControl(
					websocket.PingMessage,
					[]byte(`{"op": "pong"}`),
					time.Now().UTC().Add(10*time.Second))

				if err != nil && err != websocket.ErrCloseSent {
					s.client.Logger.Debugf("write ping: %v", err)
				}
				s.client.mu.Unlock()

			}
		}
	}()

	return err
}

func (s *Stream) SubscribeToTickers(
	ctx context.Context, symbols ...string) (chan *models.TickerResponse, error) {

	if len(symbols) == 0 {
		return nil, errors.New("symbols missing")
	}

	ct, ws := models.TickerChannel, s.WsSub

	ws.AppendRequests(ct, symbols...)

	err := s.Serve(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return s.tickersC, nil
}

func (s *Stream) SubscribeToMarkets(ctx context.Context) (chan *models.Market, error) {

	ct, ws := models.MarketsChannel, s.WsSub

	ws.AppendRequests(ct)

	err := s.Serve(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return s.marketsC, nil
}

func (s *Stream) SubscribeToTrades(
	ctx context.Context, symbols ...string) (chan *models.TradeResponse, error) {

	if len(symbols) == 0 {
		return nil, errors.New("symbols missing")
	}

	ct, ws := models.TradesChannel, s.WsSub

	ws.AppendRequests(ct, symbols...)

	err := s.Serve(ctx)
	if err != nil {
		return nil, err
	}

	return s.tradesC, nil
}

func (s *Stream) SubscribeToOrderBooks(
	ctx context.Context, symbols ...string) (chan *models.OrderBookResponse, error) {

	if len(symbols) == 0 {
		return nil, errors.New("symbols is missing")
	}

	ct, ws := models.OrderBookChannel, s.WsSub

	ws.AppendRequests(ct, symbols...)

	err := s.Serve(ctx)
	if err != nil {
		return nil, err
	}

	return s.booksC, nil
}

// TODO: Get fill and order streams to actually work right

func (s *Stream) SubscribeToFills(ctx context.Context) (chan *models.FillResponse, error) {

	ct, ws := models.FillsChannel, s.WsSub

	ws.AppendRequests(ct)

	err := s.Serve(ctx)
	if err != nil {
		return nil, err
	}

	return s.fillsC, nil
}

func (s *Stream) SubscribeToOrders(
	ctx context.Context, symbols ...string) (chan *models.OrdersResponse, error) {

	if len(symbols) == 0 {
		return nil, errors.New("symbols missing")
	}

	ct, ws := models.OrdersChannel, s.WsSub

	ws.AppendRequests(ct, symbols...)

	err := s.Serve(ctx)
	if err != nil {
		return nil, err
	}

	return s.ordersC, nil
}

func (s *Stream) WSConn() *websocket.Conn {
	return s.conn
}

func MapToMarketData(event interface{}) (map[string]*models.Market, error) {

	data, ok := event.(json.RawMessage)
	if !ok {
		return nil, fmt.Errorf("Convert fail")
	}

	var markets struct {
		Data map[string]*models.Market `json:"data"`
	}

	if err := json.Unmarshal(data, &markets); err != nil {
		return nil, fmt.Errorf("Unmarshal markets: %+v", err)
	}

	return markets.Data, nil
}

func (ws *WsSub) AppendRequests(ct models.ChannelType, symbols ...string) {

	ctypes, tm := ws.ChannelTypes, make(TrivialMap)

	if ws.Requests == nil {
		ws.Requests = make([]models.WSRequest, 0, 64)
	}

	if ctypes[ct] == nil {

		for _, s := range symbols {
			tm[s] = struct{}{}
		}

		ctypes[ct] = tm
		ws.Requests = append(ws.Requests, MakeRequests(ct, tm)...)

		return
	}

	for _, s := range symbols {
		_, ok := ctypes[ct][s]
		if !ok {
			ctypes[ct][s] = struct{}{}
			tm[s] = struct{}{}
		}
	}

	if len(tm) > 0 {
		ws.Requests = append(ws.Requests, MakeRequests(ct, tm)...)
	}
}
