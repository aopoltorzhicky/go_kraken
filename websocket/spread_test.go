package websocket

import (
	"encoding/json"
	"testing"
)

// Kraken v1 spread payload order is [bid, ask, timestamp, bidVolume, askVolume]
// (see the example at
// https://docs.kraken.com/exchange/api-reference/spot-websocket-v1/spread).
// The decoder previously mapped position 3 to AskVolume and position 4 to
// BidVolume, silently swapping the two volumes.
func TestSpreadUnmarshalVolumesNotSwapped(t *testing.T) {
	// Distinct volumes so a swap is detectable.
	payload := `["5698.40000","5700.00000","1542057299.545897","1.01234567","0.98765432"]`

	var s Spread
	if err := json.Unmarshal([]byte(payload), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if s.Bid.String() != "5698.40000" {
		t.Errorf("Bid = %s, want 5698.40000", s.Bid)
	}
	if s.Ask.String() != "5700.00000" {
		t.Errorf("Ask = %s, want 5700.00000", s.Ask)
	}
	if s.BidVolume.String() != "1.01234567" {
		t.Errorf("BidVolume = %s, want 1.01234567 (position 3)", s.BidVolume)
	}
	if s.AskVolume.String() != "0.98765432" {
		t.Errorf("AskVolume = %s, want 0.98765432 (position 4)", s.AskVolume)
	}
}
