package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uscott/go-ftx/api"
)

func TestFutures_GetFutures(t *testing.T) {

	ftx := api.New()

	futures, err := ftx.Futures.GetFutures()
	assert.NoError(t, err)
	assert.NotNil(t, futures)
	for _, p := range futures {
		fmt.Printf("Description: %s\n", p.Description)
		fmt.Printf("Expiration:  %+v\n", p.Expiry.Format(time.RFC3339))
		fmt.Printf("Name:        %s\n", p.Name)
	}
}
