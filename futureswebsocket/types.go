package futureswebsocket

type EVENT_TYPE string

const (
	SUBSCRIBE    EVENT_TYPE = "subscribe"
	UNSUBSCRIBE  EVENT_TYPE = "unsubscribe"
	SUBSCRIBED   EVENT_TYPE = "subscribed"
	UNSUBSCRIBED EVENT_TYPE = "unsubscribed"
	INFO         EVENT_TYPE = "info"
)

type BOOK_SIDE string

const (
	BUY  EVENT_TYPE = "buy"
	SELL EVENT_TYPE = "sell"
)

type FEED_TYPE string

const (
	BOOK          FEED_TYPE = "book"
	BOOK_SNAPSHOT FEED_TYPE = "book_snapshot"
)

type SubscribeBook struct {
	Event      EVENT_TYPE `json:"event,omitempty"`
	Feed       FEED_TYPE  `json:"feed,omitempty"`
	ProductIds []string   `json:"product_ids,omitempty"`
}

type EventType struct {
	Event EVENT_TYPE `json:"event"`
	Feed  FEED_TYPE  `json:"feed,omitempty"`
}

type BookDepth struct {
	Price float64 `json:"price,omitempty"`
	Qty   float64 `json:"qty,omitempty"`
}

type BookSnapshotEvent struct {
	Feed      FEED_TYPE   `json:"feed,omitempty"`
	ProductId string      `json:"product_id,omitempty"`
	Timestamp int64       `json:"timestamp,omitempty"`
	Seq       int64       `json:"seq,omitempty"`
	TickSize  float64     `json:"tick_size,omitempty"`
	Bids      []BookDepth `json:"bids,omitempty"`
	Asks      []BookDepth `json:"asks,omitempty"`
}

type BookUpdateEvent struct {
	Feed      FEED_TYPE `json:"feed,omitempty"`
	ProductId string    `json:"product_id,omitempty"`
	Side      BOOK_SIDE `json:"side,omitempty"`
	Seq       int64     `json:"seq,omitempty"`
	Price     float64   `json:"price,omitempty"`
	Qty       float64   `json:"qty,omitempty"`
	Timestamp int64     `json:"timestamp,omitempty"`
}

type Update struct {
	Feed      FEED_TYPE   `json:"feed,omitempty"`
	ProductId string      `json:"product_id,omitempty"`
	Seq       int64       `json:"seq,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}
