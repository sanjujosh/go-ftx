package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sanjujosh/go-ftx/models"
	"github.com/shopspring/decimal"
)

type Convert struct {
	client *Client
}

const (
	apiRequestQuote   = "/otc/quotes"
	apiGetQuoteStatus = "/otc/quotes/%d"
	apiAcceptQuote    = "/otc/quotes/%d/accept"
)

func (c *Convert) RequestQuote(from, to string, size decimal.Decimal) (id int64, err error) {

	params := struct {
		FromCoin *string          `json:"fromCoin"`
		ToCoin   *string          `json:"toCoin"`
		Size     *decimal.Decimal `json:"size"`
	}{FromCoin: &from, ToCoin: &to, Size: &size}

	url := fmt.Sprintf("%s%s", apiUrl, apiRequestQuote)

	response, err := c.client.Post(&params, url)
	if err != nil {
		return 0, err
	}

	var result = struct {
		QuoteID int64 `json:"quoteId"`
	}{}
	if err = json.Unmarshal(response, &result); err != nil {
		return 0, errors.WithStack(err)
	}
	return result.QuoteID, nil
}

func (c *Convert) GetQuoteStatus(id int64) (*models.ConvertQuoteStatus, error) {

	path := fmt.Sprintf(apiGetQuoteStatus, id)
	url := fmt.Sprintf("%s%s", apiUrl, path)
	response, err := c.client.Get(nil, url, true)
	if err != nil {
		return nil, err
	}

	var result models.ConvertQuoteStatus
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (c *Convert) AcceptQuote(id int64) error {

	path := fmt.Sprintf(apiAcceptQuote, id)

	request, err := c.client.prepareRequest(Request{
		Auth:       true,
		Method:     http.MethodPost,
		URL:        fmt.Sprintf("%s%s", apiUrl, path),
		SubAccount: c.client.SubAccount,
	})

	if err != nil {
		return errors.WithStack(err)
	}

	_, errReq := c.client.do(request)

	if errReq != nil {
		return errors.WithStack(errReq)
	}

	return nil
}
