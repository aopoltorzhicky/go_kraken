package rest

// KrakenResponse - template of Kraken API response
type KrakenResponse struct {
	Error  []string    `json:"error"`
	Result interface{} `json:"result"`
}

// TimeResponse - Result of Time request
type TimeResponse struct {
	Unixtime int64  `json:"unixtime"`
	Rfc1123  string `json:"rfc1123"`
}

// AssetResponse - Result of Assets request
type AssetResponse struct {
	ADA  Asset
	BCH  Asset
	BSV  Asset
	DASH Asset
	EOS  Asset
	GNO  Asset
	KFEE Asset
	QTUM Asset
	USDT Asset
	XDAO Asset
	XETC Asset
	XETH Asset
	XICN Asset
	XLTC Asset
	XMLN Asset
	XNMC Asset
	XREP Asset
	XXBT Asset
	XXDG Asset
	XXLM Asset
	XXMR Asset
	XXRP Asset
	XXTZ Asset
	XXVN Asset
	XZEC Asset
	ZCAD Asset
	ZEUR Asset
	ZGBP Asset
	ZJPY Asset
	ZKRW Asset
	ZUSD Asset
}

// Asset - asset information
type Asset struct {
	AlternateName   string `json:"altname"`
	AssetClass      string `json:"aclass"`
	Decimals        int    `json:"decimals"`
	DisplayDecimals int    `json:"display_decimals"`
}

// AssetPairsResponse - struct with asset pair informations
type AssetPairsResponse struct {
	ADACAD   AssetPair
	ADAETH   AssetPair
	ADAEUR   AssetPair
	ADAUSD   AssetPair
	ADAXBT   AssetPair
	BCHEUR   AssetPair
	BCHUSD   AssetPair
	BCHXBT   AssetPair
	DASHEUR  AssetPair
	DASHUSD  AssetPair
	DASHXBT  AssetPair
	EOSETH   AssetPair
	EOSEUR   AssetPair
	EOSUSD   AssetPair
	EOSXBT   AssetPair
	GNOETH   AssetPair
	GNOEUR   AssetPair
	GNOUSD   AssetPair
	GNOXBT   AssetPair
	QTUMCAD  AssetPair
	QTUMETH  AssetPair
	QTUMEUR  AssetPair
	QTUMUSD  AssetPair
	QTUMXBT  AssetPair
	USDTZUSD AssetPair
	XETCXETH AssetPair
	XETCXXBT AssetPair
	XETCZEUR AssetPair
	XETCZUSD AssetPair
	XETHXXBT AssetPair
	XETHZCAD AssetPair
	XETHZEUR AssetPair
	XETHZGBP AssetPair
	XETHZJPY AssetPair
	XETHZUSD AssetPair
	XICNXETH AssetPair
	XICNXXBT AssetPair
	XLTCXXBT AssetPair
	XLTCZEUR AssetPair
	XLTCZUSD AssetPair
	XMLNXETH AssetPair
	XMLNXXBT AssetPair
	XREPXETH AssetPair
	XREPXXBT AssetPair
	XREPZEUR AssetPair
	XREPZUSD AssetPair
	XTZCAD   AssetPair
	XTZETH   AssetPair
	XTZEUR   AssetPair
	XTZUSD   AssetPair
	XTZXBT   AssetPair
	XXBTZCAD AssetPair
	XXBTZEUR AssetPair
	XXBTZGBP AssetPair
	XXBTZJPY AssetPair
	XXBTZUSD AssetPair
	XXDGXXBT AssetPair
	XXLMXXBT AssetPair
	XXLMZEUR AssetPair
	XXLMZUSD AssetPair
	XXMRXXBT AssetPair
	XXMRZEUR AssetPair
	XXMRZUSD AssetPair
	XXRPXXBT AssetPair
	XXRPZCAD AssetPair
	XXRPZEUR AssetPair
	XXRPZJPY AssetPair
	XXRPZUSD AssetPair
	XZECXXBT AssetPair
	XZECZEUR AssetPair
	XZECZUSD AssetPair
}

