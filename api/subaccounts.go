package api

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sanjujosh/go-ftx/models"
)

const (
	apiSubaccounts           = "/subaccounts"
	apiChangeSubaccountName  = "/subaccounts/update_name"
	apiGetSubaccountBalances = "/subaccounts/%s/balances"
	apiTransfer              = "/subaccounts/transfer"
)

type SubAccounts struct {
	client *Client
}

func (s *SubAccounts) GetSubaccounts() ([]*models.SubAccount, error) {

	url := FormURL(apiSubaccounts)

	response, err := s.client.Get(nil, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.SubAccount

	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SubAccounts) CreateSubaccount(nickname string) (*models.SubAccount, error) {

	url := FormURL(apiSubaccounts)

	pars := &struct {
		Nickname string `json:"nickname"`
	}{Nickname: nickname}

	response, err := s.client.Post(pars, url)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result models.SubAccount

	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return &result, nil
}

func (s *SubAccounts) ChangeSubaccount(nickname, newNickname string) (result string, err error) {

	url := FormURL(apiChangeSubaccountName)

	pars := &struct {
		Nickname    string `json:"nickname"`
		NewNickname string `json:"newNickname"`
	}{Nickname: nickname, NewNickname: newNickname}

	response, err := s.client.Post(pars, url)
	if err != nil {
		return result, errors.WithStack(err)
	}

	if err = json.Unmarshal(response, &result); err != nil {
		return result, errors.WithStack(err)
	}
	return
}

func (s *SubAccounts) DeleteSubaccount(nickname string) (result string, err error) {

	url := FormURL(apiSubaccounts)

	pars := &struct {
		Nickname string `json:"nickname"`
	}{Nickname: nickname}

	response, err := s.client.Delete(pars, url)

	if err != nil {
		return result, errors.WithStack(err)
	}

	if err = json.Unmarshal(response, &result); err != nil {
		return result, errors.WithStack(err)
	}

	return
}

func (s *SubAccounts) GetSubaccountBalances(nickname string) ([]*models.Balance, error) {

	url := FormURL(fmt.Sprintf(apiGetSubaccountBalances, nickname))

	response, err := s.client.Get(nil, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.Balance

	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SubAccounts) Transfer(payload *models.TransferPayload) (*models.TransferResponse, error) {

	url := FormURL(apiTransfer)

	response, err := s.client.Post(payload, url)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result models.TransferResponse
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &result, nil
}
