package futures

import (
	"log"
	"net/url"
	"strconv"
)

func (api *KrakenFutures) SendOrder(pair string, side string, orderType string, volume float64, price float64, args map[string]interface{}) (response SendOrderResponse, err error) {
	data := url.Values{
		"pair":      {pair},
		"type":      {side},
		"ordertype": {orderType},
		"volume":    {strconv.FormatFloat(volume, 'f', 8, 64)},
		"price":     {strconv.FormatFloat(price, 'f', 8, 64)},
	}

	// Add any additional arguments provided
	for key, value := range args {
		switch v := value.(type) {
		case string:
			data.Set(key, v)
		case int64:
			data.Set(key, strconv.FormatInt(v, 10))
		case float64:
			data.Set(key, strconv.FormatFloat(v, 'f', 8, 64))
		case bool:
			data.Set(key, strconv.FormatBool(v))
		default:
			log.Printf("[WARNING] Unknown value type %v for key %s", value, key)
		}
	}

	// Send the request
	err = api.request("POST", "sendOrder", true, data, &response)
	return
}
