package api

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/uscott/go-ftx/models"
)

const (
	apiGetBorrowRates     = "/spot_margin/borrow_rates"
	apiGetLendingRates    = "/spot_margin/lending_rates"
	apiGetBorrowSummary   = "/spot_margin/borrow_summary"
	apiGetMarketInfo      = "/spot_margin/market_info"
	apiGetBorrowHistory   = "/spot_margin/borrow_history"
	apiGetLendingHistory  = "/spot_margin/lending_history"
	apiGetLendingOffers   = "/spot_margin/offers"
	apiGetLendingInfo     = "/spot_margin/lending_info"
	apiSubmitLendingOffer = "/spot_margin/offers"
)

type SpotMargin struct {
	client *Client
}

func (s *SpotMargin) GetBorrowRates() ([]*models.BorrowRate, error) {

	url := fmt.Sprintf("%s%s", apiUrl, apiGetBorrowRates)

	response, err := s.client.Get(&struct{}{}, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.BorrowRate

	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (s *SpotMargin) GetLendingRates() ([]*models.LendingRate, error) {

	url := fmt.Sprintf("%s%s", apiUrl, apiGetLendingRates)

	response, err := s.client.Get(&struct{}{}, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result []*models.LendingRate

	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) GetBorrowSummary() ([]*models.BorrowedAmount, error) {

	url := FormURL(apiGetBorrowSummary)

	response, err := s.client.Get(&struct{}{}, url, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.BorrowedAmount

	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) GetMarketInfo(market string) (*models.SpotMarginMarketInfo, error) {

	url := FormURL(apiGetMarketInfo)

	params := struct {
		Market *string `json:"market"`
	}{Market: &market}

	response, err := s.client.Get(&params, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := models.SpotMarginMarketInfo{}

	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return &result, nil
}

func (s *SpotMargin) GetBorrowHistory() ([]*models.BorrowHistory, error) {

	url := FormURL(apiGetBorrowHistory)

	response, err := s.client.Get(&struct{}{}, url, true)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.BorrowHistory

	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) GetLendingHistory() ([]*models.LendingHistory, error) {

	url := FormURL(apiGetLendingHistory)

	response, err := s.client.Get(&struct{}{}, url, true)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.LendingHistory

	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) GetLendingOffers() ([]*models.LendingOffer, error) {

	url := FormURL(apiGetLendingOffers)

	response, err := s.client.Get(&struct{}{}, url, true)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.LendingOffer

	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) GetLendingInfo() ([]*models.LendingInfo, error) {

	url := FormURL(apiGetLendingInfo)

	response, err := s.client.Get(&struct{}{}, url, true)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.LendingInfo

	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) SubmitLendingOffer(
	coin string, size decimal.Decimal, rate float64,
) (result string, err error) {

	url := FormURL(apiSubmitLendingOffer)
	params := &models.LendingOfferParams{
		Coin: &coin,
		Size: &size,
		Rate: &rate,
	}
	response, err := s.client.Get(params, url, true)
	if err != nil {
		return result, errors.WithStack(err)
	}

	if err = json.Unmarshal(response, &result); err != nil {
		return result, errors.WithStack(err)
	}

	return
}
