package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/uscott/go-ftx/models"
)

func (s *Stream) wsSignature(conn *websocket.Conn, subaccount ...string) error {

	ms := time.Now().UTC().UnixNano() / int64(time.Millisecond)

	mac := hmac.New(sha256.New, []byte(s.client.secret))
	mac.Write([]byte(fmt.Sprintf("%dwebsocket_login", ms)))
	args := map[string]interface{}{
		"key":  s.client.apiKey,
		"sign": hex.EncodeToString(mac.Sum(nil)),
		"time": ms,
	}
	if len(subaccount) > 0 {
		args["subaccount"] = subaccount[0]
	}
	err := conn.WriteJSON(&models.WSRequestPrivate{
		Op:   "login",
		Args: args,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
