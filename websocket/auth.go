package websocket

import (
	"context"
	"time"

	rest "github.com/aopoltorzhicky/go_kraken/rest"
)

// AuthClient - client for private websocket channels
type AuthClient struct {
	*Client

	restAPI           *rest.Kraken
	token             string
	tokenExpiresTimer *time.Timer
}

// NewAuth - constructor for AuthClient
func NewAuth(key, secret string) *AuthClient {
	api := rest.New(key, secret)
	data, err := api.GetWebSocketsToken()
	if err != nil {
		panic(err)
	}

	params := NewDefaultAuthParameters()
	c := &AuthClient{
		Client: &Client{
			asyncFactory: &websocketAsynchronousFactory{
				parameters: params,
			},
			isConnected:   false,
			parameters:    params,
			listener:      make(chan interface{}),
			terminal:      false,
			shutdown:      make(chan struct{}),
			asynchronous:  nil,
			heartbeat:     time.Now().Add(params.HeartbeatTimeout),
			hbChannel:     make(chan error),
			subscriptions: make(map[int64]*SubscriptionStatus),
			factories:     make(map[string]ParseFactory),
		},
		restAPI:           api,
		token:             data.Token,
		tokenExpiresTimer: time.NewTimer(time.Duration(data.Expires)),
	}
	c.createFactories()
	c.createAuthFactories()
	return c
}

func (c *AuthClient) createAuthFactories() {
	c.createFactory(ChanOwnTrades, newOwnTradesFactory())
	c.createFactory(ChanOpenOrders, newOpenOrdersFactory())
}

func (c *AuthClient) subscribeAuthChannel(channelName string) error {
	ctx, cxl := context.WithTimeout(context.Background(), c.parameters.ContextTimeout)
	defer cxl()
	s := AuthSubscriptionRequest{
		Event: EventSubscribe,
		Subs: AuthDataRequest{
			Name:  channelName,
			Token: c.token,
		},
	}
	return c.asynchronous.Send(ctx, s)
}

// SubscribeOwnTrades - method tries to subscribe on OwnTrades channel events
func (c *AuthClient) SubscribeOwnTrades() error {
	return c.subscribeAuthChannel(ChanOwnTrades)
}

// SubscribeOpenOrders - method tries to subscribe on OpenOrders channel events
func (c *AuthClient) SubscribeOpenOrders() error {
	return c.subscribeAuthChannel(ChanOpenOrders)
}

// AddOrder - method adds new order.
func (c *AuthClient) AddOrder(req AddOrderRequest) error {
	req.Event = EventAddOrder
	req.Token = c.token

	ctx, cxl := context.WithTimeout(context.Background(), c.parameters.ContextTimeout)
	defer cxl()
	return c.asynchronous.Send(ctx, req)
}

// CancelOrder - method cancels order or list of orders.
func (c *AuthClient) CancelOrder(orderIDs []string) error {
	req := CancelOrderRequest{
		AuthRequest: AuthRequest{
			Token: c.token,
			Event: EventCancelOrder,
		},
		TxID: orderIDs,
	}
	ctx, cxl := context.WithTimeout(context.Background(), c.parameters.ContextTimeout)
	defer cxl()
	return c.asynchronous.Send(ctx, req)
}
