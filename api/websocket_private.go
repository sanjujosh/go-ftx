package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/uscott/go-ftx/models"
)

func (s *Stream) wsLogin(conn *websocket.Conn, subaccounts ...string) (err error) {

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
	err = conn.WriteJSON(&models.WSRequestPrivate{
		Op:   "login",
		Args: args,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *Stream) connectPrivate(requests ...models.WSRequestPrivate) (*websocket.Conn, error) {

	conn, _, err := s.dialer.Dial(s.url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	s.printf("connected to %v", s.url)
	if err = s.subscribePrivate(conn, requests); err != nil {
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

func (s *Stream) subscribePrivate(
	conn *websocket.Conn, requests []models.WSRequestPrivate) error {

	for _, req := range requests {
		err := conn.WriteJSON(req)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (s *Stream) servePrivate(
	ctx context.Context, requests ...models.WSRequestPrivate) (chan interface{}, error) {

	conn, err := s.connectPrivate(requests...)
	if err != nil {
		return nil, errors.WithStack(err)
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
					conn, err = s.reconnectPrivate(ctx, requests)
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

	return nil, nil
}

func (s *Stream) reconnectPrivate(
	ctx context.Context, requests []models.WSRequestPrivate) (*websocket.Conn, error) {

	for i := 1; i < s.wsReconnectionCount; i++ {
		conn, err := s.connectPrivate(requests...)
		if err == nil {
			return conn, nil
		}
		select {
		case <-time.After(s.wsReconnectionInterval):
			conn, err := s.connectPrivate(requests...)
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
