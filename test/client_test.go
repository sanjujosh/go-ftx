package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	ftx "github.com/uscott/go-ftx/api"
)

func TestClient_GetServerTime(t *testing.T) {
	ftx := ftx.New()

	serverTime, err := ftx.GetServerTime()
	require.NoError(t, err)
	fmt.Println(serverTime.Sub(time.Now().UTC()))
}
