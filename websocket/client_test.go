package websocket

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"
)

var SomethingError = fmt.Errorf("Something went wrong")

func TestClient_createFactory(t *testing.T) {
	type args struct {
		name    string
		factory ParseFactory
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Create book factory",
			args: args{
				ChanBook,
				newBookFactory(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				factories: make(map[string]ParseFactory),
			}
			c.createFactory(tt.args.name, tt.args.factory)
			if len(c.factories) != 1 {
				t.Error("Can't create factory")
			}
		})
	}
}

func TestClient_createFactories(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test create all factories",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				factories: make(map[string]ParseFactory),
			}
			c.createFactories()
			if len(c.factories) != 5 {
				t.Error("Can't create factories")
			}
		})
	}
}

func TestClient_handleEvent(t *testing.T) {
	type args struct {
		msg []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test handle pong event",
			args: args{
				msg: []byte("{ \"event\": \"pong\", \"reqid\": 42 }"),
			},
			wantErr: false,
		},
		{
			name: "Test invalid message",
			args: args{
				msg: []byte("Invalid message"),
			},
			wantErr: true,
		},
		{
			name: "Test invalid pong message",
			args: args{
				msg: []byte("{ \"event\": \"pong\", \"reqid\": \"42\" }"),
			},
			wantErr: true,
		},
		{
			name: "Test handle system status event",
			args: args{
				msg: []byte("{ \"connectionID\":8628615390848610222, \"event\": \"systemStatus\", \"status\": \"online\", \"version\": \"1.0.0\" }"),
			},
			wantErr: false,
		},
		{
			name: "Test handle invalid system status event",
			args: args{
				msg: []byte("{ \"connectionID\": \"8628615390848610222\", \"event\": \"systemStatus\", \"status\": \"online\", \"version\": \"1.0.0\" }"),
			},
			wantErr: true,
		},
		{
			name: "Test handle heartbeat event",
			args: args{
				msg: []byte("{ \"event\": \"heartbeat\" }"),
			},
			wantErr: false,
		},
		{
			name: "Test handle unknown event",
			args: args{
				msg: []byte("{ \"event\": \"unknown\" }"),
			},
			wantErr: false,
		},
		{
			name: "Test handle subscription status event: subscribed",
			args: args{
				msg: []byte("{ \"channelID\": 10001, \"event\": \"subscriptionStatus\", \"status\": \"subscribed\", \"pair\": \"XBT/EUR\", \"subscription\": { \"name\": \"ticker\" }}"),
			},
			wantErr: false,
		},
		{
			name: "Test handle subscription status event: unsubscribed",
			args: args{
				msg: []byte("{ \"channelID\": 10001, \"event\": \"subscriptionStatus\", \"status\": \"unsubscribed\", \"pair\": \"XBT/EUR\", \"subscription\": { \"name\": \"ticker\" }}"),
			},
			wantErr: false,
		},
		{
			name: "Test handle subscription status event: error",
			args: args{
				msg: []byte("{ \"errorMessage\": \"Subscription depth not supported\", \"event\": \"subscriptionStatus\", \"status\": \"error\", \"pair\": \"XBT/EUR\", \"subscription\": { \"name\": \"ticker\" }}"),
			},
			wantErr: false,
		},
		{
			name: "Test handle invalid subscription status event",
			args: args{
				msg: []byte("{ \"channelID\": \"10001\", \"event\": \"subscriptionStatus\", \"status\": \"subscribed\", \"pair\": \"XBT/EUR\", \"subscription\": { \"name\": \"ticker\" }}"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				subscriptions: make(map[int64]*SubscriptionStatus),
			}
			if err := c.handleEvent(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Client.handleEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_lookupByChannelID(t *testing.T) {
	type fields struct {
		subscriptions map[int64]*SubscriptionStatus
	}
	type args struct {
		chanID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SubscriptionStatus
		wantErr bool
	}{
		{
			name: "Test `found subscription` case",
			fields: fields{
				subscriptions: map[int64]*SubscriptionStatus{
					1: &SubscriptionStatus{
						Pair: BTCCAD,
					},
				},
			},
			args: args{
				chanID: int64(1),
			},
			want: &SubscriptionStatus{
				Pair: BTCCAD,
			},
			wantErr: false,
		},
		{
			name: "Test `not found subscription` case",
			fields: fields{
				subscriptions: make(map[int64]*SubscriptionStatus),
			},
			args: args{
				chanID: 1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				subscriptions: tt.fields.subscriptions,
			}
			got, err := c.lookupByChannelID(tt.args.chanID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.lookupByChannelID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.lookupByChannelID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_handleChannel(t *testing.T) {
	type fields struct {
		subscriptions map[int64]*SubscriptionStatus
		factories     map[string]ParseFactory
		listener      chan interface{}
	}
	type args struct {
		msg []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Test invalid json",
			fields: fields{},
			args: args{
				[]byte(""),
			},
			wantErr: true,
		},
		{
			name:   "Test invalid array length",
			fields: fields{},
			args: args{
				[]byte("[ 1 ]"),
			},
			wantErr: true,
		},
		{
			name:   "Test invalid channel ID value",
			fields: fields{},
			args: args{
				[]byte("[ \"1\", [] ]"),
			},
			wantErr: true,
		},
		{
			name: "Test `not found subscription`",
			fields: fields{
				subscriptions: make(map[int64]*SubscriptionStatus),
			},
			args: args{
				[]byte("[ 1, [] ]"),
			},
			wantErr: true,
		},
		{
			name: "Test `not found factory`",
			fields: fields{
				subscriptions: map[int64]*SubscriptionStatus{
					1: &SubscriptionStatus{
						Subscription: Subscription{
							Name: ChanTicker,
						},
					},
				},
				factories: make(map[string]ParseFactory),
			},
			args: args{
				[]byte("[ 1, [] ]"),
			},
			wantErr: true,
		},
		{
			name: "Test invalid parse message",
			fields: fields{
				subscriptions: map[int64]*SubscriptionStatus{
					1: &SubscriptionStatus{
						Subscription: Subscription{
							Name: ChanTicker,
						},
					},
				},
				factories: map[string]ParseFactory{
					ChanTicker: newTickerFactory(),
				},
			},
			args: args{
				[]byte("[ 1, [] ]"),
			},
			wantErr: true,
		},
		{
			name: "Test valid message",
			fields: fields{
				subscriptions: map[int64]*SubscriptionStatus{
					1: &SubscriptionStatus{
						Subscription: Subscription{
							Name: ChanSpread,
						},
					},
				},
				factories: map[string]ParseFactory{
					ChanSpread: newSpreadFactory(),
				},
				listener: make(chan interface{}),
			},
			args: args{
				[]byte("[1 ,[ \"5698.40000\",  \"5700.00000\", \"1542057299.545897\" ]]"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				subscriptions: tt.fields.subscriptions,
				factories:     tt.fields.factories,
				listener:      make(chan interface{}),
			}
			defer close(c.listener)

			go func() {
				<-c.listener
			}()
			if err := c.handleChannel(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Client.handleChannel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_handleMessage(t *testing.T) {
	type fields struct {
		heartbeat     time.Time
		subscriptions map[int64]*SubscriptionStatus
		factories     map[string]ParseFactory
		parameters    *Parameters
	}
	type args struct {
		msg []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test valid data: event",
			fields: fields{
				heartbeat:     time.Now(),
				subscriptions: make(map[int64]*SubscriptionStatus),
				factories:     make(map[string]ParseFactory),
				parameters: &Parameters{
					HeartbeatTimeout: time.Second * 5,
				},
			},
			args: args{
				[]byte("{ \"event\": \"heartbeat\" }"),
			},
			wantErr: false,
		},
		{
			name: "Test valid data: channel",
			fields: fields{
				heartbeat: time.Now(),
				subscriptions: map[int64]*SubscriptionStatus{
					1: &SubscriptionStatus{
						ChannelID: 1,
						Subscription: Subscription{
							Name: ChanSpread,
						},
						Event: EventSubscriptionStatus,
						Pair:  BTCCAD,
					},
				},
				factories: map[string]ParseFactory{
					ChanSpread: newSpreadFactory(),
				},
				parameters: &Parameters{
					HeartbeatTimeout: time.Second * 5,
				},
			},
			args: args{
				[]byte("[1 ,[ \"5698.40000\",  \"5700.00000\", \"1542057299.545897\" ]]"),
			},
			wantErr: false,
		},
		{
			name: "Test valid data: unexpected",
			fields: fields{
				heartbeat:     time.Now(),
				subscriptions: make(map[int64]*SubscriptionStatus),
				factories:     make(map[string]ParseFactory),
				parameters: &Parameters{
					HeartbeatTimeout: time.Second * 5,
				},
			},
			args: args{
				[]byte("Unexpected"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				heartbeat:     tt.fields.heartbeat,
				listener:      make(chan interface{}),
				subscriptions: tt.fields.subscriptions,
				factories:     tt.fields.factories,
				parameters:    tt.fields.parameters,
			}
			defer close(c.listener)
			go func() {
				<-c.listener
			}()
			if err := c.handleMessage(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Client.handleMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type mockAsynchronous struct {
	isConnectionError bool
	isSendError       bool

	Finished chan error
	Data     chan []byte
}

func (m *mockAsynchronous) Connect() error {
	if m.isConnectionError {
		return SomethingError
	}
	return nil
}
func (m *mockAsynchronous) Send(ctx context.Context, msg interface{}) error {
	if m.isSendError {
		return SomethingError
	}
	return nil
}
func (m *mockAsynchronous) Close() {}

func (m *mockAsynchronous) Listen() <-chan []byte {
	return m.Data
}
func (m *mockAsynchronous) Done() <-chan error {
	return m.Finished
}

func (m *mockAsynchronous) cleanup(err error) {
	close(m.Data)
	if err != nil {
		m.Finished <- err
	}
	close(m.Finished)
}

type mockAsyncFactory struct {
	isConnectionError bool
	isSendError       bool
}

func (m *mockAsyncFactory) Create() asynchronous {
	return &mockAsynchronous{
		Finished: make(chan error),
		Data:     make(chan []byte),

		isConnectionError: m.isConnectionError,
		isSendError:       m.isSendError,
	}
}

func TestClient_connect(t *testing.T) {
	type fields struct {
		isError bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "Test correct",
			wantErr: false,
			fields: fields{
				isError: false,
			},
		},
		{
			name:    "Test does not connect",
			wantErr: false,
			fields: fields{
				isError: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				asynchronous: &mockAsynchronous{
					isConnectionError: tt.fields.isError,
				},
				isConnected: false,
			}
			if err := c.connect(); (err != nil) != tt.wantErr {
				if tt.fields.isError == c.IsConnected() {
					t.Errorf("Client.connect() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestClient_dumpParams(t *testing.T) {
	type fields struct {
		parameters *Parameters
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test dump params",
			fields: fields{
				parameters: NewDefaultParameters(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				parameters: tt.fields.parameters,
			}
			c.dumpParams()
		})
	}
}

func TestClient_Ping(t *testing.T) {
	type fields struct {
		isError bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test ping method: success",
			fields: fields{
				isError: false,
			},
			wantErr: false,
		},
		{
			name: "Test ping method: failed",
			fields: fields{
				isError: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				asynchronous: &mockAsynchronous{
					isSendError: tt.fields.isError,
				},
				parameters: NewDefaultSandboxParameters(),
			}
			if err := c.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("Client.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Unsubscribe(t *testing.T) {
	type fields struct {
		isError bool
	}
	type args struct {
		channelType string
		pairs       []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test unsubscribe method: success",
			fields: fields{
				isError: false,
			},
			args: args{
				channelType: ChanTicker,
				pairs:       []string{BTCCAD},
			},
			wantErr: false,
		},
		{
			name: "Test unsubscribe method: failed",
			fields: fields{
				isError: true,
			},
			args: args{
				channelType: ChanTicker,
				pairs:       []string{BTCCAD},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				asynchronous: &mockAsynchronous{
					isSendError: tt.fields.isError,
				},
				parameters: NewDefaultSandboxParameters(),
			}
			if err := c.Unsubscribe(tt.args.channelType, tt.args.pairs); (err != nil) != tt.wantErr {
				t.Errorf("Client.Unsubscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_SubscribeBook(t *testing.T) {
	type fields struct {
		isError bool
	}
	type args struct {
		pairs []string
		depth int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test SubscribeBook method: success",
			fields: fields{
				isError: false,
			},
			args: args{
				depth: Depth10,
				pairs: []string{BTCCAD},
			},
			wantErr: false,
		},
		{
			name: "Test SubscribeBook method: failed",
			fields: fields{
				isError: true,
			},
			args: args{
				depth: Depth10,
				pairs: []string{BTCCAD},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				asynchronous: &mockAsynchronous{
					isSendError: tt.fields.isError,
				},
				parameters: NewDefaultSandboxParameters(),
			}
			if err := c.SubscribeBook(tt.args.pairs, tt.args.depth); (err != nil) != tt.wantErr {
				t.Errorf("Client.SubscribeBook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_SubscribeSpread(t *testing.T) {
	type fields struct {
		isError bool
	}
	type args struct {
		pairs []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test SubscribeSpread method: success",
			fields: fields{
				isError: false,
			},
			args: args{
				pairs: []string{BTCCAD},
			},
			wantErr: false,
		},
		{
			name: "Test SubscribeSpread method: failed",
			fields: fields{
				isError: true,
			},
			args: args{
				pairs: []string{BTCCAD},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				asynchronous: &mockAsynchronous{
					isSendError: tt.fields.isError,
				},
				parameters: NewDefaultSandboxParameters(),
			}
			if err := c.SubscribeSpread(tt.args.pairs); (err != nil) != tt.wantErr {
				t.Errorf("Client.SubscribeSpread() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_SubscribeTrades(t *testing.T) {
	type fields struct {
		isError bool
	}
	type args struct {
		pairs []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test SubscribeTrades method: success",
			fields: fields{
				isError: false,
			},
			args: args{
				pairs: []string{BTCCAD},
			},
			wantErr: false,
		},
		{
			name: "Test SubscribeTrades method: failed",
			fields: fields{
				isError: true,
			},
			args: args{
				pairs: []string{BTCCAD},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				asynchronous: &mockAsynchronous{
					isSendError: tt.fields.isError,
				},
				parameters: NewDefaultSandboxParameters(),
			}
			if err := c.SubscribeTrades(tt.args.pairs); (err != nil) != tt.wantErr {
				t.Errorf("Client.SubscribeTrades() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_SubscribeCandles(t *testing.T) {
	type fields struct {
		isError bool
	}
	type args struct {
		pairs    []string
		interval int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test SubscribeCandles method: success",
			fields: fields{
				isError: false,
			},
			args: args{
				pairs:    []string{BTCCAD},
				interval: Interal10080,
			},
			wantErr: false,
		},
		{
			name: "Test SubscribeCandles method: failed",
			fields: fields{
				isError: true,
			},
			args: args{
				pairs:    []string{BTCCAD},
				interval: Interal10080,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				asynchronous: &mockAsynchronous{
					isSendError: tt.fields.isError,
				},
				parameters: NewDefaultSandboxParameters(),
			}
			if err := c.SubscribeCandles(tt.args.pairs, tt.args.interval); (err != nil) != tt.wantErr {
				t.Errorf("Client.SubscribeCandles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_SubscribeTicker(t *testing.T) {
	type fields struct {
		isError bool
	}
	type args struct {
		pairs []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test SubscribeTicker method: success",
			fields: fields{
				isError: false,
			},
			args: args{
				pairs: []string{BTCCAD},
			},
			wantErr: false,
		},
		{
			name: "Test SubscribeTicker method: failed",
			fields: fields{
				isError: true,
			},
			args: args{
				pairs: []string{BTCCAD},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				asynchronous: &mockAsynchronous{
					isSendError: tt.fields.isError,
				},
				parameters: NewDefaultSandboxParameters(),
			}
			if err := c.SubscribeTicker(tt.args.pairs); (err != nil) != tt.wantErr {
				t.Errorf("Client.SubscribeTicker() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_close(t *testing.T) {
	type fields struct {
		listenerIsNil bool
	}
	type args struct {
		e error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Close with error",
			args: args{
				e: SomethingError,
			},
			fields: fields{
				listenerIsNil: false,
			},
		},
		{
			name: "Close without error",
			args: args{
				e: nil,
			},
			fields: fields{
				listenerIsNil: false,
			},
		},
		{
			name: "Close with nil listener",
			args: args{
				e: nil,
			},
			fields: fields{
				listenerIsNil: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				shutdown: make(chan bool),
			}
			c.listener = make(chan interface{})

			go func() {
				<-c.listener
			}()
			c.close(tt.args.e)
		})
	}
}

func TestClient_exit(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test exit method with error",
			args: args{
				err: SomethingError,
			},
			wantErr: true,
		},
		{
			name: "Test exit method without error",
			args: args{
				err: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				terminal: false,
				shutdown: make(chan bool),
				listener: make(chan interface{}),
			}
			go func() {
				<-c.listener
			}()
			if err := c.exit(tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("Client.exit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Listen(t *testing.T) {
	l := make(chan interface{})
	defer close(l)
	tests := []struct {
		name string
		want <-chan interface{}
	}{
		{
			name: "Test Listen method",
			want: l,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				listener: l,
			}
			if got := c.Listen(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.Listen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_resubscribe(t *testing.T) {
	type fields struct {
		isError bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test without error",
			fields: fields{
				isError: false,
			},
			wantErr: false,
		},
		{
			name: "Test with error",
			fields: fields{
				isError: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				asynchronous: &mockAsynchronous{
					isSendError: tt.fields.isError,
				},
				subscriptions: map[int64]*SubscriptionStatus{
					1: &SubscriptionStatus{
						Pair:         BTCCAD,
						Subscription: Subscription{},
					},
				},
			}
			if err := c.resubscribe(); (err != nil) != tt.wantErr {
				t.Errorf("Client.resubscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_closeAsyncAndWait(t *testing.T) {
	type fields struct {
		init   bool
		isDone bool
	}
	type args struct {
		t time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test init is false",
			fields: fields{
				init:   false,
				isDone: false,
			},
			args: args{
				t: time.Second,
			},
		},
		{
			name: "Test init: true, asynchronous Done",
			fields: fields{
				init:   true,
				isDone: true,
			},
			args: args{
				t: time.Second * 10,
			},
		},
		{
			name: "Test init: true, timeout Done",
			fields: fields{
				init:   true,
				isDone: false,
			},
			args: args{
				t: time.Millisecond,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				asynchronous: &mockAsynchronous{
					Finished: make(chan error),
					Data:     make(chan []byte),
				},
				init: tt.fields.init,
			}
			if tt.fields.isDone {
				go func() {
					time.Sleep(time.Second)
					c.asynchronous.(*mockAsynchronous).Finished <- SomethingError
				}()
			}

			c.closeAsyncAndWait(tt.args.t)
		})
	}
}
func TestClient_controlHeartbeat(t *testing.T) {
	type fields struct {
		isShutdown bool
		heartbeat  time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test control heartbeat: by heartbeat",
			fields: fields{
				isShutdown: false,
				heartbeat:  time.Unix(0, 0),
			},
		},
		{
			name: "Test control heartbeat: by shutdown",
			fields: fields{
				isShutdown: true,
				heartbeat:  time.Now().Add(time.Hour),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				shutdown:  make(chan bool),
				heartbeat: time.Unix(0, 0),
				hbChannel: make(chan error),
			}
			defer close(c.hbChannel)
			defer close(c.shutdown)

			go func() {
				if tt.fields.isShutdown {
					c.shutdown <- true
				} else {
					<-c.hbChannel
				}
			}()
			time.Sleep(time.Microsecond * 500)

			c.controlHeartbeat()
		})
	}
}

func TestClient_listenUpstream(t *testing.T) {
	type fields struct {
		data []byte
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "test with data",
			fields: fields{
				data: []byte("Test data"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				asynchronous: &mockAsynchronous{
					Data: make(chan []byte),
				},
				shutdown:   make(chan bool),
				hbChannel:  make(chan error),
				heartbeat:  time.Now().Add(time.Hour),
				parameters: NewDefaultSandboxParameters(),
			}
			defer close(c.shutdown)
			defer close(c.asynchronous.(*mockAsynchronous).Data)
			defer close(c.hbChannel)

			go func() {
				c.asynchronous.(*mockAsynchronous).Data <- tt.fields.data
				c.shutdown <- true
			}()

			time.Sleep(time.Microsecond)

			c.listenUpstream()
		})
	}
}

func TestClient_reconnect(t *testing.T) {
	type fields struct {
		terminal          bool
		autoReconnect     bool
		isConnectError    bool
		isSendError       bool
		reconnectInterval time.Duration
		subscriptions     map[int64]*SubscriptionStatus
	}
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test terminal is true",
			fields: fields{
				terminal:          true,
				reconnectInterval: time.Millisecond,
				subscriptions:     make(map[int64]*SubscriptionStatus),
			},
			args: args{
				err: nil,
			},
			wantErr: false,
		},
		{
			name: "Test terminal is true with error",
			fields: fields{
				terminal:          true,
				reconnectInterval: time.Millisecond,
				subscriptions:     make(map[int64]*SubscriptionStatus),
			},
			args: args{
				err: SomethingError,
			},
			wantErr: true,
		},
		{
			name: "Test autoReconnect false",
			fields: fields{
				terminal:          false,
				autoReconnect:     false,
				reconnectInterval: time.Millisecond,
				subscriptions:     make(map[int64]*SubscriptionStatus),
			},
			args: args{
				err: nil,
			},
			wantErr: true,
		},
		{
			name: "Test autoReconnect false with error",
			fields: fields{
				terminal:          false,
				autoReconnect:     false,
				reconnectInterval: time.Millisecond,
				subscriptions:     make(map[int64]*SubscriptionStatus),
			},
			args: args{
				err: SomethingError,
			},
			wantErr: true,
		},
		{
			name: "Test autoReconnect true",
			fields: fields{
				terminal:          false,
				autoReconnect:     true,
				isConnectError:    false,
				reconnectInterval: time.Millisecond,
				subscriptions:     make(map[int64]*SubscriptionStatus),
			},
			args: args{
				err: nil,
			},
			wantErr: false,
		},
		{
			name: "Test isConnectError true",
			fields: fields{
				terminal:          false,
				autoReconnect:     true,
				isConnectError:    true,
				reconnectInterval: time.Millisecond,
				subscriptions:     make(map[int64]*SubscriptionStatus),
			},
			args: args{
				err: nil,
			},
			wantErr: true,
		},
		{
			name: "Test isSendError true",
			fields: fields{
				terminal:       false,
				autoReconnect:  true,
				isConnectError: false,
				isSendError:    true,
				subscriptions: map[int64]*SubscriptionStatus{
					1: &SubscriptionStatus{
						Pair: BTCCAD,
						Subscription: Subscription{
							Name: ChanTicker,
						},
					},
				},
			},
			args: args{
				err: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				asyncFactory: &mockAsyncFactory{
					isConnectionError: tt.fields.isConnectError,
					isSendError:       tt.fields.isSendError,
				},
				asynchronous:  nil,
				isConnected:   false,
				terminal:      tt.fields.terminal,
				init:          false,
				heartbeat:     time.Now().Add(time.Hour),
				hbChannel:     make(chan error),
				parameters:    NewDefaultSandboxParameters(),
				shutdown:      make(chan bool),
				listener:      make(chan interface{}),
				subscriptions: tt.fields.subscriptions,
				factories:     make(map[string]ParseFactory),
			}
			c.parameters.AutoReconnect = tt.fields.autoReconnect
			c.parameters.ReconnectInterval = tt.fields.reconnectInterval
			c.parameters.ReconnectAttempts = 2

			go func() {
				<-c.Listen()
			}()

			if err := c.reconnect(tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("Client.reconnect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Close(t *testing.T) {
	type fields struct {
		isShutdown bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test shutdown",
			fields: fields{
				isShutdown: true,
			},
		},
		{
			name: "Test is not shutdown",
			fields: fields{
				isShutdown: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				asynchronous: &mockAsynchronous{
					Finished: make(chan error),
					Data:     make(chan []byte),
				},
				terminal:   false,
				init:       false,
				parameters: NewDefaultSandboxParameters(),
				shutdown:   make(chan bool),
			}
			c.parameters.ShutdownTimeout = time.Second

			go func() {
				if tt.fields.isShutdown {
					c.shutdown <- true
				}
			}()

			c.Close()
		})
	}
}

func TestClient_Connect(t *testing.T) {
	tests := []struct {
		name           string
		isConnectError bool
		wantErr        bool
	}{
		{
			name:           "Test without error",
			isConnectError: false,
		},
		{
			name:           "Test with error",
			isConnectError: true,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				asyncFactory: &mockAsyncFactory{
					isConnectionError: tt.isConnectError,
				},
				heartbeat:  time.Now().Add(time.Hour),
				hbChannel:  make(chan error),
				parameters: NewDefaultSandboxParameters(),
				listener:   make(chan interface{}),
			}
			if err := c.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("Client.Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_listenDisconnect(t *testing.T) {
	type fields struct {
		isAsynchronousDone bool
		isListenHeartbeat  bool

		err error
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "test asynchronous branch",
			fields: fields{
				isAsynchronousDone: true,
				isListenHeartbeat:  false,
			},
		},
		{
			name: "test asynchronous branch with error",
			fields: fields{
				isAsynchronousDone: true,
				isListenHeartbeat:  false,
				err:                SomethingError,
			},
		},
		{
			name: "test heartbeat branch",
			fields: fields{
				isAsynchronousDone: false,
				isListenHeartbeat:  true,
			},
		},
		{
			name: "test heartbeat branch with error",
			fields: fields{
				isAsynchronousDone: false,
				isListenHeartbeat:  true,
				err:                SomethingError,
			},
		},
		{
			name: "test heartbeat branch with connect error",
			fields: fields{
				isAsynchronousDone: false,
				isListenHeartbeat:  true,
				err:                SomethingError,
			},
		},
		{
			name: "test both branches",
			fields: fields{
				isAsynchronousDone: true,
				isListenHeartbeat:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				asynchronous: &mockAsynchronous{
					Finished:          make(chan error),
					Data:              make(chan []byte),
					isConnectionError: true,
				},
				isConnected: false,
				terminal:    false,
				heartbeat:   time.Now().Add(time.Hour),
				hbChannel:   make(chan error),
				parameters:  NewDefaultSandboxParameters(),
				shutdown:    make(chan bool),
			}
			c.parameters.AutoReconnect = false
			go func() {
				if tt.fields.isListenHeartbeat {
					c.heartbeat = time.Unix(0, 0)
					c.controlHeartbeat()
				}
				if tt.fields.isAsynchronousDone {
					c.asynchronous.(*mockAsynchronous).cleanup(SomethingError)
				}
			}()

			c.listenDisconnect()
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		sandbox bool
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		{
			name: "Test sandbox",
			args: args{
				sandbox: true,
			},
		},
		{
			name: "Test prod",
			args: args{
				sandbox: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.sandbox)
			if tt.args.sandbox {
				if got.parameters.URL != sandboxBaseURL {
					t.Errorf("URL = %s, want %s", got.parameters.URL, sandboxBaseURL)
				}
			} else {
				if got.parameters.URL != prodBaseURL {
					t.Errorf("URL = %s, want %s", got.parameters.URL, prodBaseURL)
				}
			}
			if len(got.factories) != 5 {
				t.Errorf("Factories count = %d, want %d", len(got.factories), 5)
			}
		})
	}
}

func Test_websocketAsynchronousFactory_Create(t *testing.T) {
	params := NewDefaultSandboxParameters()
	tests := []struct {
		name string
		want asynchronous
	}{
		{
			name: "test creation",
			want: &ws{
				BaseURL:      params.URL,
				downstream:   make(chan []byte),
				shutdown:     make(chan struct{}),
				finished:     make(chan error),
				logTransport: params.LogTransport,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &websocketAsynchronousFactory{
				parameters: params,
			}
			got := w.Create()
			if got.(*ws).BaseURL != params.URL {
				t.Errorf("websocketAsynchronousFactory.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
