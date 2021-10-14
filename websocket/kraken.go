package websocket

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/aopoltorzhicky/go_kraken/rest"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Kraken -
type Kraken struct {
	url   string
	token string

	conn          *websocket.Conn
	subscriptions map[int64]*SubscriptionStatus

	reconnectTimeout time.Duration
	readTimeout      time.Duration
	heartbeatTimeout time.Duration

	msg     chan Update
	connect chan struct{}
	stop    chan struct{}

	wg sync.WaitGroup
}

// New -
func NewKraken(url string, opts ...KrakenOption) *Kraken {
	kraken := Kraken{
		url:              url,
		reconnectTimeout: 5 * time.Second,
		readTimeout:      15 * time.Second,
		heartbeatTimeout: 10 * time.Second,
		subscriptions:    make(map[int64]*SubscriptionStatus),
		connect:          make(chan struct{}, 1),
		msg:              make(chan Update, 1024),
		stop:             make(chan struct{}, 1),
	}

	for i := range opts {
		opts[i](&kraken)
	}

	return &kraken
}

// Connect to the Kraken API, this should only be called once.
func (k *Kraken) Connect() error {
	k.wg.Add(1)
	go k.managerThread()

	if err := k.dial(); err != nil {
		return err
	}

	k.wg.Add(1)
	go k.listenSocket()

	return nil
}

