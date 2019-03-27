package kraken_ws

import (
	"math/big"
	"time"
)

type EventType struct {
	Event string `json:"event"`
}

type Subscription struct {
	Name     string `json:"name"`
	Interval int64  `json:"interval,omitempty"`
	Depth    int64  `json:"depth,omitempty"`
}

type SubscriptionRequest struct {
	ReqID string `json:"reqid,omitempty"`
	Event string `json:"event"`

	Pairs        []string     `json:"pair"`
	Subscription Subscription `json:"subscription"`
}

type UnsubscribeRequest struct {
	Event        string       `json:"event"`
	Pairs        []string     `json:"pair"`
	Subscription Subscription `json:"subscription"`
}

type SubscriptionStatus struct {
	ChannelID    int64        `json:"channelID"`
	Event        string       `json:"event"`
	Status       string       `json:"status"`
	Pair         string       `json:"pair"`
	ReqID        string       `json:"reqid,omitempty"`
	Error        string       `json:"errorMessage,omitempty"`
	Subscription Subscription `json:"subscription"`
}

type PingRequest struct {
	Event string `json:"event"`
	ReqID int    `json:"reqid,omitempty"`
}

type PongResponse struct {
	Event string `json:"event"`
	ReqID int    `json:"reqid,omitempty"`
}

type SystemStatus struct {
	Event        string  `json:"event"`
	ConnectionID big.Int `json:"connectionID"`
	Status       string  `json:"status"`
	Version      string  `json:"version"`
}

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

type Level struct {
	Price          float64
	Volume         float64
	WholeLotVolume int
}

type Values struct {
	Today  interface{}
	Last24 interface{}
}

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

type TradeUpdate struct {
	Price     float64
	Volume    float64
	Time      time.Time
	Side      string
	OrderType string
	Misc      string
	Pair      string
}

type SpreadUpdate struct {
	Ask  float64
	Bid  float64
	Time time.Time
	Pair string
}

type OrderBookItem struct {
	Price  float64
	Volume float64
	Time   time.Time
}

type OrderBookUpdate struct {
	Asks       []OrderBookItem
	Bids       []OrderBookItem
	IsSnapshot bool
	Pair       string
}