// AssetPair - asset pair information
type AssetPair struct {
	Altname           string      `json:"altname"`
	AssetClassBase    string      `json:"aclass_base"`
	Base              string      `json:"base"`
	AssetClassQuote   string      `json:"aclass_quote"`
	Quote             string      `json:"quote"`
	Lot               string      `json:"lot"`
	PairDecimals      int         `json:"pair_decimals"`
	LotDecimals       int         `json:"lot_decimals"`
	LotMultiplier     int         `json:"lot_multiplier"`
	LeverageBuy       []float64   `json:"leverage_buy"`
	LeverageSell      []float64   `json:"leverage_sell"`
	Fees              [][]float64 `json:"fees"`
	FeesMaker         [][]float64 `json:"fees_maker"`
	FeeVolumeCurrency string      `json:"fee_volume_currency"`
	MarginCall        int         `json:"margin_call"`
	MarginStop        int         `json:"margin_stop"`
	WSName            string      `json:"wsname"`
}

// TickerResponse - all pairs in ticker response
type TickerResponse struct {
	ADACAD   Ticker
	ADAETH   Ticker
	ADAEUR   Ticker
	ADAUSD   Ticker
	ADAXBT   Ticker
	BCHEUR   Ticker
	BCHUSD   Ticker
	BCHXBT   Ticker
	DASHEUR  Ticker
	DASHUSD  Ticker
	DASHXBT  Ticker
	EOSETH   Ticker
	EOSEUR   Ticker
	EOSUSD   Ticker
	EOSXBT   Ticker
	GNOETH   Ticker
	GNOEUR   Ticker
	GNOUSD   Ticker
	GNOXBT   Ticker
	QTUMCAD  Ticker
	QTUMETH  Ticker
	QTUMEUR  Ticker
	QTUMUSD  Ticker
	QTUMXBT  Ticker
	USDTZUSD Ticker
	XETCXETH Ticker
	XETCXXBT Ticker
	XETCZEUR Ticker
	XETCZUSD Ticker
	XETHXXBT Ticker
	XETHZCAD Ticker
	XETHZEUR Ticker
	XETHZGBP Ticker
	XETHZJPY Ticker
	XETHZUSD Ticker
	XICNXETH Ticker
	XICNXXBT Ticker
	XLTCXXBT Ticker
	XLTCZEUR Ticker
	XLTCZUSD Ticker
	XMLNXETH Ticker
	XMLNXXBT Ticker
	XREPXETH Ticker
	XREPXXBT Ticker
	XREPZEUR Ticker
	XREPZUSD Ticker
	XXBTZCAD Ticker
	XXBTZEUR Ticker
	XXBTZGBP Ticker
	XXBTZJPY Ticker
	XXBTZUSD Ticker
	XXDGXXBT Ticker
	XXLMXXBT Ticker
	XXLMZEUR Ticker
	XXLMZUSD Ticker
	XXMRXXBT Ticker
	XXMRZEUR Ticker
	XXMRZUSD Ticker
	XXRPXXBT Ticker
	XXRPZCAD Ticker
	XXRPZEUR Ticker
	XXRPZJPY Ticker
	XXRPZUSD Ticker
	XTZCAD   Ticker
	XTZETH   Ticker
	XTZEUR   Ticker
	XTZUSD   Ticker
	XTZXBT   Ticker
	XZECXXBT Ticker
	XZECZEUR Ticker
	XZECZUSD Ticker
}

// Level - ticker structure for Ask and Bid
type Level struct {
	Price          float64 `json:",string"`
	WholeLotVolume float64 `json:",string"`
	Volume         float64 `json:",string"`
}

// TimeLevel - ticker structure for Volume, VolumeAveragePrice, Low, High
type TimeLevel struct {
	Today       float64 `json:",string"`
	Last24Hours float64 `json:",string"`
}

