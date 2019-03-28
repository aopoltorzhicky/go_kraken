package websocket

import (
	"reflect"
	"testing"
	"time"
)

func Test_valToFloat64(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Valid string",
			args: args{
				"21",
			},
			want: 21.0,
		},
		{
			name: "Not string",
			args: args{
				1234567890,
			},
			want: .0,
		},
		{
			name: "Bad float string",
			args: args{
				"FAILED STRING",
			},
			want: .0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := valToFloat64(tt.args.value); got != tt.want {
				t.Errorf("valToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_valToTime(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "Valid data",
			args: args{
				"1553516579",
			},
			want: time.Unix(0, 1553516579*1e9).UTC(),
		},
		{
			name: "Not string",
			args: args{
				1553516579,
			},
			want: time.Unix(0, 0).UTC(),
		},
		{
			name: "Invalid string",
			args: args{
				"FAILED STRING",
			},
			want: time.Unix(0, 0).UTC(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := valToTime(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("valToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseLevel(t *testing.T) {
	type args struct {
		data []interface{}
	}
	tests := []struct {
		name string
		args args
		want Level
	}{
		{
			name: "Valid data",
			args: args{
				data: []interface{}{"1.6", 2, "3.01"},
			},
			want: Level{
				Price:          1.6,
				Volume:         3.01,
				WholeLotVolume: 2,
			},
		},
		{
			name: "Invalid strings",
			args: args{
				data: []interface{}{1.6, 2, 3.01},
			},
			want: Level{
				Price:          .0,
				Volume:         .0,
				WholeLotVolume: 2,
			},
		},
		{
			name: "Invalid float",
			args: args{
				data: []interface{}{"1.6", "2.02", "3.01"},
			},
			want: Level{
				Price:          1.6,
				Volume:         3.01,
				WholeLotVolume: .0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseLevel(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseValues(t *testing.T) {
	type args struct {
		data []interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "valid data with string in data[0]",
			args: args{
				data: []interface{}{"123", "321"},
			},
			want: Values{
				Today:  123.,
				Last24: 321.,
			},
		},
		{
			name: "valid data with int in data[0]",
			args: args{
				data: []interface{}{123, 321},
			},
			want: Values{
				Today:  123,
				Last24: 321,
			},
		},
		{
			name: "valid data with float in data[0]",
			args: args{
				data: []interface{}{123., 321.},
			},
			want: Values{
				Today:  123.,
				Last24: 321.,
			},
		},
		{
			name: "invalid data with string in data[0]",
			args: args{
				data: []interface{}{"123.", 321.},
			},
			want: Values{
				Today:  123.,
				Last24: 0.,
			},
		},
		{
			name: "invalid data with int in data[0]",
			args: args{
				data: []interface{}{123, "321."},
			},
			want: Values{
				Today:  123,
				Last24: 0,
			},
		},
		{
			name: "invalid data with float in data[0]",
			args: args{
				data: []interface{}{123., "321."},
			},
			want: Values{
				Today:  123.,
				Last24: 0.,
			},
		},
		{
			name: "invalid data in data[0]",
			args: args{
				data: []interface{}{args{}, args{}},
			},
			want: Values{
				Today:  0,
				Last24: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseValues(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseSide(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Valid data: sell",
			args: args{"s"},
			want: "sell",
		},
		{
			name: "Valid data: buy",
			args: args{"b"},
			want: "buy",
		},
		{
			name: "Invalid data",
			args: args{"invalid data"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseSide(tt.args.data); got != tt.want {
				t.Errorf("parseSide() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseOrderType(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Valid data: limit",
			args: args{"l"},
			want: "limit",
		},
		{
			name: "Valid data: market",
			args: args{"m"},
			want: "market",
		},
		{
			name: "Invalid data",
			args: args{"Invalid data"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseOrderType(tt.args.data); got != tt.want {
				t.Errorf("parseOrderType() = %v, want %v", got, tt.want)
			}
		})
	}
}
