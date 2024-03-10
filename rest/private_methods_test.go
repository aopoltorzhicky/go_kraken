package rest

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var (
	depositMethodsJSON  = []byte(`{"error":[],"result":[{"method": "Ether (Hex)","limit": false,"fee": "0.0000000000","gen-address": true}]}`)
	depositStatusesJSON = []byte(`{"error":[],"result":[{"method": "Ether (Hex)","aclass": "currency","asset": "XETH","refid": "sometest1","txid": "sometest2","info": "sometest3","amount": "6.91","fee": "0.0000000000","time": 1617014556,"status": "Success"}]}`)
	balancesJSON        = []byte(`{"error":[],"result":{"ZUSD":"435.9135","USDT":"2.00000000","BSV":"0.0000053898"}}`)
	balancesExJSON      = []byte(`{"error":[],"result":{"ZUSD":{"balance":25435.21,"hold_trade":8249.76},"XXBT":{"balance":1.2435,"hold_trade":0.8423}}}`)
	tradeBalancesJSON   = []byte(`{"error":[],"result":{"eb":"33.50","tb":"33.50","m":"23.77","n":"4.3750","c":"11.8999","v":"12.2","e":"32.1","mf":"33.1","ml":"12.97"}}`)
	openOrdersJSON      = []byte(`{"error":[],"result":{"open":{"OR3XZM-5EN2R-LS5X51":{"refid":null,"userref":null,"status":"open","opentm":1570622342.3552,"starttm":0,"expiretm":0,"descr":{"pair":"XBTEUR","type":"sell","ordertype":"limit","price":"7712.2","price2":"0","leverage":"4:1","order":"sell 1.10000000 XBTEUR @ limit 7712.2 with 4:1 leverage","close":""},"vol":"1.10000000","vol_exec":"0.00000000","cost":"0.00000","fee":"0.00000","price":"0.00000","stopprice":"0.00000","limitprice":"0.00000","misc":"","oflags":"fciq"}}}}`)
	closedOrdersJSON    = []byte(`{"error":[],"result":{"closed":{"OK46ER-A2BXK-YOLKE1":{"refid":null,"userref":null,"status":"canceled","reason":"User requested","opentm":1570623817.6537,"closetm":1570623823.9012,"starttm":0,"expiretm":0,"descr":{"pair":"ETHEUR","type":"buy","ordertype":"limit","price":"160.87","price2":"0","leverage":"4:1","order":"buy 21.00000000 ETHEUR @ limit 160.87 with 4:1 leverage","close":""},"vol":"21.00000000","vol_exec":"0.00000000","cost":"0.00000","fee":"0.00000","price":"0.00000","stopprice":"0.00000","limitprice":"0.00000","misc":"","oflags":"fciq"}},"count":20}}`)
	queryOrdersJSON     = []byte(`{"error":[],"result":{"OLNYE1-H3BBJ-JD2LGC":{"refid":null,"userref":null,"status":"canceled","reason":"User requested","opentm":1570623816.1101,"closetm":1570623819.639,"starttm":0,"expiretm":0,"descr":{"pair":"XBTUSD","type":"buy","ordertype":"limit","price":"7920.9","price2":"0","leverage":"4:1","order":"buy 1.10000000 XBTUSD @ limit 7920.9 with 4:1 leverage","close":""},"vol":"1.10000000","vol_exec":"0.00000000","cost":"0.00000","fee":"0.00000","price":"0.00000","stopprice":"0.00000","limitprice":"0.00000","misc":"","oflags":"fciq"}}}`)
	tradeHistoryJSON    = []byte(`{"error":[],"result":{"trades":{"TO3MMA-BSBGV-XUV4A1":{"ordertxid":"OSQQQ5-MBKL6-O4YYE1","postxid":"TYE7IH-QCG76-BVMCM1","pair":"XXBTZUSD","time":1570477513.2,"type":"buy","ordertype":"limit","price":"7000.60000","cost":"1000.38301","fee":"0.00000","vol":"0.2","margin":"320","misc":"closing"}},"count":1}}`)
	queryTradesJSON     = []byte(`{"error":[],"result":{"TO3MMA-BSBGV-XUV4A1":{"ordertxid":"OSQQQ5-MBKL6-O4YYE1","postxid":"TYE7IH-QCG76-BVMCM1","pair":"XXBTZUSD","time":1570477513.2,"type":"buy","ordertype":"limit","price":"7000.60000","cost":"1000.38301","fee":"0.00000","vol":"0.2","margin":"320","misc":"closing"}}}`)
	openPositionsJSON   = []byte(`{"error":[],"result":{"TYE7IH-QCG76-BVMCM1":{"ordertxid":"OK7SOC-SGF3O-F54S51","posstatus":"open","pair":"XXBTZUSD","time":1569513333.0361,"type":"buy","ordertype":"limit","cost":"570.39712","fee":"39","vol":"7","vol_closed":"6.66208817","margin":"9.2","terms":"0.0100% per 4 hours","rollovertm":"1570638129","misc":"","oflags":""}}}`)
	getLedgersJSON      = []byte(`{"error":[],"result":{"ledger":{"LGPNZQ-2SLSA-C7QCT1":{"refid":"TI2NBU-IICD2-BAVYO1","time":1570623111.9096,"type":"rollover","aclass":"currency","asset":"ZUSD","amount":"0.0000","fee":"0.7169","balance":"1.7326"}}}}`)
	queryLedgerJSON     = []byte(`{"error":[],"result":{"LTCH4T-LG5FS-MKGVD1":{"refid":"TYE7IH-QCG76-BVMCM1","time":1570551111.2568,"type":"rollover","aclass":"currency","asset":"ZUSD","amount":"0.0000","fee":"0.4640","balance":"1.3540"}}}`)
	getTradeVolumeJSON  = []byte(`{"error":[],"result":{"currency":"ZUSD","volume":"1000","fees":{"XXBTZUSD":{"fee":"0.1600","minfee":"0.1000","maxfee":"0.2600","nextfee":"0.1400","nextvolume":"2500000.0000","tiervolume":"1000000.0000"}},"fees_maker":{"XXBTZUSD":{"fee":"0.0600","minfee":"0.0000","maxfee":"0.1600","nextfee":"0.0400","nextvolume":"2500000.0000","tiervolume":"1000000.0000"}}}}`)
	getWSTokenJSON      = []byte(`{"error":[],"result":{"token": "test", "expires": 900}}`)
)

