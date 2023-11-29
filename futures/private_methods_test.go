package futures

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestKraken_GetOrderStatus(t *testing.T) {
	json := []byte(`{"result":"success","serverTime":"2023-11-28T09:08:32.822Z","orders":[{"order":{"type":"ORDER","orderId":"848dc18a-bbf9-46ea-a88f-e7419db216f0","cliOrdId":null,"symbol":"PF_XBTUSD","side":"buy","quantity":0.001,"filled":0,"limitPrice":35000,"reduceOnly":false,"timestamp":"2023-11-28T09:06:52.784Z","lastUpdateTimestamp":"2023-11-28T09:06:52.784Z"},"status":"ENTERED_BOOK","updateReason":null,"error":null}]}`)
	// errSomething := fmt.Errorf("something went wrong")
	tests := []struct {
		name      string
		cliOrdIds []string
		orderIds  []string
		want      OrderStatusResponse
		resp      *http.Response
		wantErr   bool
		err       error
		live      bool
	}{
		{
			name:      "GetOrderStatus Valid OrderID",
			cliOrdIds: nil,
			orderIds:  []string{"848dc18a-bbf9-46ea-a88f-e7419db216f0"},
			want:      OrderStatusResponse{},
			resp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(json)),
			},
			wantErr: false,
			live:    false,
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
			got, err := api.GetOrderStatus(tt.cliOrdIds, tt.orderIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrderStatus() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(got.Orders) > 0 && got.Orders[0].Status != "ENTERED_BOOK" {
				t.Errorf("GetOrderStatus() got first order status = %v, want status = 'ENTERED_BOOK'", got.Orders[0].Status)
			}
		})
	}
}

func TestKraken_SendOrder(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		resp      *http.Response
		want      OrderBook
		wantErr   bool
		live      bool
		orderSide string
	}{
		{
			name:      "SendOrder BUY to Kraken",
			err:       nil,
			want:      OrderBook{},
			wantErr:   false,
			live:      true,
			orderSide: OrderSideBuy,
		},
		{
			name:      "SendOrder SELL to Kraken",
			err:       nil,
			want:      OrderBook{},
			wantErr:   false,
			live:      true,
			orderSide: OrderSideSell,
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
			got, err := api.SendOrder(
				"pf_xbtusd",
				tt.orderSide,
				OrderTypeLimit,
				0.001,
				35000)
			if err != nil {
				t.Errorf("SendOrder() error = %v, wantErr %v", err, tt.wantErr)
			} else if got.OrderID == "" || got.Status != "placed" {
				t.Errorf("SendOrder() got = %v, want status = 'placed'", got)
			}
		})
	}
}

func TestKraken_Accounts(t *testing.T) {
	json := []byte(`{
		"result": "success",
		"accounts": {
		  "flex": {
			"currencies": {
			  "USDT": {
				"quantity": 1000.0,
				"value": 1000.189
			  }
			},
			"initialMargin": 0.7393,
			"balanceValue": 2997.86,
			"type": "multiCollateralMarginAccount"
		  }
		},
		"serverTime": "2023-11-29T08:20:38.194Z"
	  }`)
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    AccountsResponse
		wantErr bool
		live    bool
	}{
		{
			name: "Accounts to Kraken",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(json)),
			},
			want: AccountsResponse{
				Result: "success",
				Accounts: Accounts{
					Flex: FlexAccount{
						Currencies: map[string]CurrencyDetail{
							"USDT": {
								Quantity: 1000.0,
								Value:    1000.189,
							},
						},
						InitialMargin: 0.7393,
						BalanceValue:  2997.86,
						Type:          "multiCollateralMarginAccount",
					},
				},
				ServerTime: "2023-11-29T08:20:38.194Z",
			},
			wantErr: false,
			live:    false,
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
			got, err := api.GetAccounts()
			if tt.live && err != nil {
				t.Errorf("Kraken.GetAccounts() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.GetAccounts() = %v, want %v", got, tt.want)
			}
		})
	}
}
