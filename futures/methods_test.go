package futures

import (
	"fmt"
	"net/http"
	"testing"
)

var ErrSomething = fmt.Errorf("something went wrong")

type httpMock struct {
	Response *http.Response
	Error    error
}

var krakenLive = NewFromEnv()

func (c *httpMock) Do(req *http.Request) (*http.Response, error) {
	if c.Error != nil {
		return c.Response, c.Error
	}
	return c.Response, nil
}

func TestKraken_OrderBook(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    OrderBook
		wantErr bool
		live    bool
	}{
		{
			name:    "Live from Kraken",
			err:     nil,
			want:    OrderBook{},
			wantErr: false,
			live:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var api *KrakenFutures
			if tt.live {
				api = krakenLive
			} else {
				api = &KrakenFutures{
					client: &httpMock{
						Error:    tt.err,
						Response: tt.resp,
					},
				}
			}
			got, err := api.OrderBook("PI_XBTUSD")
			if tt.live && err != nil {
				t.Errorf("Kraken.OrderBook() error = %v", err)
			}
			if got.Asks[0][0] < 1 {
				t.Errorf("Kraken.OrderBook() got issue = %v", got)
			}
		})
	}
}
