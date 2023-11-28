package futures

import (
	"fmt"
	"net/url"
	"strconv"
)

func (api *KrakenFutures) SendOrder(symbol string, side string, orderType string, volume float64, price float64) (response SendStatus, err error) {
	data := url.Values{
		"orderType":  {orderType},
		"side":       {side},
		"size":       {strconv.FormatFloat(volume, 'f', -1, 64)},
		"symbol":     {symbol},
		"limitPrice": {strconv.FormatFloat(price, 'f', -1, 64)},
		// "triggerSignal": {"mark"},
		// "cliOrdId":      {"rage"},
		// "reduceOnly":    {"false"},
	}

	// Send the request
	var resp SendOrderResponse
	err = api.request("POST", "sendorder", true, data, &resp)
	if err != nil {
		return SendStatus{}, err
	}

	// Check if the SendStatus and OrderEvents slice are present
	if resp.SendStatus.OrderEvents == nil || len(resp.SendStatus.OrderEvents) == 0 {
		return SendStatus{}, fmt.Errorf("no order events found in response")
	}

	response = resp.SendStatus
	return response, nil
}

func (api *KrakenFutures) GetOrderStatus(cliOrdIds []string, orderIds []string) (response OrderStatusResponse, err error) {
	data := url.Values{}

	for _, id := range cliOrdIds {
		data.Add("cliOrdIds", id)
	}
	for _, id := range orderIds {
		data.Add("orderIds", id)
	}

	// Send the request
	err = api.request("POST", "orders/status", true, data, &response)
	if err != nil {
		return OrderStatusResponse{}, err
	}

	return response, nil
}

func (api *KrakenFutures) GetBalances() (response SendOrderResponse, err error) {
	// Send the request
	var resp SendOrderResponse
	err = api.request("GET", "accounts", true, nil, &resp)
	if err != nil {
		return SendOrderResponse{}, err
	}

	return resp, nil
}
