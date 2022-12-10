package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
)

func getFloat64FromStr(value interface{}) (float64, error) {
	str, ok := value.(string)
	if !ok {
		return .0, errors.New("field must be a string")
	}
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return .0, err
	}
	return f, nil
}

func getFloat64(value interface{}) (float64, error) {
	f, ok := value.(float64)
	if !ok {
		return .0, errors.New("field must be a float64")
	}
	return f, nil
}

func getTimestamp(value interface{}) (int64, error) {
	f, ok := value.(float64)
	if !ok {
		return 0, errors.New("field must be a float64")
	}
	return int64(f), nil
}

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

// Asset - asset information
type Asset struct {
	AlternateName   string `json:"altname"`
	AssetClass      string `json:"aclass"`
	Decimals        int    `json:"decimals"`
	DisplayDecimals int    `json:"display_decimals"`
}

// AssetPair - asset pair information
type AssetPair struct {
	Altname           string          `json:"altname"`
	AssetClassBase    string          `json:"aclass_base"`
	Base              string          `json:"base"`
	AssetClassQuote   string          `json:"aclass_quote"`
	Quote             string          `json:"quote"`
	Lot               string          `json:"lot"`
	PairDecimals      int             `json:"pair_decimals"`
	LotDecimals       int             `json:"lot_decimals"`
	LotMultiplier     int             `json:"lot_multiplier"`
	LeverageBuy       []int           `json:"leverage_buy"`
	LeverageSell      []int           `json:"leverage_sell"`
	Fees              [][]float64     `json:"fees"`
	FeesMaker         [][]float64     `json:"fees_maker"`
	FeeVolumeCurrency string          `json:"fee_volume_currency"`
	MarginCall        int             `json:"margin_call"`
	MarginStop        int             `json:"margin_stop"`
	WSName            string          `json:"wsname"`
	OrderMin          decimal.Decimal `json:"ordermin"`
}

// Level - ticker structure for Ask and Bid
type Level struct {
	Price          decimal.Decimal
	WholeLotVolume decimal.Decimal
	Volume         decimal.Decimal
}

// UnmarshalJSON -
func (item *Level) UnmarshalJSON(buf []byte) error {
	var tmp []decimal.Decimal
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), 3; g != e {
		return fmt.Errorf("wrong number of fields in Level: %d != %d", g, e)
	}

	item.Price = tmp[0]
	item.WholeLotVolume = tmp[1]
	item.Volume = tmp[2]
	return nil
}

// TimeLevel - ticker structure
type TimeLevel struct {
	Today       int64
	Last24Hours int64
}

// UnmarshalJSON -
func (item *TimeLevel) UnmarshalJSON(buf []byte) error {
	var tmp []int64
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), 2; g != e {
		return fmt.Errorf("wrong number of fields in TimeLevel: %d != %d", g, e)
	}
	item.Today = tmp[0]
	item.Last24Hours = tmp[1]
	return nil
}

// CloseLevel - ticker structure for Close
type CloseLevel struct {
	Price     decimal.Decimal
	LotVolume decimal.Decimal
}

// UnmarshalJSON -
func (item *CloseLevel) UnmarshalJSON(buf []byte) error {
	var tmp []decimal.Decimal
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), 2; g != e {
		return fmt.Errorf("wrong number of fields in CloseLevel: %d != %d", g, e)
	}

	item.Price = tmp[0]
	item.LotVolume = tmp[1]

	return nil
}

// Ticker - struct of ticker response
type Ticker struct {
	Ask                Level      `json:"a"`
	Bid                Level      `json:"b"`
	Close              CloseLevel `json:"c"`
	Volume             CloseLevel `json:"v"`
	VolumeAveragePrice CloseLevel `json:"p"`
	Trades             TimeLevel  `json:"t"`
	Low                CloseLevel `json:"l"`
	High               CloseLevel `json:"h"`
	OpeningPrice       decimal.Decimal
}

// Candle - OHLC item
type Candle struct {
	Time      int64
	Open      decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Close     decimal.Decimal
	VolumeWAP decimal.Decimal
	Volume    decimal.Decimal
	Count     int64
}

// OHLCResponse - response of OHLC request
type OHLCResponse struct {
	Candles map[string][]Candle `json:"-"`
	Last    int64               `json:"last"`
}

// UnmarshalJSON -
func (item *OHLCResponse) UnmarshalJSON(buf []byte) error {
	res := make(map[string]interface{})
	if err := json.Unmarshal(buf, &res); err != nil {
		return err
	}

	last, err := getTimestamp(res["last"])
	if err != nil {
		return err
	}
	item.Last = last
	delete(res, "last")

	item.Candles = make(map[string][]Candle)
	for k, v := range res {
		items := v.([]interface{})
		item.Candles[k] = make([]Candle, len(items))
		for idx, c := range items {
			candle := c.([]interface{})

			ts, err := getTimestamp(candle[0])
			if err != nil {
				continue
			}
			open, err := decimal.NewFromString(candle[1].(string))
			if err != nil {
				continue
			}
			high, err := decimal.NewFromString(candle[2].(string))
			if err != nil {
				continue
			}
			low, err := decimal.NewFromString(candle[3].(string))
			if err != nil {
				continue
			}
			close, err := decimal.NewFromString(candle[4].(string))
			if err != nil {
				continue
			}
			vwap, err := decimal.NewFromString(candle[5].(string))
			if err != nil {
				continue
			}
			vol, err := decimal.NewFromString(candle[6].(string))
			if err != nil {
				continue
			}
			item.Candles[k][idx] = Candle{
				Time:      ts,
				Open:      open,
				High:      high,
				Low:       low,
				Close:     close,
				VolumeWAP: vwap,
				Volume:    vol,
				Count:     int64(candle[7].(float64)),
			}
		}
	}
	return nil
}

