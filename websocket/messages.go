package websocket

import (
	"encoding/json"
	"fmt"
	"math/big"
)

// EventType - data structure for parsing events
type EventType struct {
	Event string `json:"event"`
}

// Subscription - data structure of subscription entity
type Subscription struct {
	Name     string `json:"name"`
	Interval int64  `json:"interval,omitempty"`
	Depth    int64  `json:"depth,omitempty"`
}

// SubscriptionRequest - data structure for subscription request
type SubscriptionRequest struct {
	ReqID string `json:"reqid,omitempty"`
	Event string `json:"event"`

	Pairs        []string     `json:"pair"`
	Subscription Subscription `json:"subscription"`
}

// UnsubscribeRequest - data structure for unsubscription request
type UnsubscribeRequest struct {
	Event        string       `json:"event"`
	Pairs        []string     `json:"pair"`
	Subscription Subscription `json:"subscription"`
}

// SubscriptionStatus - data structure for subscription status event
type SubscriptionStatus struct {
	ChannelID    int64        `json:"channelID"`
	Event        string       `json:"event"`
	Status       string       `json:"status"`
	Pair         string       `json:"pair"`
	ReqID        string       `json:"reqid,omitempty"`
	Error        string       `json:"errorMessage,omitempty"`
	Subscription Subscription `json:"subscription"`
}

// PingRequest - data structure for ping request
type PingRequest struct {
	Event string `json:"event"`
	ReqID int    `json:"reqid,omitempty"`
}

// PongResponse - data structure for ping response
type PongResponse struct {
	Event string `json:"event"`
	ReqID int    `json:"reqid,omitempty"`
}

// SystemStatus - data structure for system status event
type SystemStatus struct {
	Event        string  `json:"event"`
	ConnectionID big.Int `json:"connectionID"`
	Status       string  `json:"status"`
	Version      string  `json:"version"`
}

// DataUpdate - data structure of default Kraken WS update
type DataUpdate struct {
	ChannelID   int64
	Data        interface{}
	ChannelName string
	Pair        string
	Sequence	int64
}

