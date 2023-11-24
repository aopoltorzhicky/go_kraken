package futures

import (
	"fmt"
	"net/http"
	"testing"
)

func TestKraken_SendOrder(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    OrderBook
		wantErr bool
		live    bool
	}{
		{
			name:    "SendOrder to Kraken",
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
			got, err := api.SendOrder("PI_XBTUSD", OrderSideBuy, OrderTypeMarket, 0.001, 27000)
			if tt.live && err != nil {
				t.Errorf("Kraken.OrderBook() error = %v", err)
			}
			fmt.Println(got)
		})
	}
}

func TestKraken_Accounts(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    OrderBook
		wantErr bool
		live    bool
	}{
		{
			name:    "Accounts to Kraken",
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
			got, err := api.GetBalances()
			if tt.live && err != nil {
				t.Errorf("Kraken.GetBalances() error = %v", err)
			}
			fmt.Println(got)
		})
	}
}
