package websocket

import (
	"encoding/json"
	"testing"
)

func TestOrderBookMessageBid(t *testing.T) {
	bidMessage := `[336,{"b":[["50251.20000","0.00000000","1638472269.482087"],["47190.00000","0.01000000","1638216666.203784","r"]],"c":"2271011438"},"book-1000","XBT/EUR"]`
	var msg Message
	if err := json.Unmarshal([]byte(bidMessage), &msg); err != nil {
		t.Error("could not parse message:", err)
		return
	}

	var update OrderBookUpdate
	if err := json.Unmarshal(msg.Data, &update); err != nil {
		t.Error("could not parse message:", err)
		return
	}
	if len(update.Asks) != 0 {
		t.Error("expected 0 asks, got", len(update.Asks))
	}
	if len(update.Bids) != 2 {
		t.Error("expected 2 bids, got", len(update.Bids))
	}
}

func TestOrderBookMessageAskBid(t *testing.T) {
	askBidMessage := `[336,{"a":[["51982.90000","0.00000000","1638471905.000103"],["53700.00000","0.30322268","1638286260.153968","r"]]},{"b":[["48489.80000","0.00000000","1638471905.000059"],["47112.00000","0.00212260","1638286104.076640","r"]],"c":"2144637680"},"book-1000","XBT/EUR"]`
	var msg Message

	if err := json.Unmarshal([]byte(askBidMessage), &msg); err != nil {
		t.Error("could not parse message:", err)
		return
	}

	var update OrderBookUpdate
	if err := json.Unmarshal(msg.Data, &update); err != nil {
		t.Error("could not parse order book update:", err)
		return
	}
	if len(update.Asks) != 2 {
		t.Error("expected 2 asks, got", len(update.Asks))
	}
	if len(update.Bids) != 2 {
		t.Error("expected 2 bids, got", len(update.Bids))
	}
}
