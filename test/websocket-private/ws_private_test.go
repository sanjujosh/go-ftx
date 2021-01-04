package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/uscott/go-ftx/api"
)

const sleepDuration time.Duration = 15 * time.Second

func prepForTest() (*api.Client, *context.Context, chan struct{}) {
	ftx := api.New()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		time.Sleep(sleepDuration)
		cancel()
		done <- struct{}{}
	}()
	return ftx, &ctx, done
}

func TestStream_SubscribeToFills(t *testing.T) {

	ftx, ctx, done := prepForTest()

	data, err := ftx.Stream.SubscribeToFills(*ctx)
	require.NoError(t, err)

	for {
		select {
		case <-done:
			return
		case msg := <-data:
			t.Logf("Data: %+v\n", *msg)
		default:
		}
	}
}