// CloseLevel - ticker structure for Close
type CloseLevel struct {
	Price     float64 `json:",string"`
	LotVolume float64 `json:",string"`
}

// Ticker - struct of ticker response
type Ticker struct {
	Ask                Level      `json:"a"`
	Bid                Level      `json:"b"`
	Close              CloseLevel `json:"c"`
	Volume             TimeLevel  `json:"v"`
	VolumeAveragePrice TimeLevel  `json:"p"`
	Trades             TimeLevel  `json:"t"`
	Low                TimeLevel  `json:"l"`
	High               TimeLevel  `json:"h"`
	OpeningPrice       float64    `json:"o,string"`
}

// Candle - OHLC item
type Candle struct {
	Time      int64
	Open      float64 `json:",string"`
	High      float64 `json:",string"`
	Low       float64 `json:",string"`
	Close     float64 `json:",string"`
	VolumeWAP float64 `json:",string"`
	Volume    float64 `json:",string"`
	Count     int64
}

// OHLCResponse - response of OHLC request
type OHLCResponse struct {
	Last     int64 `json:"last"`
	ADACAD   []Candle
	ADAETH   []Candle
	ADAEUR   []Candle
	ADAUSD   []Candle
	ADAXBT   []Candle
	BCHEUR   []Candle
	BCHUSD   []Candle
	BCHXBT   []Candle
	DASHEUR  []Candle
	DASHUSD  []Candle
	DASHXBT  []Candle
	EOSETH   []Candle
	EOSEUR   []Candle
	EOSUSD   []Candle
	EOSXBT   []Candle
	GNOETH   []Candle
	GNOEUR   []Candle
	GNOUSD   []Candle
	GNOXBT   []Candle
	QTUMCAD  []Candle
	QTUMETH  []Candle
	QTUMEUR  []Candle
	QTUMUSD  []Candle
	QTUMXBT  []Candle
	USDTZUSD []Candle
	XETCXETH []Candle
	XETCXXBT []Candle
	XETCZEUR []Candle
	XETCZUSD []Candle
	XETHXXBT []Candle
	XETHZCAD []Candle
	XETHZEUR []Candle
	XETHZGBP []Candle
	XETHZJPY []Candle
	XETHZUSD []Candle
	XICNXETH []Candle
	XICNXXBT []Candle
	XLTCXXBT []Candle
	XLTCZEUR []Candle
	XLTCZUSD []Candle
	XMLNXETH []Candle
	XMLNXXBT []Candle
	XREPXETH []Candle
	XREPXXBT []Candle
	XREPZEUR []Candle
	XREPZUSD []Candle
	XXBTZCAD []Candle
	XXBTZEUR []Candle
	XXBTZGBP []Candle
	XXBTZJPY []Candle
	XXBTZUSD []Candle
	XXDGXXBT []Candle
	XXLMXXBT []Candle
	XXLMZEUR []Candle
	XXLMZUSD []Candle
	XXMRXXBT []Candle
	XXMRZEUR []Candle
	XXMRZUSD []Candle
	XXRPXXBT []Candle
	XXRPZCAD []Candle
	XXRPZEUR []Candle
	XXRPZJPY []Candle
	XXRPZUSD []Candle
	XTZCAD   []Candle
	XTZETH   []Candle
	XTZEUR   []Candle
	XTZUSD   []Candle
	XTZXBT   []Candle
	XZECXXBT []Candle
	XZECZEUR []Candle
	XZECZUSD []Candle
}

// OrderBookItem - one price level in orderbook
type OrderBookItem struct {
	Price     float64 `json:",string"`
	Volume    float64 `json:",string"`
	Timestamp int64
}

// OrderBook - struct of order book levels
type OrderBook struct {
	Asks []OrderBookItem `json:"asks"`
	Bids []OrderBookItem `json:"bids"`
}

