package rest

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

// GetAccountBalances - methods returns account balances
func (api *Kraken) GetAccountBalances() (*BalanceResponse, error) {
	resp, err := api.request("Balance", true, nil, &BalanceResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*BalanceResponse), nil
}

// GetTradeBalance - returns tradable balances info
func (api *Kraken) GetTradeBalance(baseAsset string) (*TradeBalanceResponse, error) {
	if baseAsset == "" {
		baseAsset = "ZUSD"
	}
	data := url.Values{
		"asset": {baseAsset},
	}

	resp, err := api.request("TradeBalance", true, data, &TradeBalanceResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*TradeBalanceResponse), nil
}

// GetOpenOrders - returns account open order
func (api *Kraken) GetOpenOrders(needTrades bool, userRef string) (*OpenOrdersResponse, error) {
	data := url.Values{}
	if needTrades {
		data.Set("trades", "true")
	}
	if userRef != "" {
		data.Set("userref", userRef)
	}

	resp, err := api.request("OpenOrders", true, data, &OpenOrdersResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*OpenOrdersResponse), nil
}

// GetClosedOrders - returns account closed order
func (api *Kraken) GetClosedOrders(needTrades bool, userRef string, start int64, end int64) (*ClosedOrdersResponse, error) {
	data := url.Values{}
	if needTrades {
		data.Set("trades", "true")
	}
	if userRef != "" {
		data.Set("userref", userRef)
	}
	if start != 0 {
		data.Set("start", strconv.FormatInt(start, 10))
	}
	if end != 0 {
		data.Set("end", strconv.FormatInt(end, 10))
	}

	resp, err := api.request("ClosedOrders", true, data, &ClosedOrdersResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*ClosedOrdersResponse), nil
}

// QueryOrders - returns account's order by IDs
func (api *Kraken) QueryOrders(needTrades bool, userRef string, txIDs ...string) (*map[string]OrderInfo, error) {
	data := url.Values{}
	if needTrades {
		data.Set("trades", "true")
	}
	if userRef != "" {
		data.Set("userref", userRef)
	}
	if len(txIDs) > 50 {
		return nil, fmt.Errorf("Maximum count of requested orders is 50")
	} else if len(txIDs) == 0 {
		return nil, fmt.Errorf("txIDs is required")
	} else {
		data.Set("txid", strings.Join(txIDs, ","))
	}

	resp, err := api.request("QueryOrders", true, data, &map[string]OrderInfo{})
	if err != nil {
		return nil, err
	}
	return resp.(*map[string]OrderInfo), nil
}

// GetTradesHistory - returns account's trade history
func (api *Kraken) GetTradesHistory(tradeType string, needTrades bool, start int64, end int64) (*TradesHistoryResponse, error) {
	data := url.Values{
		"type": {"all"},
	}
	if needTrades {
		data.Set("trades", "true")
	}
	if tradeType != "" {
		data.Set("type", tradeType)
	}
	if start != 0 {
		data.Set("start", strconv.FormatInt(start, 10))
	}
	if end != 0 {
		data.Set("end", strconv.FormatInt(end, 10))
	}
	resp, err := api.request("TradesHistory", true, data, &TradesHistoryResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*TradesHistoryResponse), nil
}

// QueryTrades - returns trades by IDs
func (api *Kraken) QueryTrades(trades bool, txIDs ...string) (*map[string]PrivateTrade, error) {
	data := url.Values{}
	if trades {
		data.Set("trades", "true")
	}
	if len(txIDs) == 0 {
		return nil, fmt.Errorf("txIDs is required")
	}
	data.Set("txid", strings.Join(txIDs, ","))

	resp, err := api.request("QueryTrades", true, data, &map[string]PrivateTrade{})
	if err != nil {
		return nil, err
	}
	return resp.(*map[string]PrivateTrade), nil
}

// GetOpenPositions - returns list of open positions
func (api *Kraken) GetOpenPositions(docalcs bool, txIDs ...string) (*map[string]Position, error) {
	data := url.Values{}
	if docalcs {
		data.Set("docalcs", "true")
	}
	if len(txIDs) == 0 {
		return nil, fmt.Errorf("txIDs is required")
	}
	data.Set("txid", strings.Join(txIDs, ","))

	resp, err := api.request("OpenPositions", true, data, &map[string]Position{})
	if err != nil {
		return nil, err
	}
	return resp.(*map[string]Position), nil
}

// GetLedgersInfo - returns ledgers info
func (api *Kraken) GetLedgersInfo(ledgerType string, start int64, end int64, assets ...string) (*LedgerInfoResponse, error) {
	data := url.Values{}
	if ledgerType != "" {
		data.Set("type", LedgerTypeAll)
	}
	if start != 0 {
		data.Set("start", strconv.FormatInt(start, 10))
	}
	if end != 0 {
		data.Set("end", strconv.FormatInt(end, 10))
	}
	if len(assets) == 0 {
		return nil, fmt.Errorf("`assets` is required")
	}
	data.Set("assets", strings.Join(assets, ","))

	resp, err := api.request("Ledgers", true, data, &LedgerInfoResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*LedgerInfoResponse), nil
}

// QueryLedgers - get ledgers by ID
func (api *Kraken) QueryLedgers(ledgerIds ...string) (*map[string]Ledger, error) {
	data := url.Values{}
	if len(ledgerIds) == 0 {
		return nil, fmt.Errorf("`ledgerIds` is required")
	}
	data.Set("id", strings.Join(ledgerIds, ","))

	resp, err := api.request("QueryLedgers", true, data, &map[string]Ledger{})
	if err != nil {
		return nil, err
	}
	return resp.(*map[string]Ledger), nil
}

// GetTradeVolume - returns trade volumes
func (api *Kraken) GetTradeVolume(needFeeInfo bool, pairs ...string) (*TradeVolumeResponse, error) {
	data := url.Values{}
	if len(pairs) == 0 {
		return nil, fmt.Errorf("`pairs` is required")
	}
	if needFeeInfo {
		data.Set("fee-info", "true")
	}
	data.Set("pair", strings.Join(pairs, ","))

	resp, err := api.request("TradeVolume", true, data, &TradeVolumeResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*TradeVolumeResponse), nil
}

// AddOrder - method sends order to exchange
func (api *Kraken) AddOrder(pair string, side string, orderType string, volume float64, args map[string]interface{}) (*AddOrderResponse, error) {
	data := url.Values{
		"pair":      {pair},
		"volume":    {strconv.FormatFloat(volume, 'f', 8, 64)},
		"type":      {side},
		"ordertype": {orderType},
	}
	if args != nil {
		for key, value := range args {
			switch value.(type) {
			case string:
				data.Set(key, value.(string))
			case int64:
				data.Set(key, strconv.FormatInt(value.(int64), 10))
			case float64:
				data.Set(key, strconv.FormatFloat(value.(float64), 'f', 8, 64))
			case bool:
				data.Set(key, strconv.FormatBool(value.(bool)))
			default:
				log.Printf("[WARNING] Unknown value type %v for key %s", value, key)
			}
		}
	}

	resp, err := api.request("AddOrder", true, data, &AddOrderResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*AddOrderResponse), nil
}

// Cancel - method cancels order
func (api *Kraken) Cancel(orderID string) (*CancelResponse, error) {
	data := url.Values{
		"txid": {orderID},
	}
	resp, err := api.request("CancelOrder", true, data, &CancelResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*CancelResponse), nil
}
