package rest

import (
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_getTimestamp(t *testing.T) {
	tests := []struct {
		name    string
		args    interface{}
		want    int64
		wantErr bool
	}{
		{
			name:    "invalid type",
			args:    "this is string",
			want:    0,
			wantErr: true,
		}, {
			name:    "good type",
			args:    1523.3,
			want:    1523,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTimestamp(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFloat64(t *testing.T) {
	tests := []struct {
		name    string
		args    interface{}
		want    float64
		wantErr bool
	}{
		{
			name:    "invalid type",
			args:    "this is string",
			want:    0,
			wantErr: true,
		}, {
			name:    "good type",
			args:    1523.3,
			want:    1523.3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFloat64(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFloat64FromStr(t *testing.T) {
	tests := []struct {
		name    string
		args    interface{}
		want    float64
		wantErr bool
	}{
		{
			name:    "invalid type",
			args:    123,
			want:    0,
			wantErr: true,
		}, {
			name:    "good type - invalid text",
			args:    "text",
			want:    0,
			wantErr: true,
		}, {
			name:    "good",
			args:    "123.3",
			want:    123.3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFloat64FromStr(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFloat64FromStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getFloat64FromStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevel_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		buf     []byte
		wantErr bool
		result  *Level
	}{
		{
			name:    "invalid json",
			buf:     []byte(``),
			wantErr: true,
			result:  &Level{},
		}, {
			name:    "invalid array length",
			buf:     []byte(`[]`),
			wantErr: true,
			result:  &Level{},
		}, {
			buf:     []byte(`[123, 123, 123]`),
			wantErr: false,
			result: &Level{
				Price:          decimal.NewFromInt(123),
				WholeLotVolume: decimal.NewFromInt(123),
				Volume:         decimal.NewFromInt(123),
			},
		}, {
			buf:     []byte(`["123.0", 123, 123]`),
			wantErr: false,
			result: &Level{
				Price:          decimal.NewFromInt(123),
				WholeLotVolume: decimal.NewFromInt(123),
				Volume:         decimal.NewFromInt(123),
			},
		}, {
			buf:     []byte(`["123.0", "124.0", 123]`),
			wantErr: false,
			result: &Level{
				Price:          decimal.NewFromInt(123),
				WholeLotVolume: decimal.NewFromInt(124),
				Volume:         decimal.NewFromInt(123),
			},
		}, {
			buf:     []byte(`["123.0", "124.0", "125.0"]`),
			wantErr: false,
			result: &Level{
				Price:          decimal.NewFromInt(123),
				WholeLotVolume: decimal.NewFromInt(124),
				Volume:         decimal.NewFromInt(125),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &Level{}
			if err := item.UnmarshalJSON(tt.buf); (err != nil) != tt.wantErr {
				t.Errorf("Level.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, item.Price.String(), tt.result.Price.String())
			assert.Equal(t, item.WholeLotVolume.String(), tt.result.WholeLotVolume.String())
			assert.Equal(t, item.Volume.String(), tt.result.Volume.String())
		})
	}
}

func TestTimeLevel_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		buf     []byte
		wantErr bool
		result  *TimeLevel
	}{
		{
			name:    "invalid json",
			buf:     []byte(``),
			wantErr: true,
			result:  &TimeLevel{},
		}, {
			name:    "invalid array length",
			buf:     []byte(`[]`),
			wantErr: true,
			result:  &TimeLevel{},
		}, {
			name:    "invalid both types",
			buf:     []byte(`["123", "123"]`),
			wantErr: true,
			result:  &TimeLevel{},
		}, {
			name:    "invalid last 24 hour type",
			buf:     []byte(`[123, "123"]`),
			wantErr: true,
			result:  &TimeLevel{},
		}, {
			name:    "good",
			buf:     []byte(`[123, 124]`),
			wantErr: false,
			result: &TimeLevel{
				Today:       123,
				Last24Hours: 124,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := new(TimeLevel)
			if err := item.UnmarshalJSON(tt.buf); (err != nil) != tt.wantErr {
				t.Errorf("TimeLevel.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, item.Today, tt.result.Today)
			assert.Equal(t, item.Last24Hours, tt.result.Last24Hours)
		})
	}
}

func TestCloseLevel_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		buf     []byte
		wantErr bool
		result  *CloseLevel
	}{
		{
			name:    "invalid json",
			buf:     []byte(``),
			wantErr: true,
			result:  &CloseLevel{},
		}, {
			name:    "invalid array length",
			buf:     []byte(`[]`),
			wantErr: true,
			result:  &CloseLevel{},
		}, {
			buf:     []byte(`[123, 123]`),
			wantErr: false,
			result: &CloseLevel{
				Price:     decimal.NewFromInt(123),
				LotVolume: decimal.NewFromInt(123),
			},
		}, {
			buf:     []byte(`["123.0", 123]`),
			wantErr: false,
			result: &CloseLevel{
				Price:     decimal.NewFromInt(123),
				LotVolume: decimal.NewFromInt(123),
			},
		}, {
			buf:     []byte(`["123.0", "124.0"]`),
			wantErr: false,
			result: &CloseLevel{
				Price:     decimal.NewFromInt(123),
				LotVolume: decimal.NewFromInt(124),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := new(CloseLevel)
			if err := item.UnmarshalJSON(tt.buf); (err != nil) != tt.wantErr {
				t.Errorf("TimeLevel.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, item.Price.String(), tt.result.Price.String())
			assert.Equal(t, item.LotVolume.String(), tt.result.LotVolume.String())
		})
	}
}

func TestOHLCResponse_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Candles map[string][]Candle
		Last    int64
	}
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &OHLCResponse{
				Candles: tt.fields.Candles,
				Last:    tt.fields.Last,
			}
			if err := item.UnmarshalJSON(tt.args.buf); (err != nil) != tt.wantErr {
				t.Errorf("OHLCResponse.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrderBookItem_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		buf     []byte
		wantErr bool
		result  *OrderBookItem
	}{
		{
			name:    "invalid json",
			buf:     []byte(``),
			wantErr: true,
			result:  &OrderBookItem{},
		}, {
			name:    "invalid array length",
			buf:     []byte(`[]`),
			wantErr: true,
			result:  &OrderBookItem{},
		}, {
			name:    "invalid price",
			buf:     []byte(`[123, 123, 123]`),
			wantErr: true,
			result:  &OrderBookItem{},
		}, {
			name:    "invalid volume",
			buf:     []byte(`["123.0", 123, 123]`),
			wantErr: true,
			result: &OrderBookItem{
				Price: 123,
			},
		}, {
			name:    "invalid timestamp",
			buf:     []byte(`["123.0", "124.0", "123"]`),
			wantErr: true,
			result: &OrderBookItem{
				Price:  123,
				Volume: 124,
			},
		}, {
			name:    "good",
			buf:     []byte(`["123.0", "124.0", 125.0]`),
			wantErr: false,
			result: &OrderBookItem{
				Price:     123,
				Volume:    124,
				Timestamp: 125,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &OrderBookItem{}
			if err := item.UnmarshalJSON(tt.buf); (err != nil) != tt.wantErr {
				t.Errorf("Level.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(item, tt.result) {
				t.Errorf("Kraken.GetAccountBalances() = %v, want %v", item, tt.result)
			}
		})
	}
}

func TestTrade_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Price     float64
		Volume    float64
		Time      float64
		Side      string
		OrderType string
		Misc      string
	}
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &Trade{
				Price:     tt.fields.Price,
				Volume:    tt.fields.Volume,
				Time:      tt.fields.Time,
				Side:      tt.fields.Side,
				OrderType: tt.fields.OrderType,
				Misc:      tt.fields.Misc,
			}
			if err := item.UnmarshalJSON(tt.args.buf); (err != nil) != tt.wantErr {
				t.Errorf("Trade.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSpread_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Time float64
		Bid  float64
		Ask  float64
	}
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &Spread{
				Time: tt.fields.Time,
				Bid:  tt.fields.Bid,
				Ask:  tt.fields.Ask,
			}
			if err := item.UnmarshalJSON(tt.args.buf); (err != nil) != tt.wantErr {
				t.Errorf("Spread.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