// OrderBookItem - one price level in orderbook
type OrderBookItem struct {
	Price     float64
	Volume    float64
	Timestamp int64
}

// UnmarshalJSON -
func (item *OrderBookItem) UnmarshalJSON(buf []byte) error {
	var tmp []interface{}
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), 3; g != e {
		return fmt.Errorf("wrong number of fields in OrderBookItem: %d != %d", g, e)
	}

	price, err := getFloat64FromStr(tmp[0])
	if err != nil {
		return err
	}
	item.Price = price

	vol, err := getFloat64FromStr(tmp[1])
	if err != nil {
		return err
	}
	item.Volume = vol

	ts, err := getTimestamp(tmp[2])
	if err != nil {
		return err
	}
	item.Timestamp = ts

	return nil
}

// OrderBook - struct of order book levels
type OrderBook struct {
	Asks []OrderBookItem `json:"asks"`
	Bids []OrderBookItem `json:"bids"`
}

// Trade - structure of public trades
type Trade struct {
	Price     float64
	Volume    float64
	Time      float64
	Side      string
	OrderType string
	Misc      string
}

// UnmarshalJSON -
func (item *Trade) UnmarshalJSON(buf []byte) error {
	var tmp []interface{}
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), 6; g != e {
		return fmt.Errorf("wrong number of fields in CloseLevel: %d != %d", g, e)
	}

	price, err := getFloat64FromStr(tmp[0])
	if err != nil {
		return err
	}
	item.Price = price

	vol, err := getFloat64FromStr(tmp[1])
	if err != nil {
		return err
	}
	item.Volume = vol

	ts, err := getFloat64(tmp[2])
	if err != nil {
		return err
	}
	item.Time = ts

	side, ok := tmp[3].(string)
	if !ok {
		return errors.New("invalid side type")
	}
	item.Side = side

	t, ok := tmp[4].(string)
	if !ok {
		return errors.New("invalid order type")
	}
	item.OrderType = t

	misc, ok := tmp[5].(string)
	if !ok {
		return errors.New("invalid misc type")
	}
	item.Misc = misc
	return nil
}

// TradeResponse - all pairs in trade response
type TradeResponse struct {
	Last     string `json:"last"`
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
	Bid  float64
	Ask  float64
}

// UnmarshalJSON -
func (item *Spread) UnmarshalJSON(buf []byte) error {
	var tmp []interface{}
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), 3; g != e {
		return fmt.Errorf("wrong number of fields in CloseLevel: %d != %d", g, e)
	}

	ts, err := getFloat64(tmp[0])
	if err != nil {
		return err
	}
	item.Time = ts

	bid, err := getFloat64FromStr(tmp[1])
	if err != nil {
		return err
	}
	item.Bid = bid

	ask, err := getFloat64FromStr(tmp[2])
	if err != nil {
		return err
	}
	item.Ask = ask
	return nil
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
	RefID           *string          `json:"refid"`
	UserRef         interface{}      `json:"userref"`
	Status          string           `json:"status"`
	Reason          string           `json:"reason,omitempty"`
	OpenTimestamp   float64          `json:"opentm"`
	StartTimestamp  float64          `json:"starttm"`
	CloseTimestamp  float64          `json:"closetm,omitempty"`
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

// DepositMethods - respons on GetDepositMethods request
type DepositMethods struct {
	Method     string `json:"method"`
	Fee        string `json:"fee"`
	Limit      bool   `json:"limit"`
	GenAddress bool   `json:"gen-address"`
}

// GetDepositStatus - respons on GetDepositMethods request
type DepositStatuses struct {
	Method string `json:"method"`
	Aclass string `json:"aclass"`
	Asset  string `json:"asset"`
	Refid  string `json:"refid"`
	Txid   string `json:"txid"`
	Info   string `json:"info"`
	Amount string `json:"amount"`
	Fee    string `json:"fee"`
	Time   int    `json:"time"`
	Status string `json:"status"`
}

// WithdrawFunds - response on WithdrawFunds request
type WithdrawFunds struct {
	RefID string `json:"refid"`
}

// GetWithdrawStatus - response on WithdrawStatus request
type WithdrawStatus struct {
	Method string `json:"method,omitempty"`
	AClass string `json:"a_class,omitempty"`
	Asset  string `json:"asset,omitempty"`
	Refid  string `json:"refid,omitempty"`
	Txid   string `json:"txid,omitempty"`
	Info   string `json:"info,omitempty"`
	Amount string `json:"amount,omitempty"`
	Fee    string `json:"fee,omitempty"`
	Time   int    `json:"time,omitempty"`
	Status string `json:"status,omitempty"`
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

// EditOrderResponse - response on EditOrder request
type EditOrderResponse struct {
	Description     OrderDescription `json:"descr"`
	TransactionId   string           `json:"txid"`
	OrdersCancelled int64            `json:"orders_cancelled"`
	Volume          string           `json:"volume"`
	Status          string           `json:"status"`
	Price           float64          `json:"price,string"`
	Price2          float64          `json:"price2,string"`
	ErrorMessage    float64          `json:"error_message,string"`
}

// GetWebSocketTokenResponse - response on GetWebSocketsToken request
type GetWebSocketTokenResponse struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}