// BookResponse - all pairs in book response
type BookResponse struct {
	ADACAD   OrderBook
	ADAETH   OrderBook
	ADAEUR   OrderBook
	ADAUSD   OrderBook
	ADAXBT   OrderBook
	BCHEUR   OrderBook
	BCHUSD   OrderBook
	BCHXBT   OrderBook
	DASHEUR  OrderBook
	DASHUSD  OrderBook
	DASHXBT  OrderBook
	EOSETH   OrderBook
	EOSEUR   OrderBook
	EOSUSD   OrderBook
	EOSXBT   OrderBook
	GNOETH   OrderBook
	GNOEUR   OrderBook
	GNOUSD   OrderBook
	GNOXBT   OrderBook
	QTUMCAD  OrderBook
	QTUMETH  OrderBook
	QTUMEUR  OrderBook
	QTUMUSD  OrderBook
	QTUMXBT  OrderBook
	USDTZUSD OrderBook
	XETCXETH OrderBook
	XETCXXBT OrderBook
	XETCZEUR OrderBook
	XETCZUSD OrderBook
	XETHXXBT OrderBook
	XETHZCAD OrderBook
	XETHZEUR OrderBook
	XETHZGBP OrderBook
	XETHZJPY OrderBook
	XETHZUSD OrderBook
	XICNXETH OrderBook
	XICNXXBT OrderBook
	XLTCXXBT OrderBook
	XLTCZEUR OrderBook
	XLTCZUSD OrderBook
	XMLNXETH OrderBook
	XMLNXXBT OrderBook
	XREPXETH OrderBook
	XREPXXBT OrderBook
	XREPZEUR OrderBook
	XREPZUSD OrderBook
	XXBTZCAD OrderBook
	XXBTZEUR OrderBook
	XXBTZGBP OrderBook
	XXBTZJPY OrderBook
	XXBTZUSD OrderBook
	XXDGXXBT OrderBook
	XXLMXXBT OrderBook
	XXLMZEUR OrderBook
	XXLMZUSD OrderBook
	XXMRXXBT OrderBook
	XXMRZEUR OrderBook
	XXMRZUSD OrderBook
	XXRPXXBT OrderBook
	XXRPZCAD OrderBook
	XXRPZEUR OrderBook
	XXRPZJPY OrderBook
	XXRPZUSD OrderBook
	XTZCAD   OrderBook
	XTZETH   OrderBook
	XTZEUR   OrderBook
	XTZUSD   OrderBook
	XTZXBT   OrderBook
	XZECXXBT OrderBook
	XZECZEUR OrderBook
	XZECZUSD OrderBook
}

// Trade - structure of public trades
type Trade struct {
	Price     float64 `json:",string"`
	Volume    float64 `json:",string"`
	Time      float64
	Side      string
	OrderType string
	Misc      string
}

