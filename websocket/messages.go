package websocket

import "math/big"

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
	ChannelID    int64        `json:"chanId"`
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

type Ping struct {
	Event string `json:"event"`
	ReqID string `json:"reqid,omitempty"`
}

type Pong struct {
	Event string `json:"event"`
	ReqID string `json:"reqid,omitempty"`
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
	Close              Level
	Volume             Values
	VolumeAveragePrice Values
	TradeVolume        Values
	Low                Values
	High               Values
	Open               Values
}

type Level struct {
	Price          float64
	Volume         float64
	WholeLotVolume float64
}

type Values struct {
	Today  float64
	Last24 float64
}
