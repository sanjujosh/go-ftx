package api

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/uscott/go-ftx/models"
)

const (
	apiGetFundingPayments = "/funding_payments"
)

type Funding struct {
	client *Client
}

func (f *Funding) GetFundingPayments(
	future *string, start, end *int64,
) ([]*models.FundingPayment, error) {

	params := struct {
		Future    *string `json:"future,omitempty"`
		StartTime *int64  `json:"start_time,omitempty"`
		EndTime   *int64  `json:"end_time,omitempty"`
	}{Future: future, StartTime: start, EndTime: end}

	url := fmt.Sprintf("%s%s", apiUrl, apiGetFundingPayments)

	response, err := f.client.Get(&params, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result []*models.FundingPayment
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}
