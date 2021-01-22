package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/uscott/go-ftx/models"
)

const (
	apiListOptionQuoteRequests            = "/options/requests"
	apiListUserOptionQuoteRequests        = "/options/my_requests"
	apiCreateOptionQuoteRequest           = "/options/requests"
	apiCancelOptionQuoteRequest           = "/options/requests/%d"
	apiGetQuotesForUserOptionQuoteRequest = "/options/requests/%d/quotes"
	apiCreateOptionQuote                  = apiGetQuotesForUserOptionQuoteRequest
	apiUserOptionQuotes                   = "/options/my_quotes"
	apiCancelUserOptionQuote              = "/options/quotes/%d"
	apiAcceptOptionQuote                  = "/options/quotes/%d/accept"
	apiGetOptionsAccountInfo              = "/options/account_info"
	apiGetOptionsPositions                = "/options/positions"
	apiGetPublicOptionsTrades             = "/options/trades"
	apiGetOptionsFills                    = "/options/fills"
	apiGet24hOptionsVolume                = "/stats/24h_options_volume"
	apiGetOptionsHistoricalVolumes        = "/options/historical_volumes/BTC"
	apiGetOptionsOpenInterest             = "/options/open_interest/BTC"
	apiGetOptionsHistoricalOpenInterest   = "/options/historical_open_interest/BTC"
)

type Options struct {
	client *Client
}

func (o *Options) ListQuoteRequests() ([]*models.OptionQuoteRequest, error) {

	url := FormURL(apiListOptionQuoteRequests)
	response, err := o.client.Get(&struct{}{}, url, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.OptionQuoteRequest
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (o *Options) ListUserQuoteRequests() ([]*models.OptionQuoteRequest, error) {

	url := FormURL(apiListUserOptionQuoteRequests)
	response, err := o.client.Get(&struct{}{}, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.OptionQuoteRequest
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (o *Options) CreateQuoteRequest(
	params *models.OptionQuoteRequestParams,
) (*models.CreateQuoteRequest, error) {

	url := FormURL(apiCreateOptionQuoteRequest)
	response, err := o.client.Post(params, url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := models.CreateQuoteRequest{}
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (o *Options) CancelQuoteRequest(id int64) (*models.CancelQuoteRequest, error) {

	url := FormURL(fmt.Sprintf(apiCancelOptionQuoteRequest, id))
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodDelete,
		URL:    url,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := models.CancelQuoteRequest{}
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (o *Options) GetQuotesForUserQuoteRequest(
	id int64,
) ([]*models.QuotesForOptionQuoteRequest, error) {

	url := FormURL(fmt.Sprintf(apiGetQuotesForUserOptionQuoteRequest, id))
	response, err := o.client.Get(&struct{}{}, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.QuotesForOptionQuoteRequest
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (o *Options) CreateQuote(
	id int64, price decimal.Decimal,
) (*models.UserOptionQuote, error) {

	url := FormURL(fmt.Sprintf(apiCreateOptionQuote, id))

	body := &struct {
		Price *decimal.Decimal `json:"price"`
	}{Price: &price}

	response, err := o.client.Post(body, url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := models.UserOptionQuote{}
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (o *Options) GetUserQuotes() ([]*models.UserOptionQuote, error) {

	url := FormURL(apiUserOptionQuotes)
	response, err := o.client.Get(&struct{}{}, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.UserOptionQuote
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (o *Options) CancelQuote(id int64) (*models.UserOptionQuote, error) {

	url := FormURL(fmt.Sprintf(apiCancelUserOptionQuote, id))
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodDelete,
		URL:    url,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := models.UserOptionQuote{}
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (o *Options) AcceptQuote(id int64) (*models.UserOptionQuote, error) {

	url := FormURL(fmt.Sprintf(apiAcceptOptionQuote, id))
	response, err := o.client.Post(&struct{}{}, url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := models.UserOptionQuote{}
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (o *Options) GetAccountOptionsInfo() (*models.AccountOptionsInfo, error) {

	url := FormURL(apiGetOptionsAccountInfo)
	response, err := o.client.Get(&struct{}{}, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := models.AccountOptionsInfo{}
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (o *Options) GetOptionsPositions() ([]*models.OptionPosition, error) {

	url := FormURL(apiGetOptionsPositions)
	response, err := o.client.Get(&struct{}{}, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.OptionPosition
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (o *Options) GetPublicOptionsTrades(
	params *models.NumberTimeLimit,
) ([]*models.PublicOptionTrade, error) {

	url := FormURL(apiGetPublicOptionsTrades)
	response, err := o.client.Get(params, url, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.PublicOptionTrade
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (o *Options) GetOptionsFills(
	params *models.NumberTimeLimit,
) ([]*models.OptionFill, error) {

	url := FormURL(apiGetOptionsFills)
	response, err := o.client.Get(params, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.OptionFill
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (o *Options) Get24hOptionVolume() (*models.OptionsVolume, error) {

	url := FormURL(apiGet24hOptionsVolume)
	response, err := o.client.Get(&struct{}{}, url, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := models.OptionsVolume{}
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (o *Options) GetOptionsHistoricalVolumes(
	params *models.NumberTimeLimit,
) ([]*models.OptionsHistoricalVolumes, error) {

	url := FormURL(apiGetOptionsHistoricalVolumes)
	response, err := o.client.Get(params, url, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.OptionsHistoricalVolumes
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (o *Options) GetOptionsOpenInterest() (openInterest decimal.Decimal, err error) {

	url := FormURL(apiGetOptionsOpenInterest)
	response, err := o.client.Get(&struct{}{}, url, false)
	if err != nil {
		return decimal.Decimal{}, errors.WithStack(err)
	}

	result := struct {
		OpenInterest decimal.Decimal `json:"openInterest"`
	}{}
	if err = json.Unmarshal(response, &result); err != nil {
		return decimal.Decimal{}, errors.WithStack(err)
	}
	return result.OpenInterest, nil
}

func (o *Options) GetHistoricalOpenInterest(
	params *models.NumberTimeLimit,
) ([]*models.OptionsHistoricalOpenInterest, error) {

	url := FormURL(apiGetOptionsHistoricalOpenInterest)
	response, err := o.client.Get(params, url, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.OptionsHistoricalOpenInterest
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}
