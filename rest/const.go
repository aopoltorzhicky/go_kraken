package rest

const (
	// APIUrl - Kraken API Endpoint
	APIUrl = "https://api.kraken.com"
	// APIVersion - Kraken API Version Number
	APIVersion = "0"
)

// Interval values
const (
	Interval1m  = 1
	Interval5m  = 5
	Interval15m = 15
	Interval30m = 30
	Interval1h  = 60
	Interval4h  = 240
	Interval1d  = 1440
	Interval7d  = 10080
	Interval1M  = 21600
)

// Order Sides
const (
	TradeBuy  = "b"
	TradeSell = "s"
	Buy       = "buy"
	Sell      = "sell"
)

// Order types
const (
	TradeLimit  = "l"
	TradeMarket = "m"
	Limit       = "limit"
	Market      = "market"
)

// Assets
const (
	ADA  = "ADA"
	ATOM = "ATOM"
	BCH  = "BCH"
	BSV  = "BSV"
	DASH = "DASH"
	EOS  = "EOS"
	GNO  = "GNO"
	KFEE = "KFEE"
	QTUM = "QTUM"
	USDT = "USDT"
	XDAO = "XDAO"
	XETC = "XETC"
	XETH = "XETH"
	XICN = "XICN"
	XLTC = "XLTC"
	XMLN = "XMLN"
	XNMC = "XNMC"
	XREP = "XREP"
	XTZ  = "XTZ"
	XXBT = "XXBT"
	XXDG = "XXDG"
	XXLM = "XXLM"
	XXMR = "XXMR"
	XXRP = "XXRP"
	XXVN = "XXVN"
	XZEC = "XZEC"
	ZCAD = "ZCAD"
	ZEUR = "ZEUR"
	ZGBP = "ZGBP"
	ZJPY = "ZJPY"
	ZKRW = "ZKRW"
	ZUSD = "ZUSD"
)

// Trade types
const (
	TradeTypeAll             = "all"
	TradeTypeAnyPosition     = "any position"
	TradeTypeClosedPosition  = "closed position"
	TradeTypeClosingPosition = "closing position"
	TradeTypeNoPosition      = "no position"
)

// Ledger types
const (
	LedgerTypeAll        = "all"
	LedgerTypeDeposit    = "deposit"
	LedgerTypeWithdrawal = "withdrawal"
	LedgerTypeTrade      = "trade"
	LedgerTypeMargin     = "margin"
	LedgerTypeRollover   = "rollover"
)

// OrderTypes for AddOrder
const (
	OTMarket              = "market"
	OTLimit               = "limit"                  // (price = limit price)
	OTStopLoss            = "stop-loss"              // (price = stop loss price)
	OTTakeProfi           = "take-profit"            // (price = take profit price)
	OTStopLossProfit      = "stop-loss-profit"       // (price = stop loss price, price2 = take profit price)
	OTStopLossProfitLimit = "stop-loss-profit-limit" // (price = stop loss price, price2 = take profit price)
	OTStopLossLimit       = "stop-loss-limit"        // (price = stop loss trigger price, price2 = triggered limit price)
	OTTakeProfitLimit     = "take-profit-limit"      // (price = take profit trigger price, price2 = triggered limit price)
	OTTrailingStop        = "trailing-stop"          // (price = trailing stop offset)
	OTTrailingStopLimit   = "trailing-stop-limit"    // (price = trailing stop offset, price2 = triggered limit offset)
	OTStopLossAndLimit    = "stop-loss-and-limit"    // (price = stop loss price, price2 = limit price)
	OTSettlePosition      = "settle-position"
)

// OrderStatuses
const (
	StatusPending   = "pending"
	StatusOpen      = "open"
	StatusClosed    = "closed"
	StatusCancelled = "canceled"
	StatusExpired   = "expired"
)
