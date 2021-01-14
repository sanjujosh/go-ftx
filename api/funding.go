package api

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/uscott/go-ftx/models"
)

const (
	apiGetFundingPayments = "/funding_payments"
)

type Funding struct {
	client *Client
}

func (f *Funding) GetFundingPayments(start, end *int64) ([]*models.FundingPayment, error) {

	url := FormURL(apiGetFundingPayments)
	params := &models.FundingPaymentParams{StartTime: start, EndTime: end}

	response, err := f.client.Get(params, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result []*models.FundingPayment
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}