// TradeResponse - all pairs in trade response
type TradeResponse struct {
	Last     float64 `json:"last"`
	ADACAD   []Trade
	ADAETH   []Trade
	ADAEUR   []Trade
	ADAUSD   []Trade
	ADAXBT   []Trade
	BCHEUR   []Trade
	BCHUSD   []Trade
	BCHXBT   []Trade
	DASHEUR  []Trade
	DASHUSD  []Trade
	DASHXBT  []Trade
	EOSETH   []Trade
	EOSEUR   []Trade
	EOSUSD   []Trade
	EOSXBT   []Trade
	GNOETH   []Trade
	GNOEUR   []Trade
	GNOUSD   []Trade
	GNOXBT   []Trade
	QTUMCAD  []Trade
	QTUMETH  []Trade
	QTUMEUR  []Trade
	QTUMUSD  []Trade
	QTUMXBT  []Trade
	USDTZUSD []Trade
	XETCXETH []Trade
	XETCXXBT []Trade
	XETCZEUR []Trade
	XETCZUSD []Trade
	XETHXXBT []Trade
	XETHZCAD []Trade
	XETHZEUR []Trade
	XETHZGBP []Trade
	XETHZJPY []Trade
	XETHZUSD []Trade
	XICNXETH []Trade
	XICNXXBT []Trade
	XLTCXXBT []Trade
	XLTCZEUR []Trade
	XLTCZUSD []Trade
	XMLNXETH []Trade
	XMLNXXBT []Trade
	XREPXETH []Trade
	XREPXXBT []Trade
	XREPZEUR []Trade
	XREPZUSD []Trade
	XXBTZCAD []Trade
	XXBTZEUR []Trade
	XXBTZGBP []Trade
	XXBTZJPY []Trade
	XXBTZUSD []Trade
	XXDGXXBT []Trade
	XXLMXXBT []Trade
	XXLMZEUR []Trade
	XXLMZUSD []Trade
	XXMRXXBT []Trade
	XXMRZEUR []Trade
	XXMRZUSD []Trade
	XXRPXXBT []Trade
	XXRPZCAD []Trade
	XXRPZEUR []Trade
	XXRPZJPY []Trade
	XXRPZUSD []Trade
	XTZCAD   []Trade
	XTZETH   []Trade
	XTZEUR   []Trade
	XTZUSD   []Trade
	XTZXBT   []Trade
	XZECXXBT []Trade
	XZECZEUR []Trade
	XZECZUSD []Trade
}

// Spread - structure of spread data
type Spread struct {
	Time float64
	Bid  float64 `json:",string"`
	Ask  float64 `json:",string"`
}

// SpreadResponse - response of spread request
type SpreadResponse struct {
	Last     float64 `json:"last"`
	ADACAD   []Spread
	ADAETH   []Spread
	ADAEUR   []Spread
	ADAUSD   []Spread
	ADAXBT   []Spread
	BCHEUR   []Spread
	BCHUSD   []Spread
	BCHXBT   []Spread
	DASHEUR  []Spread
	DASHUSD  []Spread
	DASHXBT  []Spread
	EOSETH   []Spread
	EOSEUR   []Spread
	EOSUSD   []Spread
	EOSXBT   []Spread
	GNOETH   []Spread
	GNOEUR   []Spread
	GNOUSD   []Spread
	GNOXBT   []Spread
	QTUMCAD  []Spread
	QTUMETH  []Spread
	QTUMEUR  []Spread
	QTUMUSD  []Spread
	QTUMXBT  []Spread
	USDTZUSD []Spread
	XETCXETH []Spread
	XETCXXBT []Spread
	XETCZEUR []Spread
	XETCZUSD []Spread
	XETHXXBT []Spread
	XETHZCAD []Spread
	XETHZEUR []Spread
	XETHZGBP []Spread
	XETHZJPY []Spread
	XETHZUSD []Spread
	XICNXETH []Spread
	XICNXXBT []Spread
	XLTCXXBT []Spread
	XLTCZEUR []Spread
	XLTCZUSD []Spread
	XMLNXETH []Spread
	XMLNXXBT []Spread
	XREPXETH []Spread
	XREPXXBT []Spread
	XREPZEUR []Spread
	XREPZUSD []Spread
	XXBTZCAD []Spread
	XXBTZEUR []Spread
	XXBTZGBP []Spread
	XXBTZJPY []Spread
	XXBTZUSD []Spread
	XXDGXXBT []Spread
	XXLMXXBT []Spread
	XXLMZEUR []Spread
	XXLMZUSD []Spread
	XXMRXXBT []Spread
	XXMRZEUR []Spread
	XXMRZUSD []Spread
	XXRPXXBT []Spread
	XXRPZCAD []Spread
	XXRPZEUR []Spread
	XXRPZJPY []Spread
	XXRPZUSD []Spread
	XTZCAD   []Spread
	XTZETH   []Spread
	XTZEUR   []Spread
	XTZUSD   []Spread
	XTZXBT   []Spread
	XZECXXBT []Spread
	XZECZEUR []Spread
	XZECZUSD []Spread
}

