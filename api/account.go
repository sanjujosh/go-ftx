package api

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/uscott/go-ftx/models"
	"github.com/uscott/go-tools/errs"
)

const (
	apiGetAccountInformation = "/account"
	apiGetPositions          = "/positions"
	apiPostLeverage          = "/account/leverage"
)

type Account struct {
	client *Client
}

func (a *Account) GetAccountInformation(result *models.AccountInformation) (err error) {

	if result == nil {
		return errs.NilPtr
	}
	url := FormURL(apiGetAccountInformation)
	response, err := a.client.Get(nil, url, true)
	if err != nil {
		return errors.WithStack(err)
	}

	err = json.Unmarshal(response, result)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *Account) GetPositions() ([]*models.Position, error) {

	url := FormURL(apiGetPositions)
	response, err := a.client.Get(nil, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.Position
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (a *Account) ChangeAccountLeverage(leverage float64) (result string, err error) {

	url := FormURL(apiPostLeverage)
	l := decimal.NewFromFloat(leverage)
	params := &struct {
		Leverage *decimal.Decimal `json:"leverage"`
	}{Leverage: &l}

	response, err := a.client.Post(params, url)
	if err != nil {
		return result, errors.WithStack(err)
	}

	if err = json.Unmarshal(response, &result); err != nil {
		return result, errors.WithStack(err)
	}
	return
}
