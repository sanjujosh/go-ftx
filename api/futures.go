package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/uscott/go-ftx/models"
)

const (
	apiGetFutures         = "/futures"
	apiGetFutureStats     = "/futures/%s/stats"
	apiGetFundingRates    = "/funding_rates"
	apiGetIndexWeights    = "/indexes/%s/weights"
	apiGetExpiredFutures  = "/expired_futures"
	apiGetHistoricalIndex = "/indexes"
)

type Futures struct {
	client *Client
}

func (f *Futures) GetFutures() ([]*models.Future, error) {

	request, err := f.client.prepareRequest(Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetFutures),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := f.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.Future
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (f *Futures) GetFutureByName(name string) (*models.Future, error) {

	request, err := f.client.prepareRequest(Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s/%s", apiUrl, apiGetFutures, name),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := f.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result models.Future
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &result, nil
}

func (f *Futures) GetFutureStats(name string) (*models.FutureStats, error) {

	path := fmt.Sprintf(apiGetFutureStats, name)
	request, err := f.client.prepareRequest(Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, path),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	response, err := f.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result models.FutureStats
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (f *Futures) GetFundingRates(
	params *models.FundingRatesParams) ([]*models.FundingRates, error) {

	queryParams, err := PrepareQueryParams(params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := f.client.prepareRequest(Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetFundingRates),
		Params: queryParams,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	response, err := f.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result []*models.FundingRates
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (f *Futures) GetIndexWeights(marketName string) (*map[string]float64, error) {

	path := fmt.Sprintf(apiGetIndexWeights, marketName)
	request, err := f.client.prepareRequest(Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, path),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	response, err := f.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result map[string]float64
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (f *Futures) GetExpiredFutures() ([]*models.FutureExpired, error) {

	request, err := f.client.prepareRequest(Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetExpiredFutures),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	response, err := f.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result []*models.FutureExpired
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (f *Futures) GetHistoricalIndex(
	params *models.HistoricalIndexParams) (*models.HistoricalIndex, error) {

	queryParams, err := PrepareQueryParams(params)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	request, err := f.client.prepareRequest(Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetHistoricalIndex),
		Params: queryParams,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	response, err := f.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result models.HistoricalIndex
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}
