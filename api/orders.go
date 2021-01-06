package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/uscott/go-ftx/models"
)

const (
	apiGetOpenOrders           = "/orders"
	apiGetOrdersHistory        = "/orders/history"
	apiGetTriggerOrders        = "/conditional_orders"
	apiGetOrderTriggers        = "/conditional_orders/%d/triggers"
	apiGetTriggerOrdersHistory = "/conditional_orders/history"
	apiPlaceOrder              = apiGetOpenOrders
	apiPlaceTriggerOrder       = apiGetTriggerOrders
	apiModifyOrder             = "/orders/%d/modify"
	apiModifyTriggerOrder      = "/conditional_orders/%d/modify"
	apiGetOrderStatus          = apiGetOpenOrders
	apiCancelOrder             = apiGetOpenOrders
	apiCancelTriggerOrder      = apiGetTriggerOrders
	apiCancelAll               = apiGetOpenOrders
)

type Orders struct {
	client *Client
}

var errNilPtr = fmt.Errorf("nil pointer")

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

	queryParams, err := PrepareQueryParams(params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetOrdersHistory),
		Params: queryParams,
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

func (o *Orders) GetOpenTriggerOrders(
	params *models.OpenTriggerOrdersParams) ([]*models.TriggerOrder, error) {

	queryParams, err := PrepareQueryParams(params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetTriggerOrders),
		Params: queryParams,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
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

func (o *Orders) GetOrderTriggers(orderID int64) ([]*models.Trigger, error) {

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, fmt.Sprintf(apiGetOrderTriggers, orderID)),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
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

	queryParams, err := PrepareQueryParams(params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetTriggerOrdersHistory),
		Params: queryParams,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
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

/*
func (o *Orders) PlaceOrder(params *models.OrderParams) (*models.Order, error) {

	if params == nil {
		return nil, errNilPtr
	}
	body, err := json.Marshal(*params)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiPlaceOrder),
		Body:   body,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result models.Order
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}
*/

func (o *Orders) PlaceOrder(params *models.OrderParams) (*models.Order, error) {

	url := fmt.Sprintf("%s%s", apiUrl, apiPlaceOrder)
	response, err := o.client.Post(params, url)
	if err != nil {
		return nil, err
	}

	var result models.Order
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (o *Orders) PlaceTriggerOrder(
	params *models.TriggerOrderParams) (*models.TriggerOrder, error) {

	if params == nil {
		return nil, errNilPtr
	}
	body, err := json.Marshal(*params)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiPlaceTriggerOrder),
		Body:   body,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result models.TriggerOrder
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (o *Orders) ModifyOrder(
	orderID int64, params *models.ModifyOrderParams) (*models.Order, error) {

	if params == nil {
		return nil, errNilPtr
	}
	path := fmt.Sprintf(apiModifyOrder, orderID)

	body, err := json.Marshal(*params)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, path),
		Body:   body,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result models.Order
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (o *Orders) ModifyTriggerOrder(
	orderID int64,
	params *models.ModifyTriggerOrderParams) (*models.TriggerOrder, error) {

	if params == nil {
		return nil, errNilPtr
	}
	path := fmt.Sprintf(apiModifyTriggerOrder, orderID)
	body, err := json.Marshal(*params)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, path),
		Body:   body,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result models.TriggerOrder
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (o *Orders) GetOrderStatus(orderID int64) (*models.Order, error) {

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s/%d", apiUrl, apiGetOrderStatus, orderID),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result models.Order
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (o *Orders) CancelOrder(orderID int64) error {

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodDelete,
		URL:    fmt.Sprintf("%s%s/%d", apiUrl, apiCancelOrder, orderID),
	})
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = o.client.do(request)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (o *Orders) CancelTriggerOrder(orderID int64) error {

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodDelete,
		URL:    fmt.Sprintf("%s%s/%d", apiUrl, apiCancelTriggerOrder, orderID),
	})
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = o.client.do(request)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (o *Orders) CancelAllOrders(params *models.CancelAllParams) error {

	if params == nil {
		return errNilPtr
	}
	body, err := json.Marshal(*params)
	if err != nil {
		return errors.WithStack(err)
	}
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodDelete,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiCancelAll),
		Body:   body,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = o.client.do(request)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
