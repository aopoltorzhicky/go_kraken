package rest

import (
	"encoding/json"
	"testing"
)

// OHLCResponse.UnmarshalJSON parses a REST response (untrusted network data)
// with type assertions and fixed index accesses. Malformed or unexpected
// payloads must be skipped, not panic the client.
func TestOHLCResponseUnmarshalMalformedDoesNotPanic(t *testing.T) {
	cases := map[string]string{
		"short candle":         `{"last":123,"XBTUSD":[[123,"1","2","3","4","5","6"]]}`,
		"numeric where string": `{"last":123,"XBTUSD":[[123,1,2,3,4,5,6,7]]}`,
		"candle not an array":  `{"last":123,"XBTUSD":["notarray"]}`,
		"value not an array":   `{"last":123,"XBTUSD":42}`,
		"count not a float":    `{"last":123,"XBTUSD":[[123,"1","2","3","4","5","6","x"]]}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			var o OHLCResponse
			// Must not panic; an error is acceptable, a panic is not.
			_ = json.Unmarshal([]byte(body), &o)
		})
	}
}

// A well-formed OHLC response still decodes correctly.
func TestOHLCResponseUnmarshalValid(t *testing.T) {
	body := `{"last":1500000000,"XBTUSD":[[1500000000,"100.1","101.2","99.3","100.5","100.4","12.3",7]]}`
	var o OHLCResponse
	if err := json.Unmarshal([]byte(body), &o); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	candles, ok := o.Candles["XBTUSD"]
	if !ok || len(candles) != 1 {
		t.Fatalf("expected 1 candle for XBTUSD, got %#v", o.Candles)
	}
	c := candles[0]
	if c.Open.String() != "100.1" || c.Close.String() != "100.5" || c.Count != 7 {
		t.Errorf("unexpected candle: open=%s close=%s count=%d", c.Open, c.Close, c.Count)
	}
}
