package api

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/uscott/go-ftx/models"
	"github.com/uscott/go-tools/errs"
)

const (
	apiGetCoins             = "/wallet/coins"
	apiGetBalances          = "/wallet/balances"
	apiGetBalancesAll       = "/wallet/all_balances"
	apiGetDepositAddress    = "/wallet/deposit_address/%s"
	apiGetDepositHistory    = "/wallet/deposits"
	apiGetWithdrawalHistory = "/wallet/withdrawals"
	apiRequestWithdrawal    = apiGetWithdrawalHistory
	apiGetAirdrops          = "/wallet/airdrops"
	apiGetSavedAddresses    = "/wallet/saved_addresses"
	apiCreateSavedAddresses = apiGetSavedAddresses
	apiDeleteSavedAddresses = apiGetSavedAddresses
)

type Wallet struct {
	client *Client
}

func (w *Wallet) GetCoins() ([]*models.Coin, error) {

	url := FormURL(apiGetCoins)

	response, err := w.client.Get(nil, url, true)
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

	url := FormURL(apiGetBalances)

	response, err := w.client.Get(nil, url, true)
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

	url := FormURL(apiGetBalancesAll)

	response, err := w.client.Get(nil, url, true)
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
) (address, tag string, err error) {

	url := FormURL(fmt.Sprintf(apiGetDepositAddress, coin))

	params := &struct {
		Method *models.DepositMethod `json:"method,omitempty"`
	}{Method: method}

	response, err := w.client.Get(params, url, true)
	if err != nil {
		return address, tag, errors.WithStack(err)
	}

	result := models.DepositAddress{}
	if err = json.Unmarshal(response, &result); err != nil {
		return address, tag, errors.WithStack(err)
	}

	address, tag = result.Address, result.Tag

	return
}

func (w *Wallet) GetDepositHistory(
	params *models.DepositHistoryParams,
) ([]*models.Deposit, error) {

	url := FormURL(apiGetDepositHistory)

	response, err := w.client.Get(params, url, true)
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

	url := FormURL(apiGetWithdrawalHistory)

	response, err := w.client.Get(params, url, true)
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
	withdrawal *models.Withdrawal,
) (err error) {

	if withdrawal == nil {
		return errs.NilPtr
	}

	url := FormURL(apiRequestWithdrawal)

	response, err := w.client.Post(params, url)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = json.Unmarshal(response, withdrawal); err != nil {
		return errors.WithStack(err)
	}

	return
}

func (w *Wallet) GetAirdrops(params *models.AirDropParams) ([]*models.AirDrop, error) {

	url := FormURL(apiGetAirdrops)

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

	url := FormURL(apiGetSavedAddresses)

	params := &struct {
		Coin *string `json:"coin,omitempty"`
	}{Coin: coin}
	response, err := w.client.Get(params, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.SavedAddress
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (w *Wallet) CreateSavedAddresses(
	params *models.SavedAddressParams,
) ([]*models.SavedAddress, error) {

	url := FormURL(apiCreateSavedAddresses)

	response, err := w.client.Post(params, url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.SavedAddress
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (w *Wallet) DeleteSavedAddress(address int64) (result string, err error) {

	url := FormURL(apiDeleteSavedAddresses)

	params := &struct {
		SavedAddressID *int64 `json:"saved_address_id"`
	}{SavedAddressID: &address}

	response, err := w.client.Delete(params, url)
	if err != nil {
		return result, errors.WithStack(err)
	}

	if err = json.Unmarshal(response, &result); err != nil {
		return result, errors.WithStack(err)
	}

	return
}
