package rest

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Time - Gets server time. Note: This is to aid in approximating the skew time between the server and client.
func (api *Kraken) Time() (*TimeResponse, error) {
	resp, err := api.request("Time", false, nil, &TimeResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*TimeResponse), err
}

// Assets - Gets info about assets passed through `assets` arg.
// `assets` - array of needed assets. All by default if empty array passed or `assets` is nil.
func (api *Kraken) Assets(assets ...string) (*AssetResponse, error) {
	data := url.Values{}
	if len(assets) > 0 {
		data.Add("asset", strings.Join(assets, ","))
	} else {
		data = nil
	}
	resp, err := api.request("Assets", false, data, &AssetResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*AssetResponse), err
}

// AssetPairs - Gets array of pair names and their info passed through `pairs` arg.
// `pairs` - array of needed pairs. All by default if empty array passed or `pairs` is nil.
func (api *Kraken) AssetPairs(pairs ...string) (*AssetPairsResponse, error) {
	data := url.Values{}
	if len(pairs) > 0 {
		data.Add("pair", strings.Join(pairs, ","))
	} else {
		data = nil
	}
	resp, err := api.request("AssetPairs", false, data, &AssetPairsResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*AssetPairsResponse), err
}

// Ticker - Gets array of tickers passed through `pairs` arg.
// `pairs` - array of needed pairs. All by default if empty array passed or `pairs` is nil.
func (api *Kraken) Ticker(pairs ...string) (*TickerResponse, error) {
	var data url.Values
	if len(pairs) > 0 {
		data = url.Values{
			"pair": {strings.Join(pairs, ",")},
		}
	} else {
		return nil, fmt.Errorf("You need to set pairs on Ticker request")
	}
	resp, err := api.request("Ticker", false, data, &TickerResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*TickerResponse), err
}

// Candles - Get OHLC data
func (api *Kraken) Candles(pair string, interval int64, since int64) (*OHLCResponse, error) {
	data := url.Values{
		"pair": {pair},
	}
	if since > 0 {
		data.Set("since", strconv.FormatInt(since, 10))
	}
	if interval > 1 {
		data.Set("interval", strconv.FormatInt(interval, 10))
	}
	resp, err := api.request("OHLC", false, data, &OHLCResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*OHLCResponse), err
}

// GetOrderBook - Gets order book for `pair` with `depth`
func (api *Kraken) GetOrderBook(pair string, depth int64) (*BookResponse, error) {
	data := url.Values{
		"pair":  {pair},
		"count": {strconv.FormatInt(depth, 10)},
	}
	resp, err := api.request("Depth", false, data, &BookResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*BookResponse), nil
}

// GetTrades - returns trades on pair from since date
func (api *Kraken) GetTrades(pair string, since int64) (*TradeResponse, error) {
	data := url.Values{
		"pair": {pair},
	}
	if since > 0 {
		data.Add("since", strconv.FormatInt(since, 10))
	}
	resp, err := api.request("Trades", false, data, &TradeResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*TradeResponse), nil
}

// GetSpread - return array of pair name and recent spread data
func (api *Kraken) GetSpread(pair string, since int64) (*SpreadResponse, error) {
	data := url.Values{
		"pair": {pair},
	}
	if since > 0 {
		data.Add("since", strconv.FormatInt(since, 10))
	}
	resp, err := api.request("Spread", false, data, &SpreadResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*SpreadResponse), nil
}