// BalanceResponse - response on account balance request
type BalanceResponse struct {
	ADA  float64 `json:",string"`
	BCH  float64 `json:",string"`
	BSV  float64 `json:",string"`
	DASH float64 `json:",string"`
	EOS  float64 `json:",string"`
	GNO  float64 `json:",string"`
	KFEE float64 `json:",string"`
	QTUM float64 `json:",string"`
	USDT float64 `json:",string"`
	XDAO float64 `json:",string"`
	XETC float64 `json:",string"`
	XETH float64 `json:",string"`
	XICN float64 `json:",string"`
	XLTC float64 `json:",string"`
	XMLN float64 `json:",string"`
	XNMC float64 `json:",string"`
	XREP float64 `json:",string"`
	XXBT float64 `json:",string"`
	XXDG float64 `json:",string"`
	XXLM float64 `json:",string"`
	XXMR float64 `json:",string"`
	XXRP float64 `json:",string"`
	XXTZ float64 `json:",string"`
	XXVN float64 `json:",string"`
	XZEC float64 `json:",string"`
	ZCAD float64 `json:",string"`
	ZEUR float64 `json:",string"`
	ZGBP float64 `json:",string"`
	ZJPY float64 `json:",string"`
	ZKRW float64 `json:",string"`
	ZUSD float64 `json:",string"`
}

// TradeBalanceResponse - response of get trade balance request
type TradeBalanceResponse struct {
	EquivalentBalance float64 `json:"eb,string"`
	TradeBalance      float64 `json:"tb,string"`
	OpenMargin        float64 `json:"m,string"`
	UnrealizedProfit  float64 `json:"n,string"`
	CostPositions     float64 `json:"c,string"`
	CurrentValue      float64 `json:"v,string"`
	Equity            float64 `json:"e,string"`
	FreeMargin        float64 `json:"mf,string"`
	MarginLevel       float64 `json:"ml,string"`
}

// OpenOrdersResponse - response on OpenOrders request
type OpenOrdersResponse struct {
	Orders map[string]OrderInfo `json:"open"`
}

// ClosedOrdersResponse - response on ClosedOrders request
type ClosedOrdersResponse struct {
	Count  int64                `json:"count"`
	Orders map[string]OrderInfo `json:"closed"`
}

// OrderInfo - structure contains order information
type OrderInfo struct {
	RefID           string           `json:"refid"`
	UserRef         string           `json:"userref"`
	Status          string           `json:"status"`
	OpenTimestamp   float64          `json:"opentm"`
	StartTimestamp  float64          `json:"starttm"`
	ExpireTimestamp float64          `json:"expiretm"`
	Description     OrderDescription `json:"descr"`
	Volume          float64          `json:"vol,string"`
	VolumeExecuted  float64          `json:"vol_exec,string"`
	Cost            float64          `json:"cost,string"`
	Fee             float64          `json:"fee,string"`
	AveragePrice    float64          `json:"price,string"`
	StopPrice       float64          `json:"stopprice,string"`
	LimitPrice      float64          `json:"limitprice,string"`
	Misc            string           `json:"misc"`
	Flags           string           `json:"oflags"`
}

// TradesHistoryResponse - respons on TradesHistory request
type TradesHistoryResponse struct {
	Trades map[string]PrivateTrade `json:"trades"`
	Count  int64                   `json:"count"`
}

