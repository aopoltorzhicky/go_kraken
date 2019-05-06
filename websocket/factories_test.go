package websocket

import (
	"log"
	"reflect"
	"testing"
	"time"
)

func Test_newTickerFactory(t *testing.T) {
	tests := []struct {
		name string
		want *tickerFactory
	}{
		{
			name: "Test creation of TickerFactory",
			want: &tickerFactory{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newTickerFactory(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newTickerFactory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tickerFactory_Parse(t *testing.T) {
	type args struct {
		data interface{}
		pair string
	}
	tests := []struct {
		name    string
		f       *tickerFactory
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Test parse correct data",
			f:    &tickerFactory{},
			args: args{
				pair: BTCCAD,
				data: map[string]interface{}{
					"a": []interface{}{"5525.40000", 1, "1.000"},
					"b": []interface{}{"5525.10000", 1, "1.000"},
					"c": []interface{}{"5525.10000", "0.00398963"},
					"v": []interface{}{"2634.11501494", "3591.17907851"},
					"p": []interface{}{"5631.44067", "5653.78939"},
					"t": []interface{}{11493, 16267},
					"l": []interface{}{"5505.00000", "5505.00000"},
					"h": []interface{}{"5783.00000", "5783.00000"},
					"o": []interface{}{"5760.70000", "5763.40000"},
				},
			},
			want: TickerUpdate{
				Pair: BTCCAD,
				Ask: Level{
					Price:          5525.4,
					Volume:         1.,
					WholeLotVolume: 1,
				},
				Bid: Level{
					Price:          5525.1,
					Volume:         1.,
					WholeLotVolume: 1,
				},
				Close: Values{
					Today:  5525.1,
					Last24: 0.00398963,
				},
				Volume: Values{
					Today:  2634.11501494,
					Last24: 3591.17907851,
				},
				VolumeAveragePrice: Values{
					Today:  5631.44067,
					Last24: 5653.78939,
				},
				TradeVolume: Values{
					Today:  11493,
					Last24: 16267,
				},
				Low: Values{
					Today:  5505.00000,
					Last24: 5505.00000,
				},
				High: Values{
					Today:  5783.00000,
					Last24: 5783.00000,
				},
				Open: Values{
					Today:  5760.70000,
					Last24: 5763.40000,
				},
			},
			wantErr: false,
		},
		{
			name: "Test invalid data",
			f:    &tickerFactory{},
			args: args{
				pair: BTCCAD,
				data: nil,
			},
			want:    TickerUpdate{Pair: BTCCAD},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &tickerFactory{}
			got, err := f.Parse(tt.args.data, tt.args.pair)
			if (err != nil) != tt.wantErr {
				t.Errorf("tickerFactory.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tickerFactory.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_candlesFactory_Parse(t *testing.T) {
	t1, err := time.Parse(time.RFC3339, "2018-11-12T21:15:14+00:00")
	if err != nil {
		log.Println(err)
	}
	t2, err := time.Parse(time.RFC3339, "2018-11-12T21:16:00+00:00")
	if err != nil {
		log.Println(err)
	}
	type args struct {
		data interface{}
		pair string
	}
	tests := []struct {
		name    string
		f       *candlesFactory
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Test candles parse",
			f:    &candlesFactory{},
			args: args{
				pair: BTCCAD,
				data: []interface{}{
					"1542057314",
					"1542057360",
					"3586.70000",
					"3586.70000",
					"3586.60000",
					"3586.60000",
					"3586.68894",
					"0.03373000",
					2.,
				},
			},
			want: CandleUpdate{
				Time:      t1.UTC(),
				EndTime:   t2.UTC(),
				Open:      3586.7,
				High:      3586.7,
				Low:       3586.6,
				Close:     3586.6,
				VolumeWAP: 3586.68894,
				Volume:    0.03373000,
				Count:     2,
				Pair:      BTCCAD,
			},
			wantErr: false,
		},
		{
			name: "Test invalid data",
			f:    &candlesFactory{},
			args: args{
				pair: BTCCAD,
				data: nil,
			},
			want:    CandleUpdate{Pair: BTCCAD},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &candlesFactory{}
			got, err := f.Parse(tt.args.data, tt.args.pair)
			if (err != nil) != tt.wantErr {
				t.Errorf("candlesFactory.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("candlesFactory.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tradesFactory_Parse(t *testing.T) {
	t1, err := time.Parse(time.RFC3339, "2018-08-18T17:40:57+00:00")
	if err != nil {
		log.Println(err)
	}
	t2, err := time.Parse(time.RFC3339, "2018-08-18T17:40:57+00:00")
	if err != nil {
		log.Println(err)
	}
	type args struct {
		data interface{}
		pair string
	}
	tests := []struct {
		name    string
		f       *tradesFactory
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Test trades parse",
			args: args{
				pair: BTCCAD,
				data: []interface{}{
					[]interface{}{
						"5541.20000",
						"0.15850568",
						"1534614057",
						"s",
						"l",
						"",
					},
					[]interface{}{
						"6060.00000",
						"0.02455000",
						"1534614057",
						"b",
						"m",
						"",
					},
				},
			},
			f: &tradesFactory{},
			want: []TradeUpdate{
				TradeUpdate{
					Pair:      BTCCAD,
					Time:      t1.UTC(),
					Price:     5541.2,
					Volume:    0.15850568,
					Side:      Sell,
					OrderType: Limit,
					Misc:      "",
				},
				TradeUpdate{
					Pair:      BTCCAD,
					Time:      t2.UTC(),
					Price:     6060.,
					Volume:    0.02455000,
					Side:      Buy,
					OrderType: Market,
					Misc:      "",
				},
			},
			wantErr: false,
		},
		{
			name: "Test invalid data",
			f:    &tradesFactory{},
			args: args{
				pair: BTCCAD,
				data: nil,
			},
			want:    []TradeUpdate{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &tradesFactory{}
			got, err := f.Parse(tt.args.data, tt.args.pair)
			if (err != nil) != tt.wantErr {
				t.Errorf("tradesFactory.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tradesFactory.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_spreadFactory_Parse(t *testing.T) {
	t1, err := time.Parse(time.RFC3339, "2018-11-12T21:14:59+00:00")
	if err != nil {
		log.Println(err)
	}
	type args struct {
		data interface{}
		pair string
	}
	tests := []struct {
		name    string
		f       *spreadFactory
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Test spread update parse",
			f:    &spreadFactory{},
			args: args{
				pair: BTCCAD,
				data: []interface{}{
					"5698.40000",
					"5700.00000",
					"1542057299",
				},
			},
			want: SpreadUpdate{
				Pair: BTCCAD,
				Bid:  5698.4,
				Ask:  5700,
				Time: t1.UTC(),
			},
			wantErr: false,
		},
		{
			name: "Test invalid data",
			f:    &spreadFactory{},
			args: args{
				pair: BTCCAD,
				data: nil,
			},
			want:    SpreadUpdate{Pair: BTCCAD},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &spreadFactory{}
			got, err := f.Parse(tt.args.data, tt.args.pair)
			if (err != nil) != tt.wantErr {
				t.Errorf("spreadFactory.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("spreadFactory.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bookFactory_Parse(t *testing.T) {
	t1, err := time.Parse(time.RFC3339, "2018-08-18T17:44:08+00:00")
	if err != nil {
		log.Println(err)
	}
	t2, err := time.Parse(time.RFC3339, "2018-08-18T17:41:38+00:00")
	if err != nil {
		log.Println(err)
	}

	type args struct {
		data interface{}
		pair string
	}
	tests := []struct {
		name    string
		f       *bookFactory
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Test orderbook snapshot parse",
			f:    &bookFactory{},
			args: args{
				pair: BTCCAD,
				data: map[string]interface{}{
					"as": []interface{}{
						[]interface{}{
							"5541.30000",
							"2.50700000",
							"1534614248",
						},
						[]interface{}{
							"5541.80000",
							"0.33000000",
							"1534614098",
						},
					},
					"bs": []interface{}{
						[]interface{}{
							"5541.20000",
							"1.52900000",
							"1534614248",
						},
						[]interface{}{
							"5539.90000",
							"0.30000000",
							"1534614098",
						},
					},
				},
			},
			want: OrderBookUpdate{
				IsSnapshot: true,
				Pair:       BTCCAD,
				Asks: []OrderBookItem{
					OrderBookItem{
						Price:  5541.3,
						Volume: 2.50700000,
						Time:   t1.UTC(),
					},
					OrderBookItem{
						Price:  5541.8,
						Volume: 0.33000000,
						Time:   t2.UTC(),
					},
				},
				Bids: []OrderBookItem{
					OrderBookItem{
						Price:  5541.2,
						Volume: 1.52900000,
						Time:   t1.UTC(),
					},
					OrderBookItem{
						Price:  5539.9,
						Volume: 0.30000000,
						Time:   t2.UTC(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test orderbook update parse",
			f:    &bookFactory{},
			args: args{
				pair: BTCCAD,
				data: map[string]interface{}{
					"a": []interface{}{
						[]interface{}{
							"5541.30000",
							"2.50700000",
							"1534614248",
						},
					},
					"b": []interface{}{
						[]interface{}{
							"5541.20000",
							"1.52900000",
							"1534614248",
						},
					},
				},
			},
			want: OrderBookUpdate{
				IsSnapshot: false,
				Pair:       BTCCAD,
				Asks: []OrderBookItem{
					OrderBookItem{
						Price:  5541.3,
						Volume: 2.50700000,
						Time:   t1.UTC(),
					},
				},
				Bids: []OrderBookItem{
					OrderBookItem{
						Price:  5541.2,
						Volume: 1.52900000,
						Time:   t1.UTC(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test invalid data",
			f:    &bookFactory{},
			args: args{
				pair: BTCCAD,
				data: nil,
			},
			want:    OrderBookUpdate{Pair: BTCCAD},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &bookFactory{}
			got, err := f.Parse(tt.args.data, tt.args.pair)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookFactory.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookFactory.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newCandlesFactory(t *testing.T) {
	tests := []struct {
		name string
		want *candlesFactory
	}{
		{
			name: "Test creation candles factory",
			want: &candlesFactory{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newCandlesFactory(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCandlesFactory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newTradesFactory(t *testing.T) {
	tests := []struct {
		name string
		want *tradesFactory
	}{
		{
			name: "Test creation trades factory",
			want: &tradesFactory{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newTradesFactory(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newTradesFactory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newSpreadFactory(t *testing.T) {
	tests := []struct {
		name string
		want *spreadFactory
	}{
		{
			name: "Test creation spread factory",
			want: &spreadFactory{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSpreadFactory(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSpreadFactory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newBookFactory(t *testing.T) {
	tests := []struct {
		name string
		want *bookFactory
	}{
		{
			name: "Test creation book factory",
			want: &bookFactory{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newBookFactory(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newBookFactory() = %v, want %v", got, tt.want)
			}
		})
	}
}
