package rest

import (
	"bytes"
	"io/ioutil"
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
	type args struct {
		requestURL string
		data       url.Values
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test case 1",
			args: args{
				data: url.Values{
					"test1": {"test2"},
				},
				requestURL: "/test/branch",
			},
			want: "UpUCHWmGNrCFHktS4zqscVY1Aq+qJxkpi2fxICL5swi9IyE+jf2FpzvBlObi2FKEXzEJvVZwIF/dOiungh7q1w==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := New("key", "secret")
			if got := api.getSign(tt.args.requestURL, tt.args.data); got != tt.want {
				t.Errorf("Kraken.getSign() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_prepareRequest(t *testing.T) {
	myURL, _ := url.Parse("https://api.kraken.com/0/private/Balance")
	type args struct {
		method    string
		isPrivate bool
		data      url.Values
	}
	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr bool
	}{
		{
			name: "Error in NewRequest",
			args: args{
				method:    "!@#$%^&*()",
				isPrivate: false,
				data:      nil,
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "Private request with nil data",
			args: args{
				method:    "Balance",
				isPrivate: true,
				data:      nil,
			},
			want:    myURL,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := New("key", "secret")
			got, err := api.prepareRequest(tt.args.method, tt.args.isPrivate, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.prepareRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(got.URL, tt.want) {
				t.Errorf("Kraken.prepareRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, ErrSomething
}

func TestKraken_parseResponse(t *testing.T) {
	type args struct {
		response *http.Response
		retType  interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Invalid status code",
			args: args{
				response: &http.Response{
					StatusCode: 502,
				},
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "Response Body is nil error",
			args: args{
				response: &http.Response{
					StatusCode: 200,
					Body:       nil,
				},
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "ReadAll error",
			args: args{
				response: &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(errReader(0)),
				},
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "Decode JSON error",
			args: args{
				response: &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(``))),
				},
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "Kraken error",
			args: args{
				response: &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"error":["EGeneral:Invalid arguments"]}`))),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Kraken{
				client: &httpMock{
					Error:    nil,
					Response: &http.Response{},
				},
			}
			got, err := api.parseResponse(tt.args.response, tt.args.retType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.parseResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.parseResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKraken_request(t *testing.T) {
	type args struct {
		method    string
		isPrivate bool
		data      url.Values
		retType   interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "prepareRequest error",
			args: args{
				method:    "!@#$%^&*()",
				isPrivate: false,
				data:      nil,
				retType:   nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := New("", "")
			got, err := api.request(tt.args.method, tt.args.isPrivate, tt.args.data, tt.args.retType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kraken.request() = %v, want %v", got, tt.want)
			}
		})
	}
}
