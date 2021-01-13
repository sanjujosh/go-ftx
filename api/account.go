package api

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	request, err := a.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetPositions),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := a.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.Position
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (a *Account) ChangeAccountLeverage(leverage decimal.Decimal) error {
	body, err := json.Marshal(struct {
		Leverage decimal.Decimal `json:"leverage"`
	}{Leverage: leverage})
	if err != nil {
		return errors.WithStack(err)
	}

	request, err := a.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiPostLeverage),
		Body:   body,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = a.client.do(request)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
