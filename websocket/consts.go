package websocket

// URLs
const (
	prodBaseURL    = "wss://ws.kraken.com"
	sandboxBaseURL = "wss://ws-sandbox.kraken.com"
)

// Available channels
const (
	ChanBook    = "book"
	ChanTrades  = "trade"
	ChanTicker  = "ticker"
	ChanCandles = "ohlc"
	ChanSpread  = "spread"
	ChanAll     = "*"
)

// Events
const (
	EventSubscribe          = "subscribe"
	EventUnsubscribe        = "unsubscribe"
	EventPing               = "ping"
	EventPong               = "pong"
	EventSystemStatus       = "systemStatus"
	EventSubscriptionStatus = "subscriptionStatus"
	EventHeartbeat          = "heartbeat"
)

// Intervals
const (
	Interval1    = 1
	Interval5    = 5
	Interval15   = 15
	Interval30   = 30
	Interval60   = 60
	Interal240   = 240
	Interal1440  = 1440
	Interal10080 = 10080
	Interal21600 = 21600
)

// Depth
const (
	Depth10   = 10
	Depth25   = 25
	Depth100  = 100
	Depth500  = 500
	Depth1000 = 1000
)

// Subscription Statuses
const (
	SubscriptionStatusError        = "error"
	SubscriptionStatusSubscribed   = "subscribed"
	SubscriptionStatusUnsubscribed = "unsubscribed"
)

// Trade sides
const (
	Buy  = "buy"
	Sell = "sell"
)

// Order types
const (
	Market = "market"
	Limit  = "limit"
)
