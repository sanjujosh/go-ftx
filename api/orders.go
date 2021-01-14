package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/uscott/go-ftx/models"
	"github.com/uscott/go-tools/errs"
)

const (
	apiGetOpenOrders            = "/orders"
	apiGetOrdersHistory         = "/orders/history"
	apiGetTriggerOrders         = "/conditional_orders"
	apiGetOrderTriggers         = "/conditional_orders/%d/triggers"
	apiGetTriggerOrdersHistory  = "/conditional_orders/history"
	apiPlaceOrder               = apiGetOpenOrders
	apiPlaceTriggerOrder        = apiGetTriggerOrders
	apiModifyOrder              = "/orders/%d/modify"
	apiModifyOrderByClientID    = "/orders/by_client_id/%d/modify"
	apiModifyTriggerOrder       = "/conditional_orders/%d/modify"
	apiGetOrderStatus           = "/orders/%d"
	apiGetOrderStatusByClientID = "/orders/by_client_id/%d"
	apiCancelOrder              = apiGetOrderStatus
	apiCancelOrderByClientID    = apiGetOrderStatusByClientID
	apiCancelTriggerOrder       = "/conditional_orders/%d"
	apiCancelAll                = apiGetOpenOrders
)

type Orders struct {
	client *Client
}

func (o *Orders) GetOpenOrders(market string) ([]*models.Order, error) {

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetOpenOrders),
		Params: map[string]string{
			"market": market,
		},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.Order
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) GetOrdersHistory(
	params *models.OrdersHistoryParams) ([]*models.Order, error) {

	url := FormURL(apiGetOrdersHistory)
	response, err := o.client.Get(params, url, true)
	if err != nil {
		return nil, err
	}

	var result []*models.Order
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) GetOpenTriggerOrders(
	params *models.OpenTriggerOrdersParams) ([]*models.TriggerOrder, error) {

	url := FormURL(apiGetTriggerOrders)
	response, err := o.client.Get(params, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.TriggerOrder
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) GetTriggerOrderTriggers(orderID int64) ([]*models.Trigger, error) {

	url := FormURL(fmt.Sprintf(apiGetOrderTriggers, orderID))
	response, err := o.client.Get(nil, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.Trigger
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) GetTriggerOrdersHistory(
	params *models.TriggerOrdersHistoryParams) ([]*models.TriggerOrder, error) {

	url := FormURL(apiGetTriggerOrdersHistory)
	response, err := o.client.Get(params, url, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.TriggerOrder
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) PlaceOrder(params *models.OrderParams, order *models.Order) (err error) {

	if order == nil {
		return errs.NilPtr
	}
	url := fmt.Sprintf("%s%s", apiUrl, apiPlaceOrder)
	response, err := o.client.Post(params, url)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(response, &order); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (o *Orders) PlaceTriggerOrder(
	params *models.TriggerOrderParams, order *models.TriggerOrder) (err error) {

	if order == nil {
		return errs.NilPtr
	}
	url := FormURL(apiPlaceTriggerOrder)
	response, err := o.client.Post(params, url)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = json.Unmarshal(response, order); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (o *Orders) ModifyOrder(
	orderID int64,
	params *models.ModifyOrderParams,
	order *models.Order) (err error) {

	if params == nil || order == nil {
		return models.ErrNilPtr
	}
	url := FormURL(fmt.Sprintf(apiModifyOrder, orderID))
	response, err := o.client.Post(params, url)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = json.Unmarshal(response, order); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (o *Orders) ModifyOrderByClientID(
	clientID int64, params *models.ModifyOrderParams, order *models.Order,
) (err error) {

	if order == nil {
		return errs.NilPtr
	}
	url := FormURL(fmt.Sprintf(apiModifyOrderByClientID, clientID))
	params.ClientID = nil
	response, err := o.client.Post(params, url)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = json.Unmarshal(response, order); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (o *Orders) ModifyTriggerOrder(
	orderID int64,
	params *models.ModifyTriggerOrderParams,
	order *models.TriggerOrder) (err error) {

	if params == nil || order == nil {
		return models.ErrNilPtr
	}
	url := FormURL(fmt.Sprintf(apiModifyTriggerOrder, orderID))
	response, err := o.client.Post(params, url)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = json.Unmarshal(response, order); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (o *Orders) GetOrderStatus(orderID int64, order *models.Order) (err error) {

	if order == nil {
		return errs.NilPtr
	}
	url := FormURL(fmt.Sprintf(apiGetOrderStatus, orderID))
	response, err := o.client.Get(nil, url, true)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = json.Unmarshal(response, order); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (o *Orders) GetOrderStatusByClientID(clientID int64, order *models.Order) (err error) {

	if order == nil {
		return errs.NilPtr
	}
	url := FormURL(fmt.Sprintf(apiGetOrderStatusByClientID, clientID))
	response, err := o.client.Get(nil, url, true)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = json.Unmarshal(response, order); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (o *Orders) CancelOrder(orderID int64) (result *models.Result, err error) {

	url := FormURL(fmt.Sprintf(apiCancelOrder, orderID))
	response, err := o.client.Delete(nil, url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result = new(models.Result)
	if err = json.Unmarshal(response, result); err != nil {
		return nil, errors.WithStack(err)
	}
	return
}

func (o *Orders) CancelOrderByClientID(clientID int64) (result *models.Result, err error) {

	url := FormURL(fmt.Sprintf(apiCancelOrderByClientID, clientID))
	response, err := o.client.Delete(nil, url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result = new(models.Result)
	if err = json.Unmarshal(response, result); err != nil {
		return nil, errors.WithStack(err)
	}
	return
}

func (o *Orders) CancelTriggerOrder(orderID int64) (result *models.Result, err error) {

	url := FormURL(fmt.Sprintf(apiCancelTriggerOrder, orderID))
	response, err := o.client.Delete(nil, url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result = new(models.Result)
	if err = json.Unmarshal(response, result); err != nil {
		return nil, errors.WithStack(err)
	}
	return
}

func (o *Orders) CancelAllOrders(
	params *models.CancelAllParams) (result *models.Result, err error) {

	url := FormURL(apiCancelAll)
	response, err := o.client.Delete(params, url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result = new(models.Result)
	if err = json.Unmarshal(response, result); err != nil {
		return nil, errors.WithStack(err)
	}
	return
}
