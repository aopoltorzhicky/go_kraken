package rest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

var ErrSomething = fmt.Errorf("Something went wrong")

type httpMock struct {
	Response *http.Response
	Error    error
}

func (c *httpMock) Do(req *http.Request) (*http.Response, error) {
	if c.Error != nil {
		return c.Response, c.Error
	}
	return c.Response, nil
}

func TestKraken_Time(t *testing.T) {
	json := []byte(`{"error":[],"result":{"unixtime":1554218108,"rfc1123":"Tue,  2 Apr 19 15:15:08 +0000"}}`)
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		want    TimeResponse
		wantErr bool
	}{
		{
			name:    "Error returned from Kraken",
			err:     ErrSomething,
			resp:    &http.Response{},
			want:    TimeResponse{},
			wantErr: true,
		},
		{
			name: "Data returned from Kraken",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(json)),
			},
			want: TimeResponse{
				Unixtime: 1554218108,
				Rfc1123:  "Tue,  2 Apr 19 15:15:08 +0000",
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
			got, err := api.Time()
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.Time() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_Assets(t *testing.T) {
	json := []byte(`{"error":[],"result":{"ADA":{"aclass":"currency","altname":"ADA","decimals":8,"display_decimals":6}}}`)
	type args struct {
		assets []string
	}
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		args    args
		want    map[string]Asset
		wantErr bool
	}{
		{
			name: "Kraken returns Error",
			err:  ErrSomething,
			resp: &http.Response{},
			args: args{
				assets: nil,
			},
			want:    map[string]Asset{},
			wantErr: true,
		},
		{
			name: "Get all from kraken",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(json)),
			},
			args: args{
				assets: nil,
			},
			want: map[string]Asset{
				"ADA": Asset{
					AlternateName:   "ADA",
					AssetClass:      "currency",
					Decimals:        8,
					DisplayDecimals: 6,
				},
			},
			wantErr: false,
		},
		{
			name: "Get one asset from kraken",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(json)),
			},
			args: args{
				assets: []string{"ADA"},
			},
			want: map[string]Asset{
				"ADA": Asset{
					AlternateName:   "ADA",
					AssetClass:      "currency",
					Decimals:        8,
					DisplayDecimals: 6,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				client: &httpMock{
					Response: tt.resp,
					Error:    tt.err,
				},
			}
			got, err := api.Assets(tt.args.assets...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.Assets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.Assets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_AssetPairs(t *testing.T) {
	json := []byte(`{"error":[],"result":{"ADACAD":{"altname":"ADACAD","wsname":"ADA\/CAD","aclass_base":"currency","base":"ADA","aclass_quote":"currency","quote":"ZCAD","lot":"unit","pair_decimals":6,"lot_decimals":8,"lot_multiplier":1,"leverage_buy":[],"leverage_sell":[],"fees":[[0,0.26],[50000,0.24],[100000,0.22],[250000,0.2],[500000,0.18],[1000000,0.16],[2500000,0.14],[5000000,0.12],[10000000,0.1]],"fees_maker":[[0,0.16],[50000,0.14],[100000,0.12],[250000,0.1],[500000,0.08],[1000000,0.06],[2500000,0.04],[5000000,0.02],[10000000,0]],"fee_volume_currency":"ZUSD","margin_call":80,"margin_stop":40}}}`)
	type args struct {
		pairs []string
	}
	tests := []struct {
		name    string
		args    args
		err     error
		resp    *http.Response
		want    map[string]AssetPair
		wantErr bool
	}{
		{
			name: "Kraken returns Error",
			err:  ErrSomething,
			resp: &http.Response{},
			args: args{
				pairs: []string{"ADACAD"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Get asset pairs",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(json))},
			args: args{
				pairs: []string{"ADACAD"},
			},
			want: map[string]AssetPair{
				"ADACAD": AssetPair{
					Altname:           "ADACAD",
					WSName:            "ADA/CAD",
					AssetClassBase:    "currency",
					Base:              "ADA",
					AssetClassQuote:   "currency",
					Quote:             "ZCAD",
					Lot:               "unit",
					PairDecimals:      6,
					LotDecimals:       8,
					LotMultiplier:     1,
					LeverageBuy:       []float64{},
					LeverageSell:      []float64{},
					Fees:              [][]float64{{0, 0.26}, {50000, 0.24}, {100000, 0.22}, {250000, 0.2}, {500000, 0.18}, {1000000, 0.16}, {2500000, 0.14}, {5000000, 0.12}, {10000000, 0.1}},
					FeesMaker:         [][]float64{{0, 0.16}, {50000, 0.14}, {100000, 0.12}, {250000, 0.1}, {500000, 0.08}, {1000000, 0.06}, {2500000, 0.04}, {5000000, 0.02}, {10000000, 0}},
					FeeVolumeCurrency: "ZUSD",
					MarginCall:        80,
					MarginStop:        40,
				},
			},
			wantErr: false,
		},
		{
			name: "Pairs is nil",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(json))},
			args: args{
				pairs: nil,
			},
			want: map[string]AssetPair{
				"ADACAD": AssetPair{
					Altname:           "ADACAD",
					WSName:            "ADA/CAD",
					AssetClassBase:    "currency",
					Base:              "ADA",
					AssetClassQuote:   "currency",
					Quote:             "ZCAD",
					Lot:               "unit",
					PairDecimals:      6,
					LotDecimals:       8,
					LotMultiplier:     1,
					LeverageBuy:       []float64{},
					LeverageSell:      []float64{},
					Fees:              [][]float64{{0, 0.26}, {50000, 0.24}, {100000, 0.22}, {250000, 0.2}, {500000, 0.18}, {1000000, 0.16}, {2500000, 0.14}, {5000000, 0.12}, {10000000, 0.1}},
					FeesMaker:         [][]float64{{0, 0.16}, {50000, 0.14}, {100000, 0.12}, {250000, 0.1}, {500000, 0.08}, {1000000, 0.06}, {2500000, 0.04}, {5000000, 0.02}, {10000000, 0}},
					FeeVolumeCurrency: "ZUSD",
					MarginCall:        80,
					MarginStop:        40,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				client: &httpMock{
					Response: tt.resp,
					Error:    tt.err,
				},
			}
			got, err := api.AssetPairs(tt.args.pairs...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.AssetPairs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.AssetPairs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_Ticker(t *testing.T) {
	json := []byte(`{"error":[],"result":{"ADACAD":{"a":["0.108312","6418","6418.000"],"b":["0.090125","2688","2688.000"],"c":["0.090043","0.00000091"],"v":["115805.23341809","136512.79974015"],"p":["0.102010","0.100786"],"t":[54,67],"l":["0.090000","0.090000"],"h":["0.109000","0.109000"],"o":"0.093911"}}}`)
	type args struct {
		pairs []string
	}
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		args    args
		want    map[string]Ticker
		wantErr bool
	}{
		{
			name: "Kraken returns error",
			err:  ErrSomething,
			resp: &http.Response{},
			args: args{
				pairs: []string{"ADACAD"},
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "No pairs",
			err:  ErrSomething,
			resp: &http.Response{},
			args: args{
				pairs: nil,
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "Get ticker",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(json)),
			},
			args: args{
				pairs: []string{"ADACAD"},
			},
			want: map[string]Ticker{
				"ADACAD": Ticker{
					Ask: Level{
						Price:          0.108312,
						WholeLotVolume: 6418.,
						Volume:         6418.000,
					},
					Bid: Level{
						Price:          0.090125,
						WholeLotVolume: 2688.,
						Volume:         2688.000,
					},
					Close: CloseLevel{
						Price:     0.090043,
						LotVolume: 0.00000091,
					},
					Volume: CloseLevel{
						Price:     115805.23341809,
						LotVolume: 136512.79974015,
					},
					VolumeAveragePrice: CloseLevel{
						Price:     0.102010,
						LotVolume: 0.100786,
					},
					Trades: TimeLevel{
						Today:       54.,
						Last24Hours: 67.,
					},
					Low: CloseLevel{
						Price:     0.090000,
						LotVolume: 0.090000,
					},
					High: CloseLevel{
						Price:     0.109000,
						LotVolume: 0.109000,
					},
					OpeningPrice: 0.093911,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				client: &httpMock{
					Response: tt.resp,
					Error:    tt.err,
				},
			}
			got, err := api.Ticker(tt.args.pairs...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.Ticker() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.Ticker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_Candles(t *testing.T) {
	json := []byte(`{"error":[],"result":{"ADACAD":[[1554179640,"0.0005000","0.0005000","0.0005000","0.0005000","0.0000000","0.00000000",0]],"last":1554222360}}`)
	response := OHLCResponse{
		Last: 1554222360,
		Candles: map[string][]Candle{
			"ADACAD": []Candle{{
				Time:      1554179640,
				Open:      0.0005000,
				High:      0.0005000,
				Low:       0.0005000,
				Close:     0.0005000,
				VolumeWAP: 0.0000000,
				Volume:    0.0000000,
				Count:     0,
			}},
		},
	}
	type args struct {
		pair     string
		interval int64
		since    int64
	}
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		args    args
		want    OHLCResponse
		wantErr bool
	}{

		{
			name: "Kraken returns error",
			err:  ErrSomething,
			resp: &http.Response{},
			args: args{
				pair:     "ADACAD",
				interval: 0,
				since:    0,
			},
			want:    OHLCResponse{},
			wantErr: true,
		}, {
			name: "Get candles from kraken",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(json)),
			},
			args: args{
				pair:     "ADACAD",
				interval: 0,
				since:    0,
			},
			want:    response,
			wantErr: false,
		}, {
			name: "Get candles from kraken with interval and since",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(json)),
			},
			args: args{
				pair:     "ADACAD",
				interval: Interval15m,
				since:    123,
			},
			want:    response,
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
			got, err := api.Candles(tt.args.pair, tt.args.interval, tt.args.since)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.Candles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.Candles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_GetOrderBook(t *testing.T) {
	json := []byte(`{"error":[],"result":{"ADACAD":{"asks":[["0.109441","6741.072",1554223624],["0.109442","4950.724",1554223614]],"bids":[["0.090494","2789.652",1554223622],["0.090493","6379.886",1554223620]]}}}`)
	type args struct {
		pair  string
		depth int64
	}
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		args    args
		want    map[string]OrderBook
		wantErr bool
	}{
		{
			name: "Kraken returns error",
			err:  ErrSomething,
			resp: &http.Response{},
			args: args{
				pair:  "ADACAD",
				depth: 2,
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "get order book",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(json)),
			},
			args: args{
				pair:  "ADACAD",
				depth: 2,
			},
			want: map[string]OrderBook{
				"ADACAD": OrderBook{
					Asks: []OrderBookItem{
						OrderBookItem{
							Price:     0.109441,
							Volume:    6741.072,
							Timestamp: 1554223624,
						},
						OrderBookItem{
							Price:     0.109442,
							Volume:    4950.724,
							Timestamp: 1554223614,
						},
					},
					Bids: []OrderBookItem{
						OrderBookItem{
							Price:     0.090494,
							Volume:    2789.652,
							Timestamp: 1554223622,
						},
						OrderBookItem{
							Price:     0.090493,
							Volume:    6379.886,
							Timestamp: 1554223620,
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
			got, err := api.GetOrderBook(tt.args.pair, tt.args.depth)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.GetOrderBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.GetOrderBook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_GetTrades(t *testing.T) {
	json := []byte(`{"error":[],"result":{"ADACAD":[["0.093280","2968.26413227",1553959154.2509,"s","l",""]], "last": 1554221914617956627}}`)
	type args struct {
		pair  string
		since int64
	}
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		args    args
		want    TradeResponse
		wantErr bool
	}{
		{
			name: "Kraken returns error",
			err:  ErrSomething,
			resp: &http.Response{},
			args: args{
				pair:  "ADACAD",
				since: 2,
			},
			want:    TradeResponse{},
			wantErr: true,
		}, {
			name: "Get trades",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(json)),
			},
			args: args{
				pair:  "ADACAD",
				since: 2,
			},
			want: TradeResponse{
				Last: 1554221914617956627,
				ADACAD: []Trade{
					Trade{
						Price:     0.093280,
						Volume:    2968.26413227,
						Time:      1553959154.2509,
						Side:      "s",
						OrderType: "l",
						Misc:      "",
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
			got, err := api.GetTrades(tt.args.pair, tt.args.since)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.GetTrades() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.GetTrades() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_GetSpread(t *testing.T) {
	json := []byte(`{"error":[],"result":{"ADACAD":[[1554224145,"0.091118","0.109331"]], "last":1554224725 }}`)
	type args struct {
		pair  string
		since int64
	}
	tests := []struct {
		name    string
		err     error
		resp    *http.Response
		args    args
		want    SpreadResponse
		wantErr bool
	}{
		{
			name: "Kraken returns error",
			err:  ErrSomething,
			resp: &http.Response{},
			args: args{
				pair:  "ADACAD",
				since: 2,
			},
			want:    SpreadResponse{},
			wantErr: true,
		}, {
			name: "Get spread",
			err:  nil,
			resp: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(json)),
			},
			args: args{
				pair:  "ADACAD",
				since: 2,
			},
			want: SpreadResponse{
				Last: 1554224725,
				ADACAD: []Spread{
					Spread{
						Time: 1554224145,
						Ask:  0.109331,
						Bid:  0.091118,
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
			got, err := api.GetSpread(tt.args.pair, tt.args.since)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.GetSpread() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.GetSpread() = %v, want %v", got, tt.want)
			}
		})
	}
}
