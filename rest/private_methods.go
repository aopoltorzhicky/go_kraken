package rest

import (
	"errors"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

// GetAccountBalances - methods returns account balances
func (api *Kraken) GetAccountBalances() (map[string]decimal.Decimal, error) {
	response := make(map[string]decimal.Decimal)
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

	switch {
	case len(txIDs) > 50:
		return nil, errors.New("maximum count of requested orders is 50")
	case len(txIDs) == 0:
		return nil, errors.New("txIDs is required")
	default:
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

// GetDepositMethods - returns deposit methods
func (api *Kraken) GetDepositMethods(assets ...string) ([]DepositMethods, error) {
	data := url.Values{}
	if len(assets) > 0 {
		data.Add("asset", strings.Join(assets, ","))
	} else {
		data = nil
	}

	response := make([]DepositMethods, 0)
	if err := api.request("DepositMethods", true, data, &response); err != nil {
		return response, err
	}
	return response, nil
}

// GetDepositStatus - returns deposit status
func (api *Kraken) GetDepositStatus(method string, assets ...string) ([]DepositStatuses, error) {
	data := url.Values{}
	if len(assets) > 0 {
		data.Add("asset", strings.Join(assets, ","))
	}

	if len(method) > 0 {
		data.Add("method", method)
	}
	response := make([]DepositStatuses, 0)
	if err := api.request("DepositStatus", true, data, &response); err != nil {
		return response, err
	}
	return response, nil
}

// WithdrawFunds - returns withdrawal response
func (api *Kraken) WithdrawFunds(asset string, key string, amount float64) (response WithdrawFunds, err error) {
	data := url.Values{
		"asset":  {asset},
		"key":    {key},
		"amount": {strconv.FormatFloat(amount, 'f', 8, 64)},
	}

	if err = api.request("Withdraw", true, data, &response); err != nil {
		return response, err
	}
	return response, nil
}

// GetWithdrawStatus - returns withdrawal statuses
func (api *Kraken) GetWithdrawStatus(asset string, method string) ([]WithdrawStatus, error) {
	data := url.Values{}

	if len(asset) > 0 {
		data.Add("asset", asset)
	}

	if len(method) > 0 {
		data.Add("method", method)
	}

	response := make([]WithdrawStatus, 0)
	if err := api.request("WithdrawStatus", true, data, &response); err != nil {
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
		return nil, errors.New("txIDs is required")
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
		return nil, errors.New("txIDs is required")
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
		return nil, errors.New("`ledgerIds` is required")
	}
	if len(ledgerIds) > 20 {
		return nil, errors.New("maximum count of requested ledgers is 20")
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
		return response, errors.New("`pairs` is required")
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
func (api *Kraken) AddOrder(pair string, side string, orderType string, volume float64, args map[string]interface{}) (response AddOrderResponse, err error) {
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

	err = api.request("AddOrder", true, data, &response)
	return
}

// EditOrder - method edits an existing order in the exchange
func (api *Kraken) EditOrder(orderId string, pair string, args map[string]interface{}) (response EditOrderResponse, err error) {
	data := url.Values{
		"txid": {orderId},
		"pair": {pair},
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

	err = api.request("EditOrder", true, data, &response)
	return
}

// Cancel - method cancels order
func (api *Kraken) Cancel(orderID string) (response CancelResponse, err error) {
	data := url.Values{
		"txid": {orderID},
	}
	err = api.request("CancelOrder", true, data, &response)
	return
}

// GetWebSocketsToken - WebSockets authentication
func (api *Kraken) GetWebSocketsToken() (response GetWebSocketTokenResponse, err error) {
	err = api.request("GetWebSocketsToken", true, nil, &response)
	return
}