func TestKraken_GetDepositMethods(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    []DepositMethods
		wantErr bool
	}{
		{
			name:    "Kraken returns error",
			err:     ErrSomething,
			resp:    &http.Response{},
			want:    []DepositMethods{},
			wantErr: true,
		}, {
			name: "Get Deposit Methods",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(depositMethodsJSON)),
			},
			want:    []DepositMethods{{Method: "Ether (Hex)", Limit: false, Fee: "0.0000000000", GenAddress: true}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				client: &httpMock{
					Error:    tt.err,
					Response: tt.resp,
				},
			}
			got, err := api.GetDepositMethods()
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.GetDepositMethods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.GetDepositMethods() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_GetDepositStatuses(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    []DepositStatuses
		wantErr bool
	}{
		{
			name:    "Kraken returns error",
			err:     ErrSomething,
			resp:    &http.Response{},
			want:    []DepositStatuses{},
			wantErr: true,
		}, {
			name: "Get Deposit Statuses",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(depositStatusesJSON)),
			},
			want: []DepositStatuses{{Method: "Ether (Hex)", Aclass: "currency", Asset: "XETH", Refid: "sometest1",
				Txid: "sometest2", Info: "sometest3", Amount: "6.91", Fee: "0.0000000000", Time: 1617014556, Status: "Success"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				client: &httpMock{
					Error:    tt.err,
					Response: tt.resp,
				},
			}
			got, err := api.GetDepositStatus("Ether (Hex)", "XETH")
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.GetDepositMethods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.GetDepositMethods() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_GetAccountBalances(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    map[string]decimal.Decimal
		wantErr bool
	}{
		{
			name:    "Kraken returns error",
			err:     ErrSomething,
			resp:    &http.Response{},
			want:    make(map[string]decimal.Decimal),
			wantErr: true,
		}, {
			name: "Get Account Balances",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(balancesJSON)),
			},
			want: map[string]decimal.Decimal{
				"BSV":  decimal.NewFromFloat(0.0000053898),
				"ZUSD": decimal.NewFromFloat(435.9135),
				"USDT": decimal.NewFromInt(2),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				client: &httpMock{
					Error:    tt.err,
					Response: tt.resp,
				},
			}
			got, err := api.GetAccountBalances()
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.GetAccountBalances() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Len(t, got, len(tt.want))
			for name, balance := range got {
				wantBalance, ok := tt.want[name]
				if !ok {
					t.Errorf("Kraken.GetAccountBalances() unknown asset: %s", name)
					return
				}
				assert.Equal(t, wantBalance.String(), balance.String())
			}
		})
	}
}

