package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	"github.com/uscott/go-ftx/models"
	"github.com/uscott/go-tools/errs"
)

const (
	apiGetMarkets          = "/markets"
	apiGetOrderBook        = "/markets/%s/orderbook"
	apiGetTrades           = "/markets/%s/trades"
	apiGetHistoricalPrices = "/markets/%s/candles"
)

type Markets struct {
	client *Client
}

func (m *Markets) GetMarkets() ([]*models.Market, error) {

	url := FormURL(apiGetMarkets)
	response, err := m.client.Get(nil, url, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.Market

	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (m *Markets) GetMarketByName(name string, market *models.Market) (err error) {

	url := FormURL(fmt.Sprintf("%s/%s", apiGetMarkets, name))
	response, err := m.client.Get(nil, url, false)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = json.Unmarshal(response, market); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *Markets) GetOrderBook(market string, depth *int, ob *models.OrderBook) (err error) {

	if ob == nil {
		return errs.NilPtr
	}

	url := FormURL(fmt.Sprintf(apiGetOrderBook, market))
	var response []byte

	if depth == nil {
		response, err = m.client.Get(nil, url, false)
	} else {
		request, err := m.client.prepareRequest(Request{
			Auth:   false,
			Method: http.MethodGet,
			URL:    url,
			Params: map[string]string{"depth": strconv.FormatInt(int64(*depth), 10)},
		})
		if err != nil {
			return errors.WithStack(err)
		}
		response, err = m.client.do(request)
	}
	if err != nil {
		return errors.WithStack(err)
	}

	if err = json.Unmarshal(response, ob); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *Markets) GetTrades(
	market string, params *models.GetTradesParams) ([]*models.Trade, error) {

	url := FormURL(fmt.Sprintf(apiGetTrades, market))
	response, err := m.client.Get(params, url, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.Trade

	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (m *Markets) GetHistoricalPrices(
	market string, params *models.GetHistoricalPricesParams,
) ([]*models.HistoricalPrice, error) {

	url := FormURL(fmt.Sprintf(apiGetHistoricalPrices, market))
	response, err := m.client.Get(params, url, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.HistoricalPrice

	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}
