package websocket

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"
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
}

// UnmarshalJSON - unmarshal update
func (u *DataUpdate) UnmarshalJSON(data []byte) error {
	var raw []interface{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	} else if len(raw) < 4 {
		return fmt.Errorf("invalid data length: %#v", raw)
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
	Time      time.Time
	EndTime   time.Time
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
	Time      time.Time
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
	Time      time.Time
	Pair      string
}

// OrderBookItem - data structure for order book item
type OrderBookItem struct {
	Price     float64
	Volume    float64
	Time      time.Time
	Republish bool
}

// OrderBookUpdate - data structure for order book update
type OrderBookUpdate struct {
	Asks       []OrderBookItem
	Bids       []OrderBookItem
	IsSnapshot bool
	Pair       string
}
