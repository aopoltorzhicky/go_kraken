package futures

const (
	// FuturesAPIUrl - Kraken Futures API Endpoint
	FuturesAPIUrl = "https://futures.kraken.com/derivatives/api/v3"
)

const (
	// Order Types
	OrderTypeLimit         = "lmt"
	OrderTypePostOnlyLimit = "post"
	OrderTypeMarket        = "mkt"
	OrderTypeStop          = "stp"
	OrderTypeTakeProfit    = "take_profit"
	OrderTypeIOC           = "ioc"
	OrderTypeTrailingStop  = "trailing_stop"

	// Order Sides
	OrderSideBuy  = "buy"
	OrderSideSell = "sell"
)
