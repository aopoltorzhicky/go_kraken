package futures

type KrakenResponse struct {
	Error  []string    `json:"error"`
	Result interface{} `json:"result"`
}

// OrderBookItem - one price level in orderbook
type OrderBookItem []float64

// OrderBook - struct of order book levels
type OrderBook struct {
	Asks []OrderBookItem `json:"asks"`
	Bids []OrderBookItem `json:"bids"`
}

// OrderBookResponse - wrapper struct for the full response
type OrderBookResponse struct {
	Result    string    `json:"result"`
	OrderBook OrderBook `json:"orderBook"`
}

// Ticker - struct of ticker response
type Ticker struct {
	Ask                   float64 `json:"ask"`
	AskSize               float64 `json:"askSize"`
	Bid                   float64 `json:"bid"`
	BidSize               float64 `json:"bidSize"`
	Change24h             float64 `json:"change24h"`
	FundingRate           float64 `json:"fundingRate"`
	FundingRatePrediction float64 `json:"fundingRatePrediction"`
	IndexPrice            float64 `json:"indexPrice"`
	Last                  float64 `json:"last"`
	LastSize              float64 `json:"lastSize"`
	LastTime              string  `json:"lastTime"`
	MarkPrice             float64 `json:"markPrice"`
	Open24h               float64 `json:"open24h"`
	OpenInterest          float64 `json:"openInterest"`
	Pair                  string  `json:"pair"`
	PostOnly              bool    `json:"postOnly"`
	Suspended             bool    `json:"suspended"`
	Symbol                string  `json:"symbol"`
	Tag                   string  `json:"tag"`
	Vol24h                float64 `json:"vol24h"`
}
