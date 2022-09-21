package api

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/sanjujosh/go-ftx/models"
)

const (
	apiGetFills = "/fills"
)

type Fills struct {
	client *Client
}

func (f *Fills) GetFills(params *models.FillParams) ([]*models.Fill, error) {

	url := FormURL(apiGetFills)
	response, err := f.client.Get(params, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.Fill
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}
