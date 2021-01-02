package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/uscott/go-ftx/api"
)

func TestFutures_GetFutures(t *testing.T) {

	ftx := api.New()

	futures, err := ftx.Futures.GetFutures()
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range futures {
		if p == nil {
			t.Fatal("nil pointer")
		}
		fmt.Printf("Description: %s\n", p.Description)
		fmt.Printf("Expiration:  %+v\n", p.Expiry.Format(time.RFC3339))
		fmt.Printf("Name:        %s\n", p.Name)
	}
}
