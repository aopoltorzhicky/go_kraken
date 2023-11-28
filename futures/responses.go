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

type OrderStatusResponse struct {
	Orders []struct {
		Error        string `json:"error"`
		Order        Order  `json:"order"`
		Status       string `json:"status"`
		UpdateReason string `json:"updateReason"`
	} // `json:"orders"`
}

// OrderBookResponse - wrapper struct for the full response
type OrderBookResponse struct {
	Result    string    `json:"result"`
	OrderBook OrderBook `json:"orderBook"`
}

type SendOrderResponse struct {
	Result     string     `json:"result"`
	SendStatus SendStatus `json:"sendStatus"`
	ServerTime string     `json:"serverTime"`
}

type SendStatus struct {
	OrderID      string       `json:"order_id"`
	Status       string       `json:"status"`
	ReceivedTime string       `json:"receivedTime"`
	OrderEvents  []OrderEvent `json:"orderEvents"`
}

type OrderEvent struct {
	Order                *Order   `json:"order,omitempty"`
	ReducedQuantity      *float64 `json:"reducedQuantity,omitempty"`
	Type                 string   `json:"type"`
	ExecutionID          *string  `json:"executionId,omitempty"`
	Price                *float64 `json:"price,omitempty"`
	Amount               *float64 `json:"amount,omitempty"`
	OrderPriorEdit       *Order   `json:"orderPriorEdit,omitempty"`
	OrderPriorExecution  *Order   `json:"orderPriorExecution,omitempty"`
	TakerReducedQuantity *float64 `json:"takerReducedQuantity,omitempty"`
}

type Order struct {
	CliOrdId            *string  `json:"cliOrdId,omitempty"`
	Filled              *float64 `json:"filled,omitempty"`
	LastUpdateTimestamp string   `json:"lastUpdateTimestamp"`
	LimitPrice          *float64 `json:"limitPrice,omitempty"`
	OrderId             string   `json:"orderId"`
	Quantity            *float64 `json:"quantity,omitempty"`
	ReduceOnly          bool     `json:"reduceOnly"`
	Side                string   `json:"side"`
	Symbol              string   `json:"symbol"`
	Timestamp           string   `json:"timestamp"`
	Type                string   `json:"type"`
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
