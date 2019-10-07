package websocket

import (
	"context"
	"log"
	"net"
	"reflect"
	"sync"
	"testing"

	WS "github.com/gorilla/websocket"
)

type (
	mockWS struct {
		isWriteError bool
		isReadError  bool
		isCloseError bool
		err          error
	}
)

func (ws *mockWS) WriteMessage(messageType int, msg []byte) error {
	if ws.isWriteError {
		return ErrSomething
	}
	return nil
}

func (ws *mockWS) ReadMessage() (int, []byte, error) {
	if ws.isReadError {
		if ws.err != nil {
			return 0, nil, ws.err
		}
		return 0, nil, ErrSomething
	}
	return 1, []byte("Test message"), nil
}

func (ws *mockWS) Close() error {
	if ws.isCloseError {
		return ErrSomething
	}
	return nil
}

func Test_ws_Done(t *testing.T) {
	c := make(chan error)
	tests := []struct {
		name string
		want <-chan error
	}{
		{
			name: "Default test",
			want: c,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &ws{
				finished: c,
			}
			if got := w.Done(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ws.Done() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ws_Listen(t *testing.T) {
	downstream := make(chan []byte)
	tests := []struct {
		name string
		want <-chan []byte
	}{
		{
			name: "Default test",
			want: downstream,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &ws{
				downstream: downstream,
			}
			if got := w.Listen(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ws.Listen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ws_cleanup(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "without error",
			args: args{err: nil},
		},
		{
			name: "with error",
			args: args{err: ErrSomething},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &ws{
				downstream:   make(chan []byte),
				userShutdown: false,
				shutdown:     make(chan struct{}),
				finished:     make(chan error),
			}

			if tt.args.err != nil {
				go func() {
					<-w.finished
				}()
			}
			w.cleanup(tt.args.err)
		})
	}
}

func Test_ws_Close(t *testing.T) {
	type fields struct {
		ws mockWS
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "WS is nil",
		},
		{
			name: "WS close return error",
			fields: fields{
				ws: mockWS{
					isCloseError: true,
				},
			},
		},
		{
			name: "WS close return nil",
			fields: fields{
				ws: mockWS{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &ws{
				ws:           &tt.fields.ws,
				wsLock:       sync.Mutex{},
				userShutdown: false,
			}
			w.Close()
			if !w.userShutdown {
				t.Errorf("ws.Listen() = %v, want %v", w.userShutdown, true)
			}
		})
	}
}

func Test_ws_listenWs(t *testing.T) {
	type fields struct {
		ws           connInterface
		logTransport bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "WS is nil",
		},
		{
			name: "Shutdown by CloseError",
			fields: fields{
				ws: &mockWS{
					isReadError: true,
					err:         &WS.CloseError{},
				},
			},
		},
		{
			name: "Shutdown by OpError",
			fields: fields{
				ws: &mockWS{
					isReadError: true,
					err:         &net.OpError{},
				},
			},
		},
		{
			name: "Good test",
			fields: fields{
				ws:           &mockWS{},
				logTransport: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &ws{
				ws:           tt.fields.ws,
				downstream:   make(chan []byte),
				userShutdown: false,
				logTransport: tt.fields.logTransport,
				shutdown:     make(chan struct{}),
				finished:     make(chan error),
			}
			if w.logTransport {
				go func() {
					<-w.downstream
					w.cleanup(nil)
				}()
			} else {
				go func() {
					<-w.finished
				}()
			}
			w.listenWs()
		})
	}
}

func Test_ws_Send(t *testing.T) {
	type fields struct {
		ws           connInterface
		logTransport bool
		isCtxDone    bool
		isCtxNil     bool
		isShutdown   bool
	}
	type args struct {
		msg interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "WS is nil",
			wantErr: true,
			fields: fields{
				isCtxNil: true,
			},
		},
		{
			name:    "JSON error",
			wantErr: true,
			args: args{
				msg: make(chan string),
			},
			fields: fields{
				ws: &mockWS{},
			},
		},
		{
			name:    "Ctx is Done",
			wantErr: true,
			args: args{
				msg: "test message",
			},
			fields: fields{
				ws:        &mockWS{},
				isCtxDone: true,
			},
		},
		{
			name:    "Write message error",
			wantErr: true,
			args: args{
				msg: map[string]string{},
			},
			fields: fields{
				ws: &mockWS{
					isWriteError: true,
				},
				logTransport: true,
				isCtxDone:    false,
			},
		},
		{
			name:    "Good test",
			wantErr: false,
			args: args{
				msg: map[string]string{},
			},
			fields: fields{
				ws:        &mockWS{},
				isCtxDone: false,
			},
		},
		{
			name:    "Shutdown",
			wantErr: true,
			args: args{
				msg: map[string]string{},
			},
			fields: fields{
				ws:         &mockWS{},
				isCtxDone:  false,
				isShutdown: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &ws{
				ws:           tt.fields.ws,
				wsLock:       sync.Mutex{},
				downstream:   make(chan []byte),
				userShutdown: false,
				logTransport: tt.fields.logTransport,
				shutdown:     make(chan struct{}),
				finished:     make(chan error),
			}
			var ctx context.Context
			if !tt.fields.isCtxNil {
				var cnt context.CancelFunc
				ctx, cnt = context.WithCancel(context.Background())
				if tt.fields.isCtxDone {
					cnt()
				}
				log.Println(cnt)
			}
			if tt.fields.logTransport {
				go func() {
					<-w.finished
				}()
			}
			if tt.fields.isShutdown {
				close(w.shutdown)
			}
			if err := w.Send(ctx, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("ws.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ws_Connect(t *testing.T) {
	type fields struct {
		ws      connInterface
		baseURL string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "WS is nil",
			wantErr: false,
			fields: fields{
				baseURL: SandboxBaseURL,
			},
		},
		{
			name:    "WS is nil",
			wantErr: true,
			fields: fields{
				baseURL: "",
			},
		},
		{
			name:    "WS is nil",
			wantErr: true,
			fields: fields{
				baseURL: "ws://ws-sandbox.kraken.com",
			},
		},
		{
			name: "WS is not nil",
			fields: fields{
				ws: &mockWS{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &ws{
				ws:            tt.fields.ws,
				wsLock:        sync.Mutex{},
				BaseURL:       tt.fields.baseURL,
				TLSSkipVerify: true,
				downstream:    make(chan []byte),
				userShutdown:  false,
				logTransport:  true,
				shutdown:      make(chan struct{}),
				finished:      make(chan error),
			}
			if err := w.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("ws.Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
