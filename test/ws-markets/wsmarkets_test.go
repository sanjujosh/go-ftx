package wsmarketstest

import (
	"context"
	"testing"
	"time"

	"github.com/uscott/go-ftx/api"
)

func Test_WsMarkets(t *testing.T) {

	client := api.New()
	client.Stream.SetDebugMode(true)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	marketsC, err := client.Stream.SubscribeToMarkets(ctx)
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		time.Sleep(5 * time.Second)
		done <- struct{}{}
	}()

	for {
		select {
		case <-done:
			t.Log("Exiting")
			return
		case market := <-marketsC:
			if market != nil {
				t.Logf("Market: %+v", *market)
			}
		}
	}
}
