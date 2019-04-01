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
	if assets != nil && len(assets) > 0 {
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
	if pairs != nil && len(pairs) > 0 {
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
	if pairs != nil && len(pairs) > 0 {
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

func (api *Kraken) parseCandle(candle []interface{}) Candle {
	open, _ := strconv.ParseFloat(candle[1].(string), 64)
	high, _ := strconv.ParseFloat(candle[2].(string), 64)
	low, _ := strconv.ParseFloat(candle[3].(string), 64)
	close, _ := strconv.ParseFloat(candle[4].(string), 64)
	vwap, _ := strconv.ParseFloat(candle[5].(string), 64)
	volume, _ := strconv.ParseFloat(candle[6].(string), 64)

	return Candle{
		Time:      int64(candle[0].(float64)),
		Open:      open,
		High:      high,
		Low:       low,
		Close:     close,
		VolumeWAP: vwap,
		Volume:    volume,
		Count:     int64(candle[7].(float64)),
	}
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
	resp, err := api.request("OHLC", false, data, nil)
	if err != nil {
		return nil, err
	}
	parsedResp := resp.(map[string]interface{})
	ret := OHLCResponse{
		Last:    int64(parsedResp["last"].(float64)),
		Candles: make([]Candle, 0),
	}
	candles := parsedResp[pair].([]interface{})
	arr := make([]Candle, 0)
	for _, c := range candles {
		candle := api.parseCandle(c.([]interface{}))
		arr = append(arr, candle)
	}
	ret.Candles = arr
	return &ret, err
}

// GetOrderBook - Gets order book for `pair` with `depth`
func (api *Kraken) GetOrderBook(pair string, depth int64) (*OrderBook, error) {
	data := url.Values{
		"pair":  {pair},
		"depth": {strconv.FormatInt(depth, 10)},
	}
	resp, err := api.request("Depth", false, data, nil)
	if err != nil {
		return nil, err
	}
	parsedResp := resp.(map[string]interface{})
	book := parsedResp[pair]

}
