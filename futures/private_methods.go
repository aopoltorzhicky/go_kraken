package futures

import (
	"fmt"
	"net/url"
	"strconv"
)

func (api *KrakenFutures) SendOrder(pair string, side string, orderType string, volume float64, price float64) (response Order, err error) {
	data := url.Values{
		"symbol":    {pair},
		"side":      {side},
		"orderType": {orderType},
		"size":      {strconv.FormatFloat(volume, 'f', 8, 64)},
		// "limitPrice": {strconv.FormatFloat(price, 'f', 8, 64)},
		// "triggerSignal": {"mark"},
		// "cliOrdId":      {"rage"},
		// "reduceOnly":    {"false"},
	}

	// Send the request
	var resp SendOrderResponse
	err = api.request("POST", "sendorder", true, data, &resp)
	if err != nil {
		return Order{}, err
	}

	// Check if the SendStatus and OrderEvents slice are present
	if resp.SendStatus.OrderEvents == nil || len(resp.SendStatus.OrderEvents) == 0 {
		return Order{}, fmt.Errorf("no order events found in response")
	}

	// Check if Order is present
	if resp.SendStatus.OrderEvents[0].Order == (Order{}) {
		return Order{}, fmt.Errorf("order is empty in the first order event")
	}

	response = resp.SendStatus.OrderEvents[0].Order
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
