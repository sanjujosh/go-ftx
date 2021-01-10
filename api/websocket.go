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
	wsUrl             = "wss://ftx.com/ws/"
	websocketTimeout  = time.Second * 60
	pingPeriod        = (websocketTimeout * 9) / 10
	reconnectCount    = int(10)
	reconnectInterval = time.Second
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
	if !s.isDebugMode {
		return
	}
	log.Printf(format+"\n", v)
}

func (s *Stream) connect(requests ...models.WSRequest) (*websocket.Conn, error) {

	conn, _, err := s.dialer.Dial(s.url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	s.printf("connected to %v", s.url)

	if err = s.subscribe(conn, requests); err != nil {
		return nil, errors.WithStack(err)
	}
	lastPong := time.Now()
	conn.SetPongHandler(
		func(msg string) error {
			lastPong = time.Now()
			if time.Now().Sub(lastPong) > websocketTimeout {
				// TODO handle this case
				errmsg := "PONG response time has been exceeded"
				s.printf(errmsg)
				return fmt.Errorf(errmsg) // Handled?
			} else {
				s.printf("PONG")
			}
			return nil
		})
	return conn, nil
}

func (s *Stream) serve(
	ctx context.Context, requests ...models.WSRequest) (chan interface{}, error) {

	conn, err := s.connect(requests...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, req := range requests {
		switch req.Channel {
		case models.FillsChannel, models.OrdersChannel:
			if err = s.Authorize(conn); err != nil {
				return nil, errors.WithStack(err)
			}
			break
		}
	}
	doneC := make(chan struct{})
	eventsC := make(chan interface{}, 1)

	go func() {
		go func() {

			defer close(doneC)

			for {
				message := &models.WsResponse{}
				err = conn.ReadJSON(&message)
				if err != nil {
					s.printf("read msg: %v", err)
					if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
						return
					}
					conn, err = s.reconnect(ctx, requests)
					if err != nil {
						s.printf("reconnect: %+v", err)
						return
					}
					continue
				}

				switch message.Type {
				case models.Subscribed, models.UnSubscribed:
					continue
				}

				var response interface{}
				switch message.Channel {
				case models.TickerChannel:
					response, err = message.MapToTickerResponse()
				case models.TradesChannel:
					response, err = message.MapToTradesResponse()
				case models.OrderBookChannel:
					response, err = message.MapToOrderBookResponse()
				case models.MarketsChannel:
					response = message.Data
				case models.FillsChannel:
					response, err = message.MapToFillResponse()
				case models.OrdersChannel:
					response, err = message.MapToOrdersResponse()
				}

				eventsC <- response
			}
		}()

		for {
			select {
			case <-ctx.Done():
				err := conn.WriteMessage(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					s.printf("write close msg: %v", err)
					return
				}
				select {
				case <-doneC:
					return
				case <-time.After(time.Second):
					return
				}
			case <-doneC:
				return
			case <-time.After(pingPeriod):
				s.printf("PING")
				err := conn.WriteControl(
					websocket.PingMessage,
					[]byte(`{"op": "pong"}`),
					time.Now().Add(10*time.Second))
				if err != nil && err != websocket.ErrCloseSent {
					s.printf("write ping: %v", err)
				}
			}
		}
	}()

	return eventsC, nil
}

func (s *Stream) reconnect(
	ctx context.Context, requests []models.WSRequest) (*websocket.Conn, error) {

	for i := 1; i < s.wsReconnectionCount; i++ {
		conn, err := s.connect(requests...)
		if err == nil {
			return conn, nil
		}

		select {
		case <-time.After(s.wsReconnectionInterval):
			conn, err := s.connect(requests...)
			if err != nil {
				continue
			}

			return conn, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return nil, errors.New("reconnection failed")
}

func (s *Stream) subscribe(conn *websocket.Conn, requests []models.WSRequest) error {
	for _, req := range requests {
		err := conn.WriteJSON(req)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (s *Stream) Authorize(conn *websocket.Conn, subaccounts ...string) (err error) {

	if s.isLoggedIn {
		return
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
	if len(subaccounts) > 0 {
		args["subaccount"] = subaccounts[0]
	}
	err = conn.WriteJSON(&models.WSRequestAuthorize{
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

	requests := make([]models.WSRequest, 0, len(symbols))
	for _, symbol := range symbols {
		requests = append(requests, models.WSRequest{
			Channel: models.TickerChannel,
			Market:  symbol,
			Op:      models.Subscribe,
		})
	}

	eventsC, err := s.serve(ctx, requests...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tickersC := make(chan *models.TickerResponse, 1)
	go func() {
		defer close(tickersC)
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-eventsC:
				if !ok {
					return
				}
				ticker, ok := event.(*models.TickerResponse)
				if !ok {
					return
				}
				tickersC <- ticker
			}
		}
	}()

	return tickersC, nil
}

func (s *Stream) SubscribeToMarkets(ctx context.Context) (chan *models.Market, error) {

	eventsC, err := s.serve(ctx, models.WSRequest{
		Channel: models.MarketsChannel,
		Op:      models.Subscribe,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	marketsC := make(chan *models.Market, 1)
	go func() {
		defer close(marketsC)
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-eventsC:
				if !ok {
					return
				}
				data, ok := event.(json.RawMessage)
				if !ok {
					return
				}
				var markets struct {
					Data map[string]*models.Market `json:"data"`
				}
				err = json.Unmarshal(data, &markets)
				if err != nil {
					s.printf("unmarshal markets: %+v", err)
					return
				}
				for _, market := range markets.Data {
					marketsC <- market
				}
			}
		}
	}()

	return marketsC, nil
}

func (s *Stream) SubscribeToTrades(
	ctx context.Context, symbols ...string,
) (chan *models.TradeResponse, error) {

	if len(symbols) == 0 {
		return nil, errors.New("symbols missing")
	}

	requests := make([]models.WSRequest, 0, len(symbols))
	for _, symbol := range symbols {
		requests = append(requests, models.WSRequest{
			Channel: models.TradesChannel,
			Market:  symbol,
			Op:      models.Subscribe,
		})
	}

	eventsC, err := s.serve(ctx, requests...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tradesC := make(chan *models.TradeResponse, 1)
	go func() {
		defer close(tradesC)
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-eventsC:
				if !ok {
					return
				}
				trades, ok := event.(*models.TradesResponse)
				if !ok {
					return
				}
				for _, trade := range trades.Trades {
					tradesC <- &models.TradeResponse{
						Trade:        trade,
						BaseResponse: trades.BaseResponse,
					}
				}
			}
		}
	}()

	return tradesC, nil
}

func (s *Stream) SubscribeToOrderBooks(
	ctx context.Context, symbols ...string,
) (chan *models.OrderBookResponse, error) {

	if len(symbols) == 0 {
		return nil, errors.New("symbols is missing")
	}

	requests := make([]models.WSRequest, 0, len(symbols))
	for _, symbol := range symbols {
		requests = append(requests, models.WSRequest{
			Channel: models.OrderBookChannel,
			Market:  symbol,
			Op:      models.Subscribe,
		})
	}

	eventsC, err := s.serve(ctx, requests...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	booksC := make(chan *models.OrderBookResponse, 1)
	go func() {
		defer close(booksC)
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventsC:
				book, ok := event.(*models.OrderBookResponse)
				if !ok {
					return
				}
				booksC <- book
			}
		}
	}()

	return booksC, nil
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
	fillsC := make(chan *models.FillResponse, 1)

	go func() {
		defer close(fillsC)
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventsC:
				fill, ok := event.(*models.FillResponse)
				if !ok {
					return
				}
				fillsC <- fill
			}
		}
	}()

	return fillsC, nil
}

func (s *Stream) SubscribeToOrders(ctx context.Context) (chan *models.OrdersResponse, error) {

	eventsC, err := s.serve(ctx, models.WSRequest{
		Channel: models.OrdersChannel,
		Op:      models.Subscribe,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ordersC := make(chan *models.OrdersResponse, 1)

	go func() {
		defer close(ordersC)
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventsC:
				order, ok := event.(*models.OrdersResponse)
				if !ok {
					return
				}
				ordersC <- order
			}
		}
	}()

	return ordersC, nil
}
