package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-tools/tm"
)

func TestFutures_GetFutures(t *testing.T) {

	ftx := api.New()

	futures, err := ftx.Futures.GetFutures()
	assert.NoError(t, err)
	assert.NotNil(t, futures)
	for _, p := range futures {
		fmt.Printf("Description: %s\n", p.Description)
		fmt.Printf("Expiration:  %+v\n", tm.Format0(p.Expiry))
		fmt.Printf("Name:        %s\n", p.Name)
	}
}
