package websocket

import (
	"encoding/json"
	"testing"
)

// An order book configured with a depth larger than the number of levels a side
// actually has (a thin market, or a book that hasn't filled to `depth` levels
// yet) must not panic: levels[o.depth:] / levels[:o.depth] would be out of range.
func TestOrderBookApplyUpdateFewerLevelsThanDepth(t *testing.T) {
	ob := NewOrderBook(10, 1, 8) // depth 10

	upd := OrderBookUpdate{
		Asks: []OrderBookItem{
			{Price: json.Number("100.0"), Volume: json.Number("1.0")},
			{Price: json.Number("101.0"), Volume: json.Number("2.0")},
		},
		Bids: []OrderBookItem{
			{Price: json.Number("99.0"), Volume: json.Number("1.5")},
		},
		IsSnapshot: true,
	}

	if err := ob.ApplyUpdate(upd, false); err != nil {
		t.Fatalf("ApplyUpdate returned error: %v", err)
	}

	// Sanity: both sides kept exactly the levels provided (fewer than depth).
	if got := len(ob.Asks.sorted); got != 2 {
		t.Errorf("asks: got %d levels, want 2", got)
	}
	if got := len(ob.Bids.sorted); got != 1 {
		t.Errorf("bids: got %d levels, want 1", got)
	}
}

// Once a side exceeds the configured depth, it must be trimmed to depth.
func TestOrderBookApplyUpdateTrimsToDepth(t *testing.T) {
	ob := NewOrderBook(2, 1, 8) // depth 2

	upd := OrderBookUpdate{
		Asks: []OrderBookItem{
			{Price: json.Number("100.0"), Volume: json.Number("1.0")},
			{Price: json.Number("101.0"), Volume: json.Number("1.0")},
			{Price: json.Number("102.0"), Volume: json.Number("1.0")},
			{Price: json.Number("103.0"), Volume: json.Number("1.0")},
		},
		IsSnapshot: true,
	}

	if err := ob.ApplyUpdate(upd, false); err != nil {
		t.Fatalf("ApplyUpdate returned error: %v", err)
	}
	if got := len(ob.Asks.sorted); got != 2 {
		t.Errorf("asks: got %d levels, want 2 (trimmed to depth)", got)
	}
}
