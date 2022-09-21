package testclient

import (
	"fmt"
	"testing"
	"time"

	"github.com/sanjujosh/go-ftx/api"
)

func TestClient_GetServerTime(t *testing.T) {

	ftx := api.New()
	serverTime, err := ftx.GetServerTime()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(serverTime.Sub(time.Now().UTC()))
}
