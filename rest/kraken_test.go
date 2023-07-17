package rest

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
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

func TestKraken_getSign(t *testing.T) {
	type fields struct {
		key        string
		isNotValid bool
		client     clientInterface
	}
	type args struct {
		requestURL string
		data       url.Values
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Good test",
			fields: fields{
				key:    "api-key",
				client: nil,
			},
			args: args{
				requestURL: "test-url",
				data: url.Values{
					"a": {"b"},
				},
			},
			want:    "U7jUymCBhD6Q+GBoUZRshexaUNfkVoHZsdkbLTqcT5502b4Qx7HDxLMhlVMN9tIit+Ir6UmzzjStKlSLURy4Xg==",
			wantErr: false,
		}, {
			name: "Invalid base64 decode",
			fields: fields{
				key:        "api-key",
				isNotValid: true,
				client:     nil,
			},
			args: args{
				requestURL: "test-url",
				data: url.Values{
					"a": {"b"},
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				key:    tt.fields.key,
				secret: deadbeaf,
				client: tt.fields.client,
			}

			if tt.fields.isNotValid {
				api.secret = invalid
			}

			got, err := api.getSign(tt.args.requestURL, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.getSign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Kraken.getSign() = %v, want %v", got, tt.want)
			}
		})
	}
}

const (
	deadbeaf = "deadbeaf"
	invalid  = "invalid"
)

func TestKraken_prepareRequest(t *testing.T) {
	type fields struct {
		key        string
		isNotValid bool
	}
	type args struct {
		method    string
		isPrivate bool
		data      url.Values
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Request
		wantErr bool
	}{
		{
			name: "Good test",
			fields: fields{
				key: "api-key",
			},
			args: args{
				method:    "time",
				isPrivate: false,
				data: url.Values{
					"a": {"b"},
				},
			},
			wantErr: false,
			want:    nil,
		}, {
			name: "Data is nil",
			fields: fields{
				key: "api-key",
			},
			args: args{
				method:    "time",
				isPrivate: false,
				data:      nil,
			},
			wantErr: false,
			want:    nil,
		}, {
			name: "Invalid request creation",
			fields: fields{
				key: "api-key",
			},
			args: args{
				method:    "^%&RDRER^WER*XRW&",
				isPrivate: false,
				data:      nil,
			},
			wantErr: true,
			want:    nil,
		}, {
			name: "Private request",
			fields: fields{
				key: "api-key",
			},
			args: args{
				method:    "time",
				isPrivate: true,
				data:      nil,
			},
			wantErr: false,
		}, {
			name: "Private request: invalid signature",
			fields: fields{
				key:        "api-key",
				isNotValid: true,
			},
			args: args{
				method:    "time",
				isPrivate: true,
				data:      nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := deadbeaf
			if tt.fields.isNotValid {
				s = invalid
			}
			api := New(tt.fields.key, s)
			got, err := api.prepareRequest(tt.args.method, tt.args.isPrivate, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.prepareRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got != nil {
					if got.Method != "POST" {
						t.Errorf("Kraken.prepareRequest() expected = POST, got %v", got.Method)
						return
					}
					if got.URL.Path != "/0/public/time" && got.URL.Path != "/0/private/time" {
						t.Errorf("Kraken.prepareRequest() expected = /0/public/time, got %v", got.URL.Path)
						return
					}
					if got.URL.Host != "api.kraken.com" {
						t.Errorf("Kraken.prepareRequest() expected = api.kraken.com, got %v", got.URL.Host)
						return
					}
					if got.URL.Scheme != "https" {
						t.Errorf("Kraken.prepareRequest() expected = https, got %v", got.URL.Scheme)
						return
					}
				} else {
					t.Errorf("Kraken.prepareRequest() got = %v, wantErr %v", got, tt.wantErr)
					return
				}
			}
		})
	}
}

func TestKraken_parseResponse(t *testing.T) {
	type fields struct {
		key string
	}
	type args struct {
		response *http.Response
		retType  interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Status code != 200",
			fields: fields{
				key: "api-key",
			},
			args: args{
				response: &http.Response{
					StatusCode: 400,
				},
				retType: nil,
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "Response body is nil",
			fields: fields{
				key: "api-key",
			},
			args: args{
				response: &http.Response{
					StatusCode: 200,
				},
				retType: nil,
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "Response marshalling error",
			fields: fields{
				key: "api-key",
			},
			args: args{
				response: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				},
				retType: nil,
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "Response with Kraken error",
			fields: fields{
				key: "api-key",
			},
			args: args{
				response: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString("{\"error\": [\"test\"]}")),
				},
				retType: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := New(tt.fields.key, deadbeaf)
			err := api.parseResponse(tt.args.response, tt.args.retType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.parseResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestKraken_request(t *testing.T) {
	type fields struct {
		key string
	}
	type args struct {
		method    string
		isPrivate bool
		data      url.Values
		retType   interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "error request",
			args: args{
				method:    "2782r73rt72teX&Tv^&r@#&r",
				isPrivate: true,
				data:      nil,
				retType:   nil,
			},
			fields: fields{
				key: "key",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := New(tt.fields.key, invalid)
			err := api.request(tt.args.method, tt.args.isPrivate, tt.args.data, tt.args.retType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
