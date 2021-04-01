package websocket

// URLs
const (
	ProdBaseURL    = "wss://ws.kraken.com"
	AuthBaseURL    = "wss://ws-auth.kraken.com"
	SandboxBaseURL = "wss://beta-ws.kraken.com"
)

// Available channels
const (
	ChanBook       = "book"
	ChanTrades     = "trade"
	ChanTicker     = "ticker"
	ChanCandles    = "ohlc"
	ChanSpread     = "spread"
	ChanOpenOrders = "openOrders"
	ChanOwnTrades  = "ownTrades"
	ChanAll        = "*"
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
	EventAddOrder           = "addOrder"
	EventAddOrderStatus     = "addOrderStatus"
	EventCancelOrder        = "cancelOrder"
	EventCancelOrderStatus  = "cancelOrderStatus"
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

// Pairs
const (
	ADACAD  = "ADA/CAD"
	ADAETH  = "ADA/ETH"
	ADAEUR  = "ADA/EUR"
	ADAUSD  = "ADA/USD"
	ADABTC  = "ADA/XBT"
	BCHEUR  = "BCH/EUR"
	BCHUSD  = "BCH/USD"
	BCHBTC  = "BCH/XBT"
	BSVEUR  = "BSV/EUR"
	BSVUSD  = "BSV/USD"
	BSVBTC  = "BSV/XBT"
	BTCEUR  = "XBT/EUR"
	BTCUSD  = "XBT/USD"
	BTCCAD  = "XBT/CAD"
	BTCJPY  = "XBT/JPY"
	BTCGBP  = "XBT/GBP"
	DASHEUR = "DASH/EUR"
	DASHUSD = "DASH/USD"
	DASHBTC = "DASH/XBT"
	DOGEBTC = "XDG/XBT"
    DOTEUR  = "DOT/EUR"
    DOTUSD  = "DOT/USD"
	EOSETH  = "EOS/ETH"
	EOSEUR  = "EOS/EUR"
	EOSUSD  = "EOS/USD"
	EOSBTC  = "EOS/XBT"
	ETCETH  = "ETC/ETH"
	ETCEUR  = "ETC/EUR"
	ETCUSD  = "ETC/USD"
	ETCBTC  = "ETC/XBT"
	ETHCAD  = "ETH/CAD"
	ETHEUR  = "ETH/EUR"
	ETHUSD  = "ETH/USD"
	ETHBTC  = "ETH/XBT"
	ETHJPY  = "ETH/JPY"
	ETHGBP  = "ETH/GBP"
	GNOETH  = "GNO/ETH"
	GNOEUR  = "GNO/EUR"
	GNOUSD  = "GNO/USD"
	GNOBTC  = "GNO/XBT"
	LTCEUR  = "LTC/EUR"
	LTCUSD  = "LTC/USD"
	LTCBTC  = "LTC/XBT"
	MLNETH  = "MLN/ETH"
	MLNBTC  = "MLN/XBT"
	QTUMCAD = "QTUM/CAD"
	QTUMETH = "QTUM/ETH"
	QTUMEUR = "QTUM/EUR"
	QTUMUSD = "QTUM/USD"
	QTUMBTC = "QTUM/XBT"
	REPETH  = "REP/ETH"
	REPEUR  = "REP/EUR"
	REPUSD  = "REP/USD"
	REPBTC  = "REP/XBT"
	USDTUSD = "USDT/USD"
	XLMEUR  = "XLM/EUR"
	XLMUSD  = "XLM/USD"
	XLMBTC  = "XLM/XBT"
	XMREUR  = "XMR/EUR"
	XMRUSD  = "XMR/USD"
	XMRBTC  = "XMR/XBT"
	XRPCAD  = "XRP/CAD"
	XRPEUR  = "XRP/EUR"
	XRPJPY  = "XRP/JPY"
	XRPUSD  = "XRP/USD"
	XRPBTC  = "XRP/XBT"
	XTZCAD  = "XTZ/CAD"
	XTZETH  = "XTZ/ETH"
	XTZEUR  = "XTZ/EUR"
	XTZUSD  = "XTZ/USD"
	XTZBTC  = "XTZ/XBT"
	ZECEUR  = "ZEC/EUR"
	ZECJPY  = "ZEC/JPY"
	ZECUSD  = "ZEC/USD"
)

// Statuses
const (
	StatusOK    = "ok"
	StatusError = "error"
)
