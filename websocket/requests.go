package websocket

import (
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
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// CancelOrderRequest -
type CancelOrderRequest struct {
	AuthRequest
	TxID []string `json:"txid"`
}

// CancelOrderResponse -
type CancelOrderResponse struct {
	ErrorMessage string `json:"errorMessage,omitempty"`
	Event        string `json:"event"`
	Status       string `json:"status"`
}

// CancelAllResponse -
type CancelAllResponse struct {
	Count        int    `json:"count"`
	Event        string `json:"event"`
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}
