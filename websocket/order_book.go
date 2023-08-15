package websocket

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"strings"

	"github.com/pkg/errors"
)

// OrderBook -
type OrderBook struct {
	Asks *OrderBookSide
	Bids *OrderBookSide
}

// NewOrderBook - creates order book.
//
//	depth - is a requested depth from Kraken
//
//	pricePrecision - count of valuable signs after dot in price, which is required for checksum verification
//
//	volumePrecision - count of valuable signs after dot in volume, which is required for checksum verification
func NewOrderBook(depth, pricePrecision, volumePrecision int) *OrderBook {
	return &OrderBook{
		Asks: newOrderBookSide(depth, pricePrecision, volumePrecision, true),
		Bids: newOrderBookSide(depth, pricePrecision, volumePrecision, false),
	}
}

// ApplyUpdate - applies updates from kraken websocket.
// If you need to verify checksum, set verify to true.
func (o *OrderBook) ApplyUpdate(upd OrderBookUpdate, verify bool) error {
	if err := o.Asks.applyUpdates(upd.Asks); err != nil {
		return err
	}
	if err := o.Bids.applyUpdates(upd.Bids); err != nil {
		return err
	}

	if verify && !upd.IsSnapshot {
		if cs := o.Checksum(); cs != upd.CheckSum {
			return errors.Errorf("invalid checksum: local %s != remote %s", cs, upd.CheckSum)
		}
	}
	return nil
}

// Checksum - computes order book checksum. Details https://docs.kraken.com/websockets/#book-checksum
func (o *OrderBook) Checksum() string {
	var str bytes.Buffer
	str.Write(o.Asks.checksum())
	str.Write(o.Bids.checksum())
	return fmt.Sprint(crc32.ChecksumIEEE(str.Bytes()))
}

// String - returns full order book as a string
func (o *OrderBook) String() string {
	var builder strings.Builder
	builder.WriteString("\r\n==== ASKS ====\r\n")
	builder.WriteString(o.Asks.String())
	builder.WriteString("==== BIDS ====\r\n")
	builder.WriteString(o.Bids.String())
	return builder.String()
}
