package futures

import "net/url"

// OrderBook - Fetches the order book for a specified futures symbol.
func (api *KrakenFutures) OrderBook(symbol string) (OrderBook, error) {
	var data url.Values
	if symbol != "" {
		data = url.Values{"symbol": {symbol}}
	}
	var response OrderBookResponse
	if err := api.request("GET", "orderbook", false, data, &response); err != nil {
		return response.OrderBook, err
	}
	return response.OrderBook, nil
}

// Tickers - Gets current market data for all listed futures contracts and indices.
func (api *KrakenFutures) Tickers() (Ticker, error) {
	var response Ticker
	if err := api.request("GET", "tickers", false, nil, &response); err != nil {
		return response, err
	}
	return response, nil
}