func TestKraken_GetAccountBalancesEx(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    map[string]BalanceEx
		wantErr bool
	}{
		{
			name:    "Kraken returns error",
			err:     ErrSomething,
			resp:    &http.Response{},
			want:    make(map[string]BalanceEx),
			wantErr: true,
		}, {
			name: "Get Account Balances Ex",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(balancesExJSON)),
			},
			want: map[string]BalanceEx{
				"ZUSD": {
					Balance:   decimal.NewFromFloat(25435.21),
					HoldTrade: decimal.NewFromFloat(8249.76),
				},
				"XXBT": {
					Balance:   decimal.NewFromFloat(1.2435),
					HoldTrade: decimal.NewFromFloat(0.8423),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				client: &httpMock{
					Error:    tt.err,
					Response: tt.resp,
				},
			}
			got, err := api.GetAccountBalancesEx()
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.GetAccountBalances() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Len(t, got, len(tt.want))
			for name, balanceEx := range got {
				wantBalanceEx, ok := tt.want[name]
				if !ok {
					t.Errorf("Kraken.GetAccountBalances() unknown asset: %s", name)
					return
				}
				assert.Equal(t, wantBalanceEx.Balance, balanceEx.Balance)
			}
		})
	}
}

func TestKraken_GetTradeBalance(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    TradeBalanceResponse
		wantErr bool
	}{
		{
			name:    "Kraken returns error",
			err:     ErrSomething,
			resp:    &http.Response{},
			want:    TradeBalanceResponse{},
			wantErr: true,
		}, {
			name: "Get Account Trade Balances",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(tradeBalancesJSON)),
			},
			want: TradeBalanceResponse{
				EquivalentBalance: 33.5,
				TradeBalance:      33.5,
				OpenMargin:        23.77,
				UnrealizedProfit:  4.375,
				CostPositions:     11.8999,
				CurrentValue:      12.2,
				Equity:            32.1,
				FreeMargin:        33.1,
				MarginLevel:       12.97,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				client: &httpMock{
					Error:    tt.err,
					Response: tt.resp,
				},
			}
			got, err := api.GetTradeBalance("")
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.GetTradeBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.GetTradeBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_GetOpenOrders(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    OpenOrdersResponse
		wantErr bool
	}{
		{
			name:    "Kraken returns error",
			err:     ErrSomething,
			resp:    &http.Response{},
			want:    OpenOrdersResponse{},
			wantErr: true,
		}, {
			name: "Get Opened Orders",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(openOrdersJSON)),
			},
			want: OpenOrdersResponse{
				Orders: map[string]OrderInfo{
					"OR3XZM-5EN2R-LS5X51": {
						RefID:           nil,
						UserRef:         nil,
						Status:          "open",
						OpenTimestamp:   1570622342.3552,
						StartTimestamp:  0,
						ExpireTimestamp: 0,
						Volume:          1.1,
						VolumeExecuted:  0,
						Cost:            0,
						Fee:             0,
						AveragePrice:    0,
						StopPrice:       0,
						LimitPrice:      0,
						Misc:            "",
						Flags:           "fciq",
						Description: OrderDescription{
							Pair:           "XBTEUR",
							Side:           "sell",
							OrderType:      "limit",
							Price:          7712.2,
							Price2:         0,
							Leverage:       "4:1",
							Info:           "sell 1.10000000 XBTEUR @ limit 7712.2 with 4:1 leverage",
							CloseCondition: "",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				client: &httpMock{
					Error:    tt.err,
					Response: tt.resp,
				},
			}
			got, err := api.GetOpenOrders(false, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.GetOpenOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.GetOpenOrders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_GetClosedOrders(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    ClosedOrdersResponse
		wantErr bool
	}{
		{
			name:    "Kraken returns error",
			err:     ErrSomething,
			resp:    &http.Response{},
			want:    ClosedOrdersResponse{},
			wantErr: true,
		}, {
			name: "Get Closed orders",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(closedOrdersJSON)),
			},
			want: ClosedOrdersResponse{
				Count: 20,
				Orders: map[string]OrderInfo{
					"OK46ER-A2BXK-YOLKE1": {
						RefID:           nil,
						UserRef:         nil,
						Status:          "canceled",
						Reason:          "User requested",
						OpenTimestamp:   1570623817.6537,
						CloseTimestamp:  1570623823.9012,
						StartTimestamp:  0,
						ExpireTimestamp: 0,
						Volume:          21,
						VolumeExecuted:  0,
						Cost:            0,
						Fee:             0,
						AveragePrice:    0,
						StopPrice:       0,
						LimitPrice:      0,
						Misc:            "",
						Flags:           "fciq",
						Description: OrderDescription{
							Pair:           "ETHEUR",
							Side:           "buy",
							OrderType:      "limit",
							Price:          160.87,
							Price2:         0,
							Leverage:       "4:1",
							Info:           "buy 21.00000000 ETHEUR @ limit 160.87 with 4:1 leverage",
							CloseCondition: "",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				client: &httpMock{
					Error:    tt.err,
					Response: tt.resp,
				},
			}
			got, err := api.GetClosedOrders(false, "", 0, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.GetClosedOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.GetClosedOrders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_QueryOrders(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    map[string]OrderInfo
		wantErr bool
	}{
		{
			name:    "Kraken returns error",
			err:     ErrSomething,
			resp:    &http.Response{},
			want:    nil,
			wantErr: true,
		}, {
			name: "Query orders",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(queryOrdersJSON)),
			},
			want: map[string]OrderInfo{
				"OLNYE1-H3BBJ-JD2LGC": {
					RefID:           nil,
					UserRef:         nil,
					Status:          "canceled",
					Reason:          "User requested",
					OpenTimestamp:   1570623816.1101,
					CloseTimestamp:  1570623819.639,
					StartTimestamp:  0,
					ExpireTimestamp: 0,
					Volume:          1.1,
					VolumeExecuted:  0,
					Cost:            0,
					Fee:             0,
					AveragePrice:    0,
					StopPrice:       0,
					LimitPrice:      0,
					Misc:            "",
					Flags:           "fciq",
					Description: OrderDescription{
						Pair:           "XBTUSD",
						Side:           "buy",
						OrderType:      "limit",
						Price:          7920.9,
						Price2:         0,
						Leverage:       "4:1",
						Info:           "buy 1.10000000 XBTUSD @ limit 7920.9 with 4:1 leverage",
						CloseCondition: "",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				client: &httpMock{
					Error:    tt.err,
					Response: tt.resp,
				},
			}
			got, err := api.QueryOrders(false, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.QueryOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.QueryOrders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_GetTradesHistory(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    TradesHistoryResponse
		wantErr bool
	}{
		{
			name:    "Kraken returns error",
			err:     ErrSomething,
			resp:    &http.Response{},
			want:    TradesHistoryResponse{},
			wantErr: true,
		}, {
			name: "Get trades history",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(tradeHistoryJSON)),
			},
			want: TradesHistoryResponse{
				Count: 1,
				Trades: map[string]PrivateTrade{
					"TO3MMA-BSBGV-XUV4A1": {
						OrderID:    "OSQQQ5-MBKL6-O4YYE1",
						PositionID: "TYE7IH-QCG76-BVMCM1",
						Pair:       "XXBTZUSD",
						Time:       1570477513.2,
						Side:       "buy",
						OrderType:  "limit",
						Price:      7000.6,
						Cost:       1000.38301,
						Fee:        0,
						Volume:     0.2,
						Margin:     320,
						Misc:       "closing",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				client: &httpMock{
					Error:    tt.err,
					Response: tt.resp,
				},
			}
			got, err := api.GetTradesHistory("", false, 0, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.GetTradesHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.GetTradesHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_QueryTrades(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    map[string]PrivateTrade
		wantErr bool
	}{
		{
			name:    "Kraken returns error",
			err:     ErrSomething,
			resp:    &http.Response{},
			want:    nil,
			wantErr: true,
		}, {
			name: "Query trades",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(queryTradesJSON)),
			},
			want: map[string]PrivateTrade{
				"TO3MMA-BSBGV-XUV4A1": {
					OrderID:    "OSQQQ5-MBKL6-O4YYE1",
					PositionID: "TYE7IH-QCG76-BVMCM1",
					Pair:       "XXBTZUSD",
					Time:       1570477513.2,
					Side:       "buy",
					OrderType:  "limit",
					Price:      7000.6,
					Cost:       1000.38301,
					Fee:        0,
					Volume:     0.2,
					Margin:     320,
					Misc:       "closing",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				client: &httpMock{
					Error:    tt.err,
					Response: tt.resp,
				},
			}
			got, err := api.QueryTrades(false, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.QueryTrades() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.QueryTrades() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_GetOpenPositions(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    map[string]Position
		wantErr bool
	}{
		{
			name:    "Kraken returns error",
			err:     ErrSomething,
			resp:    &http.Response{},
			want:    nil,
			wantErr: true,
		}, {
			name: "get open positions",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(openPositionsJSON)),
			},
			want: map[string]Position{
				"TYE7IH-QCG76-BVMCM1": {
					OrderID:      "OK7SOC-SGF3O-F54S51",
					Status:       "open",
					Pair:         "XXBTZUSD",
					Time:         1569513333.0361,
					Side:         "buy",
					OrderType:    "limit",
					Cost:         570.39712,
					Fee:          39,
					Volume:       7,
					VolumeClosed: 6.66208817,
					Margin:       9.2,
					Terms:        "0.0100% per 4 hours",
					RolloverTime: 1570638129,
					Misc:         "",
					Flags:        "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				client: &httpMock{
					Error:    tt.err,
					Response: tt.resp,
				},
			}
			got, err := api.GetOpenPositions(false, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.GetOpenPositions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.GetOpenPositions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_GetLedgersInfo(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    LedgerInfoResponse
		wantErr bool
	}{
		{
			name:    "Kraken returns error",
			err:     ErrSomething,
			resp:    &http.Response{},
			want:    LedgerInfoResponse{},
			wantErr: true,
		}, {
			name: "Get ledgers info",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(getLedgersJSON)),
			},
			want: LedgerInfoResponse{
				Ledgers: map[string]Ledger{
					"LGPNZQ-2SLSA-C7QCT1": {
						RefID:      "TI2NBU-IICD2-BAVYO1",
						Time:       1570623111.9096,
						LedgerType: "rollover",
						AssetClass: "currency",
						Asset:      "ZUSD",
						Amount:     0,
						Fee:        0.7169,
						Balance:    1.7326,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				client: &httpMock{
					Error:    tt.err,
					Response: tt.resp,
				},
			}
			got, err := api.GetLedgersInfo("", 0, 0, "ZUSD")
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.GetLedgersInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.GetLedgersInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_QueryLedgers(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    map[string]Ledger
		wantErr bool
	}{
		{
			name:    "Kraken returns error",
			err:     ErrSomething,
			resp:    &http.Response{},
			want:    nil,
			wantErr: true,
		}, {
			name: "Query ledgers",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(queryLedgerJSON)),
			},
			want: map[string]Ledger{
				"LTCH4T-LG5FS-MKGVD1": {
					RefID:      "TYE7IH-QCG76-BVMCM1",
					Time:       1570551111.2568,
					LedgerType: "rollover",
					AssetClass: "currency",
					Asset:      "ZUSD",
					Amount:     0,
					Fee:        0.4640,
					Balance:    1.3540,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				client: &httpMock{
					Error:    tt.err,
					Response: tt.resp,
				},
			}
			got, err := api.QueryLedgers("")
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.QueryLedgers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.QueryLedgers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_GetTradeVolume(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    TradeVolumeResponse
		wantErr bool
	}{
		{
			name:    "Kraken returns error",
			err:     ErrSomething,
			resp:    &http.Response{},
			want:    TradeVolumeResponse{},
			wantErr: true,
		}, {
			name: "Get Trade Volume",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(getTradeVolumeJSON)),
			},
			want: TradeVolumeResponse{
				Currency: "ZUSD",
				Volume:   1000,
				Fees: map[string]Fees{
					"XXBTZUSD": {
						Fee:        0.16,
						MinFee:     0.1,
						MaxFee:     0.26,
						NextFee:    0.14,
						NextVolume: 2500000,
						TierVolume: 1000000,
					},
				},
				FeesMaker: map[string]Fees{
					"XXBTZUSD": {
						Fee:        0.06,
						MinFee:     0,
						MaxFee:     0.16,
						NextFee:    0.04,
						NextVolume: 2500000,
						TierVolume: 1000000,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				client: &httpMock{
					Error:    tt.err,
					Response: tt.resp,
				},
			}
			got, err := api.GetTradeVolume(false, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.GetTradeVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.GetTradeVolume() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_GetWebSocketsToken(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    GetWebSocketTokenResponse
		wantErr bool
	}{
		{
			name:    "Kraken returns error",
			err:     ErrSomething,
			resp:    &http.Response{},
			want:    GetWebSocketTokenResponse{},
			wantErr: true,
		}, {
			name: "Get WS Token",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(getWSTokenJSON)),
			},
			want: GetWebSocketTokenResponse{
				Token:   "test",
				Expires: 900,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				client: &httpMock{
					Error:    tt.err,
					Response: tt.resp,
				},
			}
			got, err := api.GetWebSocketsToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.GetWebSocketsToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.GetWebSocketsToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
