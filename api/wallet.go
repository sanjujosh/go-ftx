package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/uscott/go-ftx/models"
)

const (
	apiGetCoins             = "/wallet/coins"
	apiGetBalances          = "/wallet/balances"
	apiGetBalancesAll       = "/wallet/all_balances"
	apiGetDepositAddress    = "/wallet/deposit_address/%s"
	apiGetDepositHistory    = "/wallet/deposits"
	apiGetWithdrawalHistory = "/wallet/withdrawals"
	apiRequestWithdrawal    = apiGetWithdrawalHistory
	apiGetAirdrops          = "wallet/airdrops"
	apiGetSavedAddresses    = "/wallet/saved_addresses"
	apiCreateSavedAddresses = apiGetSavedAddresses
	apiDeleteSavedAddresses = apiGetSavedAddresses
)

type Wallet struct {
	client *Client
}

func (w *Wallet) GetCoins() ([]*models.Coin, error) {

	request, err := w.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetCoins),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := w.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.Coin
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (w *Wallet) GetBalances() ([]*models.Balance, error) {

	request, err := w.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetBalances),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := w.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.Balance
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (w *Wallet) GetBalancesAllAccts() (map[string][]*models.Balance, error) {

	request, err := w.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetBalancesAll),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := w.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result map[string][]*models.Balance
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (w *Wallet) GetDepositAddress(
	coin string, method *models.DepositMethod,
) (*models.DepositAddress, error) {

	path := fmt.Sprintf(apiGetDepositAddress, coin)
	queryParams, err := PrepareQueryParams(&struct {
		Method *models.DepositMethod `json:"method"`
	}{Method: method})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := w.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, path),
		Params: queryParams,
	})

	response, err := w.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result models.DepositAddress
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return &result, nil
}

func (w *Wallet) GetDepositHistory(
	params *models.DepositHistoryParams,
) ([]*models.Deposit, error) {

	queryParams, err := PrepareQueryParams(params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := w.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetDepositHistory),
		Params: queryParams,
	})

	response, err := w.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.Deposit
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (w *Wallet) GetWithdrawalHistory(
	params *models.WithdrawalHistoryParams,
) ([]*models.Withdrawal, error) {

	queryParams, err := PrepareQueryParams(params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := w.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetWithdrawalHistory),
		Params: queryParams,
	})

	response, err := w.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.Withdrawal
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (w *Wallet) RequestWithdrawal(
	params *models.RequestWithdrawalParams,
) (*models.Withdrawal, error) {

	url := fmt.Sprintf("%s%s", apiUrl, apiRequestWithdrawal)
	response, err := w.client.Post(params, url)
	if err != nil {
		return nil, err
	}

	var result models.Withdrawal
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (w *Wallet) GetAirdrops(params *models.AirDropParams) ([]*models.AirDrop, error) {

	url := fmt.Sprintf("%s%s", apiUrl, apiGetAirdrops)
	response, err := w.client.Get(params, url, true)
	if err != nil {
		return nil, err
	}

	var result []*models.AirDrop
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (w *Wallet) GetSavedAddresses(coin *string) ([]*models.SavedAddress, error) {

	url := fmt.Sprintf("%s%s", apiUrl, apiGetSavedAddresses)
	params := struct {
		Coin *string `json:"coin,omitempty"`
	}{Coin: coin}
	response, err := w.client.Get(&params, url, true)
	if err != nil {
		return nil, err
	}

	var result []*models.SavedAddress
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (w *Wallet) CreateSavedAddresses(params *models.SavedAddressParams) ([]*models.SavedAddress, error) {

	url := fmt.Sprintf("%s%s", apiUrl, apiCreateSavedAddresses)
	response, err := w.client.Post(params, url)
	if err != nil {
		return nil, err
	}

	var result []*models.SavedAddress
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (w *Wallet) DeleteSavedAddress(address int64) error {

	params := struct {
		SavedAddressID *int64 `json:"saved_address_id"`
	}{SavedAddressID: &address}

	body, err := json.Marshal(&params)
	if err != nil {
		return errors.WithStack(err)
	}

	request, err := w.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodDelete,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiDeleteSavedAddresses),
		Body:   body,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = w.client.do(request)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
