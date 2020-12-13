package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/uscott/goftx"
)

func TestClient_GetServerTime(t *testing.T) {
	ftx := goftx.New()

	serverTime, err := ftx.GetServerTime()
	require.NoError(t, err)
	fmt.Println(serverTime.Sub(time.Now().UTC()))
}
