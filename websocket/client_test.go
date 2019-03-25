package websocket

import (
	"reflect"
	"testing"
)

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				subscriptions: tt.fields.subscriptions,
				factories:     tt.fields.factories,
			}
			if err := c.handleChannel(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Client.handleChannel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
