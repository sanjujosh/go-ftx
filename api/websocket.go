package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
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
	dialer                 *websocket.Dialer
	wsReconnectionCount    int
	wsReconnectionInterval time.Duration
	isDebugMode            bool
	isLoggedIn             bool
	conn                   *websocket.Conn
	tickersC               chan *models.TickerResponse
	marketsC               chan *models.Market
	tradesC                chan *models.TradeResponse
	booksC                 chan *models.OrderBookResponse
	fillsC                 chan *models.FillResponse
	ordersC                chan *models.OrdersResponse
}

func (s *Stream) SetReconnectionCount(count int) {
	s.mu.Lock()
	s.wsReconnectionCount = count
	s.mu.Unlock()
}

func (s *Stream) SetDebugMode(isDebugMode bool) {
	s.mu.Lock()
	s.isDebugMode = isDebugMode
	s.mu.Unlock()
}

func (s *Stream) SetReconnectionInterval(interval time.Duration) {
	s.mu.Lock()
	s.wsReconnectionInterval = interval
	s.mu.Unlock()
}

func (s *Stream) printf(format string, v ...interface{}) {
	if s.isDebugMode {
		log.Printf(format+"\n", v)
	}
}

func (s *Stream) connect(requests ...models.WSRequest) (err error) {

	if s.conn == nil {
		s.conn, _, err = s.dialer.Dial(s.url, nil)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	s.printf("connected to %v", s.url)

	if err = s.subscribe(requests); err != nil {
		return errors.WithStack(err)
	}
	lastPong := time.Now()
	s.conn.SetPongHandler(
		func(msg string) error {
			lastPong = time.Now()
			if time.Since(lastPong) > websocketTimeout {
				// TODO handle this case
				errmsg := "PONG response time has been exceeded"
				s.printf(errmsg)
				return fmt.Errorf(errmsg) // Handled?
			} else {
				s.printf("PONG")
			}
			return nil
		})
	return nil
}

func (s *Stream) getEventResponse(
	ctx context.Context,
	eventsC chan interface{},
	msg *models.WsResponse,
	requests ...models.WSRequest) (err error) {

	err = s.conn.ReadJSON(&msg)

	if err != nil {

		s.printf("read msg: %v", err)

		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			return
		}

		err = s.reconnect(ctx, requests)
		if err != nil {
			s.printf("reconnect: %+v", err)
			return
		}

		return nil
	}

	if msg.Type == models.Subscribed || msg.Type == models.UnSubscribed {
		return
	}

	var response interface{}

	switch msg.Channel {
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

	eventsC <- response

	return
}

func (s *Stream) serve(
	ctx context.Context, requests ...models.WSRequest) (chan interface{}, error) {

	for _, req := range requests {
		if req.Channel == models.FillsChannel || req.Channel == models.OrdersChannel {
			if err := s.Authorize(); err != nil {
				return nil, errors.WithStack(err)
			}
			break
		}
	}

	err := s.connect(requests...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	eventsC := make(chan interface{})
	msg := models.WsResponse{}

	go func() {

		go func() {
			for {
				s.client.mu.Lock()
				if err = s.getEventResponse(ctx, eventsC, &msg, requests...); err != nil {
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
					s.printf("write close msg: %v", err)
					s.client.mu.Unlock()
					return
				}
				s.client.mu.Unlock()

				time.Sleep(time.Second)

				return

			case <-time.After(pingPeriod):

				s.printf("PING")

				s.client.mu.Lock()
				err = s.conn.WriteControl(
					websocket.PingMessage,
					[]byte(`{"op": "pong"}`),
					time.Now().UTC().Add(10*time.Second))

				if err != nil && err != websocket.ErrCloseSent {
					s.printf("write ping: %v", err)
				}
				s.client.mu.Unlock()

			}
		}
	}()

	return eventsC, err
}

func (s *Stream) reconnect(
	ctx context.Context, requests []models.WSRequest) (err error) {

	for i := 0; i < s.wsReconnectionCount; i++ {
		if err = s.connect(requests...); err == nil {
			return nil
		}
		select {
		case <-time.After(s.wsReconnectionInterval):
			if err = s.connect(requests...); err != nil {
				continue
			}
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return errors.New("Reconnection failed")
}

func (s *Stream) subscribe(requests []models.WSRequest) (err error) {
	for _, req := range requests {
		if err = s.conn.WriteJSON(req); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (s *Stream) Authorize() (err error) {

	if s.isLoggedIn {
		return
	}

	if s.conn == nil {
		s.conn, _, err = s.dialer.Dial(s.url, nil)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	ms := time.Now().UTC().UnixNano() / int64(time.Millisecond)
	mac := hmac.New(sha256.New, []byte(s.client.secret))

	_, err = mac.Write([]byte(fmt.Sprintf("%dwebsocket_login", ms)))
	if err != nil {
		return errors.WithStack(err)
	}

	args := map[string]interface{}{
		"key":  s.client.apiKey,
		"sign": hex.EncodeToString(mac.Sum(nil)),
		"time": ms,
	}

	if s.client.SubAccount != nil {
		args["subaccount"] = *s.client.SubAccount
	}

	err = s.conn.WriteJSON(&models.WSRequestAuthorize{
		Op:   "login",
		Args: args,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	s.isLoggedIn = true

	return
}

func (s *Stream) SubscribeToTickers(
	ctx context.Context, symbols ...string) (chan *models.TickerResponse, error) {

	if len(symbols) == 0 {
		return nil, errors.New("symbols missing")
	}

	requests := make([]models.WSRequest, len(symbols))
	for i, symb := range symbols {
		requests[i] = models.WSRequest{
			Channel: models.TickerChannel,
			Market:  symb,
			Op:      models.Subscribe,
		}
	}

	eventsC, err := s.serve(ctx, requests...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	s.tickersC = make(chan *models.TickerResponse)

	go func() {

		for {

			select {
			case <-ctx.Done():
				return
			case event := <-eventsC:
				ticker, ok := event.(*models.TickerResponse)
				if !ok {
					return
				}
				s.tickersC <- ticker
			}
		}
	}()

	return s.tickersC, nil
}

func (s *Stream) SubscribeToMarkets(
	ctx context.Context) (chan *models.Market, error) {

	eventsC, err := s.serve(ctx, models.WSRequest{
		Channel: models.MarketsChannel,
		Op:      models.Subscribe,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	s.marketsC = make(chan *models.Market)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventsC:
				data, ok := event.(json.RawMessage)
				if !ok {
					return
				}
				var markets struct {
					Data map[string]*models.Market `json:"data"`
				}

				if err = json.Unmarshal(data, &markets); err != nil {
					s.printf("unmarshal markets: %+v", err)
					return
				}

				for _, market := range markets.Data {
					s.marketsC <- market
				}
			}
		}
	}()

	return s.marketsC, nil
}

func (s *Stream) SubscribeToTrades(
	ctx context.Context, symbols ...string) (chan *models.TradeResponse, error) {

	if len(symbols) == 0 {
		return nil, errors.New("symbols missing")
	}

	requests := make([]models.WSRequest, len(symbols))
	for i, symb := range symbols {
		requests[i] = models.WSRequest{
			Channel: models.TradesChannel,
			Market:  symb,
			Op:      models.Subscribe,
		}
	}

	eventsC, err := s.serve(ctx, requests...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	s.tradesC = make(chan *models.TradeResponse)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventsC:
				trades, ok := event.(*models.TradesResponse)
				if !ok {
					return
				}
				for _, trade := range trades.Trades {
					s.tradesC <- &models.TradeResponse{
						Trade:        trade,
						BaseResponse: trades.BaseResponse,
					}
				}
			}
		}
	}()

	return s.tradesC, nil
}

func (s *Stream) SubscribeToOrderBooks(
	ctx context.Context, symbols ...string,
) (chan *models.OrderBookResponse, error) {

	if len(symbols) == 0 {
		return nil, errors.New("symbols is missing")
	}

	requests := make([]models.WSRequest, len(symbols))
	for i, symb := range symbols {
		requests[i] = models.WSRequest{
			Channel: models.OrderBookChannel,
			Market:  symb,
			Op:      models.Subscribe,
		}
	}

	eventsC, err := s.serve(ctx, requests...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	s.booksC = make(chan *models.OrderBookResponse)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventsC:
				book, ok := event.(*models.OrderBookResponse)
				if !ok {
					return
				}
				s.booksC <- book
			}
		}
	}()

	return s.booksC, nil
}

// TODO: Get fill and order streams to actually work right

func (s *Stream) SubscribeToFills(ctx context.Context) (chan *models.FillResponse, error) {

	eventsC, err := s.serve(ctx, models.WSRequest{
		Channel: models.FillsChannel,
		Op:      models.Subscribe,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	s.fillsC = make(chan *models.FillResponse)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventsC:
				fill, ok := event.(*models.FillResponse)
				if !ok {
					return
				}
				s.fillsC <- fill
			}
		}
	}()

	return s.fillsC, nil
}

func (s *Stream) SubscribeToOrders(ctx context.Context) (chan *models.OrdersResponse, error) {

	eventsC, err := s.serve(ctx, models.WSRequest{
		Channel: models.OrdersChannel,
		Op:      models.Subscribe,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	s.ordersC = make(chan *models.OrdersResponse)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventsC:
				order, ok := event.(*models.OrdersResponse)
				if !ok {
					return
				}
				s.ordersC <- order
			}
		}
	}()

	return s.ordersC, nil
}
