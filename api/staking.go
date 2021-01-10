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
	apiGetStakes            = "/staking/stakes"
	apiGetUnstakeRequests   = "/staking/unstake_requests"
	apiGetStakeBalances     = "/staking/balances"
	apiRequestUnstake       = apiGetUnstakeRequests
	apiCancelUnstakeRequest = "/staking/unstake_requests/%d"
	apiGetStakingRewards    = "/staking/staking_rewards"
	apiRequestStake         = "/srm_stakes/stakes"
)

type Staking struct {
	client *Client
}

func (s *Staking) GetStakes() ([]*models.Stake, error) {

	url := FormURL(apiGetStakes)
	response, err := s.client.Get(&struct{}{}, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.Stake
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (s *Staking) GetUnstakeRequests() ([]*models.UnstakeRequest, error) {

	url := FormURL(apiGetUnstakeRequests)
	response, err := s.client.Get(&struct{}{}, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.UnstakeRequest
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (s *Staking) GetStakeBalances() ([]*models.StakeBalance, error) {

	url := FormURL(apiGetStakeBalances)
	response, err := s.client.Get(&struct{}{}, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.StakeBalance
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (s *Staking) RequestUnstake(
	coin string, size decimal.Decimal,
) (*models.UnstakeRequest, error) {

	url := FormURL(apiRequestUnstake)
	params := &struct {
		Coin *string          `json:"coin"`
		Size *decimal.Decimal `json:"size"`
	}{Coin: &coin, Size: &size}

	response, err := s.client.Post(params, url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := models.UnstakeRequest{}
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (s *Staking) CancelUnstakeRequest(id int64) (*models.Succeeded, error) {

	url := FormURL(fmt.Sprintf(apiCancelUnstakeRequest, id))
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodDelete,
		URL:    url,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := models.Succeeded{}
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (s *Staking) GetStakingRewards() ([]*models.StakingReward, error) {

	url := FormURL(apiGetStakingRewards)
	response, err := s.client.Get(&struct{}{}, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.StakingReward
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (s *Staking) RequestStake(coin string, size decimal.Decimal) (*models.Stake, error) {

	url := FormURL(apiRequestStake)
	params := &struct {
		Coin *string          `json:"coin"`
		Size *decimal.Decimal `json:"size"`
	}{Coin: &coin, Size: &size}
	response, err := s.client.Post(params, url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := models.Stake{}
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}
