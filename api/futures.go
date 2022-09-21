package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sanjujosh/go-ftx/models"
	"github.com/uscott/go-tools/errs"
)

const (
	apiGetFutures         = "/futures"
	apiGetFutureStats     = "/futures/%s/stats"
	apiGetFundingRates    = "/funding_rates"
	apiGetIndexWeights    = "/indexes/%s/weights"
	apiGetExpiredFutures  = "/expired_futures"
	apiGetHistoricalIndex = "/indexes/%s/candles"
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

func (f *Futures) GetFutureByName(name string, future *models.Future) (err error) {

	if future == nil {
		return errs.NilPtr
	}
	url := FormURL(fmt.Sprintf("%s/%s", apiGetFutures, name))
	response, err := f.client.Get(nil, url, false)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = json.Unmarshal(response, future); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (f *Futures) GetFutureStats(future string, stats *models.FutureStats) (err error) {

	if stats == nil {
		panic(errs.NilPtrArg)
	}

	url := FormURL(fmt.Sprintf(apiGetFutureStats, future))

	response, err := f.client.Get(nil, url, false)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = json.Unmarshal(response, stats); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (f *Futures) GetFundingRates() ([]*models.FundingRates, error) {

	url := FormURL(apiGetFundingRates)

	response, err := f.client.Get(nil, url, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.FundingRates
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (f *Futures) GetIndexWeights(index string) (*map[string]float64, error) {

	url := FormURL(fmt.Sprintf(apiGetIndexWeights, index))

	response, err := f.client.Get(nil, url, false)
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

	url := FormURL(apiGetExpiredFutures)

	response, err := f.client.Get(nil, url, false)
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
	indexName string,
	params *models.HistoricalIndexParams) ([]*models.HistoricalIndex, error) {

	url := FormURL(fmt.Sprintf(apiGetHistoricalIndex, indexName))

	response, err := f.client.Get(params, url, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.HistoricalIndex
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}