// PrivateTrade - structure of account's trades
type PrivateTrade struct {
	OrderID              string   `json:"ordertxid"`
	PositionID           string   `json:"postxid"`
	Pair                 string   `json:"pair"`
	Time                 float64  `json:"time"`
	Side                 string   `json:"type"`
	OrderType            string   `json:"ordertype"`
	Price                float64  `json:"price,string"`
	Cost                 float64  `json:"cost,string"`
	Fee                  float64  `json:"fee,string"`
	Volume               float64  `json:"vol,string"`
	Margin               float64  `json:"margin,string"`
	Misc                 string   `json:"misc"`
	PositionStatus       string   `json:"posstatus,omitempty"`
	PositionAveragePrice float64  `json:"cprice,omitempty,string"`
	PositionCost         float64  `json:"ccost,omitempty,string"`
	PositionFee          float64  `json:"cfee,omitempty,string"`
	PositionVolume       float64  `json:"cvol,omitempty,string"`
	PositionMargin       float64  `json:"cmargin,omitempty,string"`
	PositionProfit       float64  `json:"net,omitempty,string"`
	PositionTrades       []string `json:"trades,omitempty"`
}

// Position - structure of account position
type Position struct {
	OrderID      string  `json:"ordertxid"`
	Status       string  `json:"posstatus"`
	Pair         string  `json:"pair"`
	Time         float64 `json:"time"`
	Side         string  `json:"type"`
	OrderType    string  `json:"ordertype"`
	Price        float64 `json:"price,string"`
	Cost         float64 `json:"cost,string"`
	Fee          float64 `json:"fee,string"`
	Volume       float64 `json:"vol,string"`
	VolumeClosed float64 `json:"vol_closed,string"`
	Margin       float64 `json:"margin,string"`
	Misc         string  `json:"misc"`
	Value        float64 `json:"value,omitempty,string"`
	Profit       float64 `json:"net,omitempty,string"`
	Terms        string  `json:"terms,omitempty"`
	RolloverTime float64 `json:"rollovertm,omitempty,string"`
	Flags        string  `json:"oflags"`
}

// LedgerInfoResponse - response on ledger request
type LedgerInfoResponse struct {
	Ledgers map[string]Ledger `json:"ledger"`
}

// Ledger - structure of account's ledger
type Ledger struct {
	RefID      string  `json:"refid"`
	Time       float64 `json:"time"`
	LedgerType string  `json:"type"`
	AssetClass string  `json:"aclass"`
	Asset      string  `json:"asset"`
	Amount     float64 `json:"amount,string"`
	Fee        float64 `json:"fee,string"`
	Balance    float64 `json:"balance,string"`
}

// TradeVolumeResponse - response on TradeVolume request
type TradeVolumeResponse struct {
	Currency  string          `json:"currency"`
	Volume    float64         `json:"volume,string"`
	Fees      map[string]Fees `json:"fees,omitempty"`
	FeesMaker map[string]Fees `json:"fees_maker,omitempty"`
}

// Fees - structure of fees info
type Fees struct {
	Fee        float64 `json:"fee,string"`
	MinFee     float64 `json:"minfee,string"`
	MaxFee     float64 `json:"maxfee,string"`
	NextFee    float64 `json:"nextfee,string"`
	NextVolume float64 `json:"nextvolume,string"`
	TierVolume float64 `json:"tiervolume,string"`
}

// CancelResponse - response on CancelOrder request
type CancelResponse struct {
	Count   int64 `json:"count"`
	Pending bool  `json:"pending,omitempty"`
}

// OrderDescription - structure of order description
type OrderDescription struct {
	Pair           string  `json:"pair"`
	Side           string  `json:"type"`
	OrderType      string  `json:"ordertype"`
	Price          float64 `json:"price,string"`
	Price2         float64 `json:"price2,string"`
	Leverage       string  `json:"leverage"`
	Info           string  `json:"order"`
	CloseCondition string  `json:"close"`
}

// AddOrderResponse - response on AddOrder request
type AddOrderResponse struct {
	Description    OrderDescription `json:"descr"`
	TransactionIds []string         `json:"txid"`
}

// GetWebSocketTokenResponse - response on GetWebSocketsToken request
type GetWebSocketTokenResponse struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}
