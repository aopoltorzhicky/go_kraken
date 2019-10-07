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
func NewAuth(key, secret string, sandbox bool) *AuthClient {
	api := rest.New(key, secret)
	data, err := api.GetWebSocketsToken()
	if err != nil {
		panic(err)
	}
	return &AuthClient{
		Client:            New(sandbox),
		restAPI:           api,
		token:             data.Token,
		tokenExpiresTimer: time.NewTimer(time.Duration(data.Expires)),
	}
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
