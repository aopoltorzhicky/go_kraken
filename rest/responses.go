package rest

// KrakenResponse - template of Kraken API response
type KrakenResponse struct {
	Error  []string    `json:"error"`
	Result interface{} `json:"result"`
}

// TimeResponse - Result of Time request
type TimeResponse struct {
	Unixtime int64
	Rfc1123  string
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

// Ticker - struct of ticker response
type Ticker struct {
	Ask                []string `json:"a"`
	Bid                []string `json:"b"`
	Close              []string `json:"c"`
	Volume             []string `json:"v"`
	VolumeAveragePrice []string `json:"p"`
	Trades             []int    `json:"t"`
	Low                []string `json:"l"`
	High               []string `json:"h"`
	OpeningPrice       float64  `json:"o,string"`
}

// Candle - OHLC item
type Candle struct {
	Time      int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	VolumeWAP float64
	Volume    float64
	Count     int64
}

// OHLCResponse - response of OHLC request
type OHLCResponse struct {
	Last    int64
	Candles []Candle
}

// OrderBookItem - one price level in orderbook
type OrderBookItem struct {
	Price     float64
	Volume    float64
	Timestamp int64
}

// OrderBook - struct of order book levels
type OrderBook struct {
	Asks []OrderBookItem `json:"asks"`
	Bids []OrderBookItem `json:"bids"`
}