func (k *Kraken) dial() error {
	dialer := websocket.Dialer{
		Subprotocols:    []string{"p1", "p2"},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Proxy:           http.ProxyFromEnvironment,
	}

	c, resp, err := dialer.Dial(k.url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	k.conn = c
	return nil
}

func (k *Kraken) managerThread() {
	defer k.wg.Done()

	heartbeat := time.NewTicker(k.heartbeatTimeout)
	defer heartbeat.Stop()

	for {
		select {
		case <-k.stop:
			return
		case <-k.connect:
			time.Sleep(k.reconnectTimeout)

			log.Warnf("reconnecting...")

			if err := k.dial(); err != nil {
				log.Error(err)
				k.connect <- struct{}{}
				continue
			}

			k.wg.Add(1)
			go k.listenSocket()

			if err := k.resubscribe(); err != nil {
				log.Error(err)
			}
		case <-heartbeat.C:
			if err := k.send(PingRequest{
				Event: EventPing,
			}); err != nil {
				log.Println(err)
				k.connect <- struct{}{}
			}
		}
	}
}

func (k *Kraken) resubscribe() error {
	for _, sub := range k.subscriptions {
		switch sub.Subscription.Name {
		// Private Channels
		case ChanOwnTrades, ChanOpenOrders:
			return k.subscribeToPrivate(sub.Subscription.Name)
		default:
			if err := k.send(SubscriptionRequest{
				Event:        EventSubscribe,
				Pairs:        []string{sub.Pair},
				Subscription: sub.Subscription,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

// Listen provides an atomic interface for receiving API messages.
// When a websocket connection is terminated, the publisher channel will close.
func (k *Kraken) Listen() <-chan Update {
	return k.msg
}

// Close - provides an interface for a user initiated shutdown.
func (k *Kraken) Close() error {
	for i := 0; i < 2; i++ {
		k.stop <- struct{}{}
	}
	k.wg.Wait()

	if k.conn != nil {
		if err := k.conn.Close(); err != nil {
			return err
		}
	}

	close(k.stop)
	close(k.msg)
	close(k.connect)
	return nil
}

func (k *Kraken) send(msg interface{}) error {
	if k.conn == nil {
		return nil
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log.Tracef("client->server: %s", string(data))
	return k.conn.WriteMessage(websocket.TextMessage, data)
}

func (k *Kraken) listenSocket() {
	defer k.wg.Done()

	if k.conn == nil {
		return
	}

	if err := k.conn.SetReadDeadline(time.Now().Add(k.readTimeout)); err != nil {
		log.Error(err)
		return
	}

	for {
		select {
		case <-k.stop:
			return
		default:
			_, msg, err := k.conn.ReadMessage()
			if err != nil {
				log.Error(err)
				k.connect <- struct{}{}
				return
			}

			if err := k.conn.SetReadDeadline(time.Now().Add(k.readTimeout)); err != nil {
				log.Error(err)
				k.connect <- struct{}{}
				return
			}

			log.Tracef("server->client: %s", string(msg))

			if err := k.handleMessage(msg); err != nil {
				log.Error(err)
			}
		}
	}
}

func (k *Kraken) handleMessage(data []byte) error {
	if len(data) == 0 {
		return errors.Errorf("Empty response: %s", string(data))
	}
	switch data[0] {
	case '[':
		return k.handleChannel(data)
	case '{':
		return k.handleEvent(data)
	default:
		return errors.Errorf("Unexpected message: %s", string(data))
	}
}

// SubscribeTicker - Ticker information includes best ask and best bid prices, 24hr volume, last trade price, volume weighted average price, etc for a given currency pair. A ticker message is published every time a trade or a group of trade happens.
func (k *Kraken) SubscribeTicker(pairs []string) error {
	return k.send(SubscriptionRequest{
		Event: EventSubscribe,
		Pairs: pairs,
		Subscription: Subscription{
			Name: ChanTicker,
		},
	})
}

// SubscribeCandles - Open High Low Close (Candle) feed for a currency pair and interval period.
func (k *Kraken) SubscribeCandles(pairs []string, interval int64) error {
	return k.send(SubscriptionRequest{
		Event: EventSubscribe,
		Pairs: pairs,
		Subscription: Subscription{
			Name:     ChanCandles,
			Interval: interval,
		},
	})
}

// SubscribeTrades - Trade feed for a currency pair.
func (k *Kraken) SubscribeTrades(pairs []string) error {
	return k.send(SubscriptionRequest{
		Event: EventSubscribe,
		Pairs: pairs,
		Subscription: Subscription{
			Name: ChanTrades,
		},
	})
}

// SubscribeSpread - Spread feed to show best bid and ask price for a currency pair
func (k *Kraken) SubscribeSpread(pairs []string) error {
	return k.send(SubscriptionRequest{
		Event: EventSubscribe,
		Pairs: pairs,
		Subscription: Subscription{
			Name: ChanSpread,
		},
	})
}

// SubscribeBook - Order book levels. On subscription, a snapshot will be published at the specified depth, following the snapshot, level updates will be published.
func (k *Kraken) SubscribeBook(pairs []string, depth int64) error {
	return k.send(SubscriptionRequest{
		Event: EventSubscribe,
		Pairs: pairs,
		Subscription: Subscription{
			Name:  ChanBook,
			Depth: depth,
		},
	})
}

// Unsubscribe - Unsubscribe from single subscription, can specify multiple currency pairs.
func (k *Kraken) Unsubscribe(channelType string, pairs []string) error {
	return k.send(UnsubscribeRequest{
		Event: EventUnsubscribe,
		Pairs: pairs,
		Subscription: Subscription{
			Name: channelType,
		},
	})
}

// UnsubscribeCandles - Unsubscribe from candles subscription, can specify multiple currency pairs.
func (k *Kraken) UnsubscribeCandles(pairs []string, interval int64) error {
	return k.send(UnsubscribeRequest{
		Event: EventUnsubscribe,
		Pairs: pairs,
		Subscription: Subscription{
			Name:     ChanCandles,
			Interval: interval,
		},
	})
}

// UnsubscribeBook - Unsubscribe from order book subscription, can specify multiple currency pairs.
func (k *Kraken) UnsubscribeBook(pairs []string, depth int64) error {
	return k.send(UnsubscribeRequest{
		Event: EventUnsubscribe,
		Pairs: pairs,
		Subscription: Subscription{
			Name:  ChanBook,
			Depth: depth,
		},
	})
}

// Authenticate - authenticate in private Websocket API
func (k *Kraken) Authenticate(key, secret string) error {
	data, err := rest.New(key, secret).GetWebSocketsToken()
	if err != nil {
		return err
	}
	k.token = data.Token
	return nil
}

func (k *Kraken) subscribeToPrivate(channelName string) error {
	return k.send(AuthSubscriptionRequest{
		Event: EventSubscribe,
		Subs: AuthDataRequest{
			Name:  channelName,
			Token: k.token,
		},
	})
}

// SubscribeOwnTrades - method tries to subscribe on OwnTrades channel events
func (k *Kraken) SubscribeOwnTrades() error {
	return k.subscribeToPrivate(ChanOwnTrades)
}

// SubscribeOpenOrders - method tries to subscribe on OpenOrders channel events
func (k *Kraken) SubscribeOpenOrders() error {
	return k.subscribeToPrivate(ChanOpenOrders)
}

// AddOrder - method adds new order.
func (k *Kraken) AddOrder(req AddOrderRequest) error {
	req.Event = EventAddOrder
	req.Token = k.token
	return k.send(req)
}

// CancelOrder - method cancels order or list of orders.
func (k *Kraken) CancelOrder(orderIDs []string) error {
	return k.send(CancelOrderRequest{
		AuthRequest: AuthRequest{
			Token: k.token,
			Event: EventCancelOrder,
		},
		TxID: orderIDs,
	})
}

// CancelAll - method cancels order or list of orders.
func (k *Kraken) CancelAll() error {
	return k.send(AuthRequest{
		Token: k.token,
		Event: EventCancelAll,
	})
}

// CancelAllOrdersAfter -  provides a `Dead Man's Switch` mechanism to protect the client from network malfunction, extreme latency or unexpected matching engine downtime. The client can send a request with a timeout (in seconds), that will start a countdown timer which will cancel *all* client orders when the timer expires.
func (k *Kraken) CancelAllOrdersAfter(timeout int64) error {
	return k.send(CancelAllOrdersAfterRequest{
		AuthRequest: AuthRequest{
			Token: k.token,
			Event: EventCancelAllOrdersAfter,
		},
		Timeout: timeout,
	})
}
