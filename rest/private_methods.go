package rest

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

// GetAccountBalances - methods returns account balances
func (api *Kraken) GetAccountBalances() (BalanceResponse, error) {
	response := BalanceResponse{}
	if err := api.request("Balance", true, nil, &response); err != nil {
		return response, err
	}
	return response, nil
}

// GetTradeBalance - returns tradable balances info
func (api *Kraken) GetTradeBalance(baseAsset string) (TradeBalanceResponse, error) {
	data := url.Values{}
	if baseAsset != "" {
		data.Set("asset", baseAsset)
	}

	response := TradeBalanceResponse{}
	if err := api.request("TradeBalance", true, data, &response); err != nil {
		return response, err
	}
	return response, nil
}

// GetOpenOrders - returns account open order
func (api *Kraken) GetOpenOrders(needTrades bool, userRef string) (OpenOrdersResponse, error) {
	data := url.Values{}
	if needTrades {
		data.Set("trades", "true")
	}
	if userRef != "" {
		data.Set("userref", userRef)
	}

	response := OpenOrdersResponse{}
	if err := api.request("OpenOrders", true, data, &response); err != nil {
		return response, err
	}
	return response, nil
}

// GetClosedOrders - returns account closed order
func (api *Kraken) GetClosedOrders(needTrades bool, userRef string, start int64, end int64) (ClosedOrdersResponse, error) {
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

	response := ClosedOrdersResponse{}
	if err := api.request("ClosedOrders", true, data, &response); err != nil {
		return response, err
	}
	return response, nil
}

// QueryOrders - returns account's order by IDs
func (api *Kraken) QueryOrders(needTrades bool, userRef string, txIDs ...string) (map[string]OrderInfo, error) {
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

	response := make(map[string]OrderInfo)
	if err := api.request("QueryOrders", true, data, &response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetTradesHistory - returns account's trade history
func (api *Kraken) GetTradesHistory(tradeType string, needTrades bool, start int64, end int64) (TradesHistoryResponse, error) {
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
	response := TradesHistoryResponse{}
	if err := api.request("TradesHistory", true, data, &response); err != nil {
		return response, err
	}
	return response, nil
}

// QueryTrades - returns trades by IDs
func (api *Kraken) QueryTrades(trades bool, txIDs ...string) (map[string]PrivateTrade, error) {
	data := url.Values{}
	if trades {
		data.Set("trades", "true")
	}
	if len(txIDs) == 0 {
		return nil, fmt.Errorf("txIDs is required")
	}
	data.Set("txid", strings.Join(txIDs, ","))

	response := make(map[string]PrivateTrade)
	if err := api.request("QueryTrades", true, data, &response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetOpenPositions - returns list of open positions
func (api *Kraken) GetOpenPositions(docalcs bool, txIDs ...string) (map[string]Position, error) {
	data := url.Values{}
	if docalcs {
		data.Set("docalcs", "true")
	}
	if len(txIDs) == 0 {
		return nil, fmt.Errorf("txIDs is required")
	}
	data.Set("txid", strings.Join(txIDs, ","))

	response := make(map[string]Position)
	if err := api.request("OpenPositions", true, data, &response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetLedgersInfo - returns ledgers info
func (api *Kraken) GetLedgersInfo(ledgerType string, start int64, end int64, assets ...string) (LedgerInfoResponse, error) {
	response := LedgerInfoResponse{}
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
	if len(assets) > 0 {
		data.Set("assets", strings.Join(assets, ","))
	}

	if err := api.request("Ledgers", true, data, &response); err != nil {
		return response, err
	}
	return response, nil
}

// QueryLedgers - get ledgers by ID
func (api *Kraken) QueryLedgers(ledgerIds ...string) (map[string]Ledger, error) {
	data := url.Values{}
	if len(ledgerIds) == 0 {
		return nil, fmt.Errorf("`ledgerIds` is required")
	}
	if len(ledgerIds) > 20 {
		return nil, fmt.Errorf("Maximum count of requested ledgers is 20")
	}
	data.Set("id", strings.Join(ledgerIds, ","))

	response := make(map[string]Ledger)
	if err := api.request("QueryLedgers", true, data, &response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetTradeVolume - returns trade volumes
func (api *Kraken) GetTradeVolume(needFeeInfo bool, pairs ...string) (TradeVolumeResponse, error) {
	response := TradeVolumeResponse{}
	data := url.Values{}
	if len(pairs) == 0 {
		return response, fmt.Errorf("`pairs` is required")
	}
	if needFeeInfo {
		data.Set("fee-info", "true")
	}
	data.Set("pair", strings.Join(pairs, ","))

	if err := api.request("TradeVolume", true, data, &response); err != nil {
		return response, err
	}
	return response, nil
}

// AddOrder - method sends order to exchange
func (api *Kraken) AddOrder(pair string, side string, orderType string, volume float64, args map[string]interface{}) (AddOrderResponse, error) {
	data := url.Values{
		"pair":      {pair},
		"volume":    {strconv.FormatFloat(volume, 'f', 8, 64)},
		"type":      {side},
		"ordertype": {orderType},
	}
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

	response := AddOrderResponse{}
	if err := api.request("AddOrder", true, data, &response); err != nil {
		return response, err
	}
	return response, nil
}

// Cancel - method cancels order
func (api *Kraken) Cancel(orderID string) (CancelResponse, error) {
	data := url.Values{
		"txid": {orderID},
	}
	response := CancelResponse{}
	if err := api.request("CancelOrder", true, data, &response); err != nil {
		return response, err
	}
	return response, nil
}

// GetWebSocketsToken - WebSockets authentication
func (api *Kraken) GetWebSocketsToken() (GetWebSocketTokenResponse, error) {
	response := GetWebSocketTokenResponse{}
	if err := api.request("GetWebSocketsToken", true, nil, &response); err != nil {
		return response, err
	}
	return response, nil
}
