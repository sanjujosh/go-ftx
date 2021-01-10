package api

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/uscott/go-ftx/models"
)

const (
	apiListLeveragedTokens                  = "/lt/tokens"
	apiGetTokenInfo                         = "/lt/%s"
	apiGetLeveragedTokenBalances            = "/lt/balances"
	apiListLeveragedTokenCreationRequests   = "/lt/creations"
	apiRequestLeveragedTokenCreation        = "/lt/%s/create"
	apiListLeveragedTokenRedemptionRequests = "/lt/redemptions"
	apiRequestLeveragedTokenRedemption      = "/lt/%s/redeem"
)

type LeveragedTokens struct {
	client *Client
}

func (l *LeveragedTokens) ListLeveragedTokens() ([]*models.LeveragedToken, error) {

	url := FormURL(apiListLeveragedTokens)

	response, err := l.client.Get(&struct{}{}, url, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.LeveragedToken
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (l *LeveragedTokens) GetTokenInfo(token string) (*models.TokenInfo, error) {

	url := FormURL(fmt.Sprintf(apiGetTokenInfo, token))

	response, err := l.client.Get(&struct{}{}, url, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := models.TokenInfo{}
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (l *LeveragedTokens) GetLeveragedTokenBalances() (
	[]*models.LeveragedTokenBalance, error) {

	url := FormURL(apiGetLeveragedTokenBalances)

	response, err := l.client.Get(&struct{}{}, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.LeveragedTokenBalance
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (l *LeveragedTokens) ListLeveragedTokenCreationRequests() (
	[]*models.LeveragedTokenCreationRequest, error) {

	url := FormURL(apiListLeveragedTokenCreationRequests)

	response, err := l.client.Get(&struct{}{}, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.LeveragedTokenCreationRequest
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (l *LeveragedTokens) RequestLeveragedTokenCreation(
	token string, size decimal.Decimal,
) (*models.LeveragedTokenCreation, error) {

	url := FormURL(fmt.Sprintf(apiRequestLeveragedTokenCreation, token))
	body := struct {
		Size *decimal.Decimal `json:"size"`
	}{Size: &size}
	response, err := l.client.Post(&body, url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := models.LeveragedTokenCreation{}
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (l *LeveragedTokens) ListLeveragedTokenRedemptionRequests() (
	[]*models.LeveragedTokenRedemptionRequest, error) {

	url := FormURL(apiListLeveragedTokenRedemptionRequests)

	response, err := l.client.Get(&struct{}{}, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.LeveragedTokenRedemptionRequest
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (l *LeveragedTokens) RequestLeveragedTokenRedemption(
	token string, size decimal.Decimal,
) (*models.LeveragedTokenRedemption, error) {

	url := FormURL(fmt.Sprintf(apiRequestLeveragedTokenRedemption, token))
	body := struct {
		Size *decimal.Decimal `json:"size"`
	}{Size: &size}
	response, err := l.client.Post(&body, url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	result := models.LeveragedTokenRedemption{}
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}