// UnmarshalJSON - unmarshal update
func (u *DataUpdate) UnmarshalJSON(data []byte) error {
	var raw []interface{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	if len(raw) < 3 {
		return fmt.Errorf("invalid data length: %#v", raw)
	}

	if len(raw) == 3 {
		var ok bool
		u.Data = raw[0]

		if u.ChannelName, ok = raw[1].(string); !ok {
			return fmt.Errorf("expected message to have channel name as 2nd element but got %#v instead", raw[1])
		}

		var sequenceMap map[string]interface{}
		if sequenceMap, ok = raw[2].(map[string]interface{}); !ok {
			return fmt.Errorf("expected message to have JSON object as 3rd element but got %#v instead", raw[2])
		}

		var sequenceRaw interface{}
		if sequenceRaw, ok = sequenceMap["sequence"]; !ok {
			return fmt.Errorf("expected message to have sequence in JSON object as 3rd element but got %#v instead", raw[2])
		}

		var seq float64
		if seq, ok = sequenceRaw.(float64); !ok {
			return fmt.Errorf("expected message to have sequence integer in JSON object as 3rd element but got %#v instead", raw[2])
		}

		u.Sequence = int64(seq)
		return nil
	}

	chID, ok := raw[0].(float64)
	if !ok {
		return fmt.Errorf("expected message to start with a channel id but got %#v instead", raw[0])
	}

	u.ChannelID = int64(chID)
	u.ChannelName, ok = raw[len(raw)-2].(string)
	if !ok {
		return fmt.Errorf("expected message with (n - 2) element channel name but got %#v instead", raw[len(raw)-2])
	}
	u.Pair, ok = raw[len(raw)-1].(string)
	if !ok {
		return fmt.Errorf("expected message  with (n - 2) element pair but got %#v instead", raw[len(raw)-1])
	}
	u.Data = raw[1 : len(raw)-2][0]

	return nil
}

// TickerUpdate - data structure for ticker update
type TickerUpdate struct {
	Ask                Level
	Bid                Level
	Close              Values
	Volume             Values
	VolumeAveragePrice Values
	TradeVolume        Values
	Low                Values
	High               Values
	Open               Values
	Pair               string
}

// Level - data structure for ticker data ask/bid
type Level struct {
	Price          float64
	Volume         float64
	WholeLotVolume int
}

// Values - data structure for ticker others data
type Values struct {
	Today  interface{}
	Last24 interface{}
}

// CandleUpdate - data structure for candles update
type CandleUpdate struct {
	Time      float64
	EndTime   float64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	VolumeWAP float64
	Volume    float64
	Count     int
	Pair      string
}

// TradeUpdate - data structure for trade update
type TradeUpdate struct {
	Price     float64
	Volume    float64
	Time      float64
	Side      string
	OrderType string
	Misc      string
	Pair      string
}

// SpreadUpdate - data structure for spread update
type SpreadUpdate struct {
	Ask       float64
	Bid       float64
	AskVolume float64
	BidVolume float64
	Time      float64
	Pair      string
}

// OrderBookItem - data structure for order book item
type OrderBookItem struct {
	Price     float64
	Volume    float64
	Time      float64
	Republish bool
}

// OrderBookUpdate - data structure for order book update
type OrderBookUpdate struct {
	Asks       []OrderBookItem
	Bids       []OrderBookItem
	IsSnapshot bool
	Pair       string
	CheckSum   string
}

// AuthDataRequest - data structure for private subscription request
type AuthDataRequest struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

// AuthSubscriptionRequest - data structure for private subscription request
type AuthSubscriptionRequest struct {
	Event string          `json:"event"`
	Subs  AuthDataRequest `json:"subscription"`
}

// OwnTrade - Own trades.
type OwnTrade struct {
	Cost      float64 `json:"cost,string"`
	Fee       float64 `json:"fee,string"`
	Margin    float64 `json:"margin,string"`
	OrderID   string  `json:"ordertxid"`
	OrderType string  `json:"ordertype"`
	Pair      string  `json:"pair"`
	PosTxID   string  `json:"postxid"`
	Price     float64 `json:"price,string"`
	Time      float64 `json:"time"`
	Type      string  `json:"type"`
	Vol       float64 `json:"vol,string"`
}

// OpenOrderDescr -
type OpenOrderDescr struct {
	Close     string  `json:"close"`
	Leverage  string  `json:"leverage"`
	Order     string  `json:"order"`
	Ordertype string  `json:"ordertype"`
	Pair      string  `json:"pair"`
	Price     float64 `json:"price,string"`
	Price2    float64 `json:"price2,string"`
	Type      string  `json:"type"`
}

// OpenOrder -
type OpenOrder struct {
	Cost       float64        `json:"cost,string"`
	Descr      OpenOrderDescr `json:"descr"`
	Fee        float64        `json:"fee,string"`
	LimitPrice float64        `json:"limitprice,string"`
	Misc       string         `json:"misc"`
	Oflags     string         `json:"oflags"`
	OpenTime   float64        `json:"opentm"`
	StartTime  float64        `json:"starttm"`
	ExpireTime float64        `json:"expiretm"`
	Price      float64        `json:"price,string"`
	Refid      string         `json:"refid"`
	Status     string         `json:"status"`
	StopPrice  float64        `json:"stopprice,string"`
	UserRef    int            `json:"userref"`
	Vol        float64        `json:"vol,string"`
	VolExec    float64        `json:"vol_exec,string"`
}

// OwnTradesUpdate -
type OwnTradesUpdate struct {
	Trades      map[string]OwnTrade
	ChannelName string
}

// OpenOrdersUpdate -
type OpenOrdersUpdate struct {
	Order       map[string]OpenOrder
	ChannelName string
}

// AuthRequest -
type AuthRequest struct {
	Token string `json:"token"`
	Event string `json:"event"`
}

// AddOrderRequest -
type AddOrderRequest struct {
	AuthRequest
	ClosePrice string `json:"close[price]"`
	Ordertype  string `json:"ordertype"`
	Pair       string `json:"pair"`
	Price      string `json:"price"`
	Type       string `json:"type"`
	Volume     string `json:"volume"`
}

// AddOrderResponse -
type AddOrderResponse struct {
	Description  string `json:"descr"`
	Event        string `json:"event"`
	Status       string `json:"status"`
	TxID         string `json:"txid"`
	ErrorMessage string `json:"errorMessage,omiempty"`
}

// CancelOrderRequest -
type CancelOrderRequest struct {
	AuthRequest
	TxID []string `json:"txid"`
}

// CancelOrderResponse -
type CancelOrderResponse struct {
	ErrorMessage string `json:"errorMessage,omiempty"`
	Event        string `json:"event"`
	Status       string `json:"status"`
}
