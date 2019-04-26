package rest

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		key    string
		secret string
	}
	tests := []struct {
		name string
		args args
		want *Kraken
	}{
		{
			name: "Without keys",
			args: args{
				key:    "",
				secret: "",
			},
			want: &Kraken{
				key:    "",
				secret: "",
				client: http.DefaultClient,
			},
		},
		{
			name: "With keys",
			args: args{
				key:    "key",
				secret: "secret",
			},
			want: &Kraken{
				key:    "key",
				secret: "secret",
				client: http.DefaultClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.key, tt.args.secret); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
