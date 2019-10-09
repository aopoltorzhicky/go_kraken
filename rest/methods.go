package rest

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Time - Gets server time. Note: This is to aid in approximating the skew time between the server and client.
func (api *Kraken) Time() (TimeResponse, error) {
	response := TimeResponse{}
	if err := api.request("Time", false, nil, &response); err != nil {
		return response, err
	}
	return response, nil
}

// Assets - Gets info about assets passed through `assets` arg.
// `assets` - array of needed assets. All by default if empty array passed or `assets` is nil.
func (api *Kraken) Assets(assets ...string) (map[string]Asset, error) {
	data := url.Values{}
	if len(assets) > 0 {
		data.Add("asset", strings.Join(assets, ","))
	} else {
		data = nil
	}
	response := make(map[string]Asset)
	if err := api.request("Assets", false, data, &response); err != nil {
		return response, err
	}
	return response, nil
}

// AssetPairs - Gets array of pair names and their info passed through `pairs` arg.
// `pairs` - array of needed pairs. All by default if empty array passed or `pairs` is nil.
func (api *Kraken) AssetPairs(pairs ...string) (map[string]AssetPair, error) {
	data := url.Values{}
	if len(pairs) > 0 {
		data.Add("pair", strings.Join(pairs, ","))
	} else {
		data = nil
	}
	response := make(map[string]AssetPair)
	if err := api.request("AssetPairs", false, data, &response); err != nil {
		return nil, err
	}
	return response, nil
}

// Ticker - Gets array of tickers passed through `pairs` arg.
// `pairs` - array of needed pairs. All by default if empty array passed or `pairs` is nil.
func (api *Kraken) Ticker(pairs ...string) (map[string]Ticker, error) {
	var data url.Values
	if len(pairs) > 0 {
		data = url.Values{
			"pair": {strings.Join(pairs, ",")},
		}
	} else {
		return nil, fmt.Errorf("You need to set pairs on Ticker request")
	}
	response := make(map[string]Ticker)
	if err := api.request("Ticker", false, data, &response); err != nil {
		return nil, err
	}
	return response, nil
}

// Candles - Get OHLC data
func (api *Kraken) Candles(pair string, interval int64, since int64) (OHLCResponse, error) {
	data := url.Values{
		"pair": {pair},
	}
	if since > 0 {
		data.Set("since", strconv.FormatInt(since, 10))
	}
	if interval > 1 {
		data.Set("interval", strconv.FormatInt(interval, 10))
	}
	response := OHLCResponse{}
	if err := api.request("OHLC", false, data, &response); err != nil {
		return response, err
	}
	return response, nil
}

// GetOrderBook - Gets order book for `pair` with `depth`
func (api *Kraken) GetOrderBook(pair string, depth int64) (map[string]OrderBook, error) {
	data := url.Values{
		"pair":  {pair},
		"count": {strconv.FormatInt(depth, 10)},
	}
	response := make(map[string]OrderBook)
	if err := api.request("Depth", false, data, &response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetTrades - returns trades on pair from since date
func (api *Kraken) GetTrades(pair string, since int64) (TradeResponse, error) {
	data := url.Values{
		"pair": {pair},
	}
	if since > 0 {
		data.Add("since", strconv.FormatInt(since, 10))
	}
	response := TradeResponse{}
	if err := api.request("Trades", false, data, &response); err != nil {
		return response, err
	}
	return response, nil
}

// GetSpread - return array of pair name and recent spread data
func (api *Kraken) GetSpread(pair string, since int64) (SpreadResponse, error) {
	data := url.Values{
		"pair": {pair},
	}
	if since > 0 {
		data.Add("since", strconv.FormatInt(since, 10))
	}
	response := SpreadResponse{}
	if err := api.request("Spread", false, data, &response); err != nil {
		return response, err
	}
	return response, nil
}
