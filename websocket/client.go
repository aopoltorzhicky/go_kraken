package websocket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
	"unicode"
)

type asynchronous interface {
	Connect() error
	Send(ctx context.Context, msg interface{}) error
	Listen() <-chan []byte
	Close()
	Done() <-chan error
}

type asynchronousFactory interface {
	Create() asynchronous
}

type websocketAsynchronousFactory struct {
	parameters *Parameters
}

// Client provides a unified interface for users to interact with the Kraken Websocket API.
// nolint:megacheck,structcheck
type Client struct {
	asyncFactory asynchronousFactory // for re-creating transport during reconnects

	timeout      int64 // read timeout
	asynchronous asynchronous
	heartbeat    time.Time
	hbChannel    chan error

	// connection & operational behavior
	parameters *Parameters

	// channel to stop all routines
	shutdown chan struct{}

	// downstream listener channel to deliver API objects
	listener chan interface{}

	subscriptions map[int64]*SubscriptionStatus
	factories     map[string]ParseFactory

	isConnected bool
	terminal    bool
	init        bool

	// race management
	waitGroup sync.WaitGroup

	isConnectedMux sync.Mutex
	heartbeatMux   sync.Mutex
}

// Create returns a new websocket transport.
func (w *websocketAsynchronousFactory) Create() asynchronous {
	return newWs(w.parameters.URL, w.parameters.LogTransport)
}

// New creates a default client.
func New() *Client {
	params := NewDefaultParameters()
	c := &Client{
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
	}
	c.createFactories()
	return c
}

// NewSandbox - creates a default sandbox client.
func NewSandbox() *Client {
	params := NewDefaultSandboxParameters()
	c := &Client{
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
	}
	c.createFactories()
	return c
}

func (c *Client) controlHeartbeat() {
	for {
		select {
		case <-c.shutdown:
			log.Printf("Shutdown control heartbeat")
			return
		default:
		}
		c.heartbeatMux.Lock()
		if time.Now().After(c.heartbeat) {
			c.hbChannel <- fmt.Errorf("Heartbeat disconnect on client")
			c.heartbeatMux.Unlock()
			return
		}
		c.heartbeatMux.Unlock()
		time.Sleep(c.parameters.HeartbeatCheckPeriod)
	}
}

func (c *Client) listenHeartbeat() <-chan error {
	return c.hbChannel
}

func (c *Client) updateHeartbeat() {
	c.heartbeatMux.Lock()
	defer c.heartbeatMux.Unlock()
	c.heartbeat = time.Now().Add(c.parameters.HeartbeatTimeout)
}

func (c *Client) listenDisconnect() {
	select {
	case e := <-c.asynchronous.Done():
		log.Printf("socket disconnect")
		if e != nil {
			log.Printf("socket disconnect: %s", e.Error())
		}
		c.setIsConnected()

		if err := c.reconnect(e); err != nil {
			log.Printf("socket disconnect: %s", err.Error())
		}
	case <-c.shutdown:
		log.Printf("Shutdown listen disconnect")
		c.setIsConnected()
		return

	case e := <-c.listenHeartbeat():
		log.Printf("Heartbeat")
		if e != nil {
			c.closeAsyncAndWait(c.parameters.ShutdownTimeout)

			if err := c.reconnect(nil); err != nil {
				log.Printf("socket disconnect: %s", err.Error())
			}
		}
		c.setIsConnected()
	}
}

func (c *Client) dumpParams() {
	log.Print("----Kraken Client Parameters----")
	log.Printf("AutoReconnect=%t", c.parameters.AutoReconnect)
	log.Printf("ReconnectInterval=%s", c.parameters.ReconnectInterval)
	log.Printf("ReconnectAttempts=%d", c.parameters.ReconnectAttempts)
	log.Printf("ShutdownTimeout=%s", c.parameters.ShutdownTimeout)
	log.Printf("ResubscribeOnReconnect=%t", c.parameters.ResubscribeOnReconnect)
	log.Printf("HeartbeatTimeout=%s", c.parameters.HeartbeatTimeout)
	log.Printf("URL=%s", c.parameters.URL)
}

func (c *Client) reset() {
	if c.asynchronous == nil {
		c.asynchronous = c.asyncFactory.Create()
		c.init = true
	}
	c.updateHeartbeat()

	go c.listenDisconnect()

	go c.listenUpstream()

	go c.controlHeartbeat()
}

func (c *Client) connect() error {
	err := c.asynchronous.Connect()
	c.setIsConnected()
	return err
}

func (c *Client) resubscribe() error {
	for _, sub := range c.subscriptions {
		s := SubscriptionRequest{
			Event:        EventSubscribe,
			Pairs:        []string{sub.Pair},
			Subscription: sub.Subscription,
		}

		if err := c.asynchronous.Send(context.Background(), s); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (c *Client) reconnect(err error) error {
	if c.terminal {
		return c.exit(err)
	}
	if !c.parameters.AutoReconnect {
		s := "Empty error"
		if err != nil {
			s = err.Error()
		}
		return c.exit(fmt.Errorf("AutoReconnect setting is disabled, do not reconnect: %s", s))
	}
	reconnectTry := 0
	for ; reconnectTry < c.parameters.ReconnectAttempts; reconnectTry++ {
		log.Printf("waiting %s until reconnect...", c.parameters.ReconnectInterval)
		time.Sleep(c.parameters.ReconnectInterval)
		log.Printf("reconnect attempt %d/%d", reconnectTry+1, c.parameters.ReconnectAttempts)
		c.reset()

		if err = c.connect(); err == nil {
			if err = c.resubscribe(); err == nil {
				log.Print("reconnect OK")
				return nil
			}
			log.Printf("reconnect failed: %s", err.Error())
			return c.exit(err)
		}

		log.Printf("reconnect failed: %s", err.Error())
	}
	if err != nil {
		log.Printf("could not reconnect: %s", err.Error())
	}
	return c.exit(err)
}

func (c *Client) exit(err error) error {
	c.terminal = true
	c.close(err)
	return err
}

func (c *Client) listenUpstream() {
	for {
		select {
		case <-c.shutdown:
			log.Printf("Shutdown listen upstream")
			return
		case msg := <-c.asynchronous.Listen():
			if msg != nil {
				// log.Printf("[DEBUG]: %s\n", msg)
				if err := c.handleMessage(msg); err != nil {
					log.Printf("[WARN]: %s\n", err)
				}
			}
		}
	}
}

func (c *Client) close(e error) {
	if c.listener != nil {
		if e != nil {
			c.listener <- e
		}
		close(c.listener)
	}

	if c.shutdown != nil {
		close(c.shutdown)
	}
}

func (c *Client) closeAsyncAndWait(t time.Duration) {
	if !c.init {
		return
	}
	timeout := make(chan bool)
	c.waitGroup.Add(1)
	go func() {
		select {
		case <-c.asynchronous.Done():
			c.waitGroup.Done()
		case <-timeout:
			c.waitGroup.Done()
		}
	}()
	c.asynchronous.Close()
	go func() {
		time.Sleep(t)
		close(timeout)
	}()
	c.waitGroup.Wait()
}

func (c *Client) handleMessage(msg []byte) error {
	t := bytes.TrimLeftFunc(msg, unicode.IsSpace)
	c.updateHeartbeat()

	if len(t) == 0 {
		return fmt.Errorf("Empty response: %s", string(msg))
	}
	switch t[0] {
	case '[':
		return c.handleChannel(msg)
	case '{':
		return c.handleEvent(msg)
	default:
		return fmt.Errorf("Unexpected message: %s", string(msg))
	}
}

func (c *Client) createFactories() {
	c.createFactory(ChanTicker, newTickerFactory())
	c.createFactory(ChanCandles, newCandlesFactory())
	c.createFactory(ChanTrades, newTradesFactory())
	c.createFactory(ChanSpread, newSpreadFactory())
	c.createFactory(ChanBook, newBookFactory())
}

func (c *Client) createFactory(name string, factory ParseFactory) {
	c.factories[name] = factory
}

func (c *Client) handleEvent(msg []byte) error {
	var event EventType
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}

	switch event.Event {
	case EventPong:
		var pong PongResponse
		if err := json.Unmarshal(msg, &pong); err != nil {
			return err
		}
		log.Print("Pong received")

	case EventSystemStatus:
		var systemStatus SystemStatus
		if err := json.Unmarshal(msg, &systemStatus); err != nil {
			return err
		}
		log.Printf("Status: %s", systemStatus.Status)
		log.Printf("Connection ID: %s", systemStatus.ConnectionID.String())
		log.Printf("Version: %s", systemStatus.Version)

	case EventSubscriptionStatus:
		var status SubscriptionStatus
		if err := json.Unmarshal(msg, &status); err != nil {
			return err
		}

		if status.Status == SubscriptionStatusError {
			log.Printf("[ERROR] %s: %s", status.Error, status.Pair)
		} else {
			log.Printf("\tStatus: %s", status.Status)
			log.Printf("\tPair: %s", status.Pair)
			log.Printf("\tSubscription: %s", status.Subscription.Name)
			log.Printf("\tChannel ID: %d", status.ChannelID)
			log.Printf("\tReq ID: %s", status.ReqID)

			if status.Status == SubscriptionStatusSubscribed {
				c.subscriptions[status.ChannelID] = &status
			} else if status.Status == SubscriptionStatusUnsubscribed {
				delete(c.subscriptions, status.ChannelID)
			}
		}
	case EventCancelOrderStatus:
		var cancelOrderResponse CancelOrderResponse
		if err := json.Unmarshal(msg, &cancelOrderResponse); err != nil {
			return err
		}

		switch cancelOrderResponse.Status {
		case StatusError:
			log.Printf("[ERROR] %s", cancelOrderResponse.ErrorMessage)
		case StatusOK:
			log.Print("[INFO] Order successfully cancelled")
			c.listener <- DataUpdate{
				ChannelName: EventCancelOrder,
				Data:        cancelOrderResponse,
			}
		default:
			log.Printf("[ERROR] Unknown status: %s", cancelOrderResponse.Status)
		}
	case EventAddOrderStatus:
		var addOrderResponse AddOrderResponse
		if err := json.Unmarshal(msg, &addOrderResponse); err != nil {
			return err
		}

		switch addOrderResponse.Status {
		case StatusError:
			log.Printf("[ERROR] %s", addOrderResponse.ErrorMessage)
		case StatusOK:
			log.Print("[INFO] Order successfully sent")
			c.listener <- DataUpdate{
				ChannelName: EventAddOrder,
				Data:        addOrderResponse,
			}
		default:
			log.Printf("[ERROR] Unknown status: %s", addOrderResponse.Status)
		}
	case EventHeartbeat:
	default:
		fmt.Printf("unknown event: %s", msg)
	}
	return nil
}

func (c *Client) handleChannel(msg []byte) error {
	var data DataUpdate
	if err := json.Unmarshal(msg, &data); err != nil {
		return err
	}

	channel := strings.Split(data.ChannelName, "-")[0]
	factory, ok := c.factories[channel]
	if !ok {
		return fmt.Errorf("Unknown message type: %s", data.ChannelName)
	}

	result, err := factory.Parse(data.Data, data.Pair)
	if err != nil {
		return err
	}
	data.Data = result
	c.listener <- data
	return nil
}

// IsConnected returns true if the underlying asynchronous transport is connected to an endpoint.
func (c *Client) IsConnected() bool {
	c.isConnectedMux.Lock()
	defer c.isConnectedMux.Unlock()

	return c.isConnected
}

func (c *Client) setIsConnected() {
	c.isConnectedMux.Lock()
	defer c.isConnectedMux.Unlock()

	c.isConnected = true
}

// Connect to the Kraken API, this should only be called once.
func (c *Client) Connect() error {
	c.dumpParams()
	c.reset()
	return c.connect()
}

// Listen provides an atomic interface for receiving API messages.
// When a websocket connection is terminated, the publisher channel will close.
func (c *Client) Listen() <-chan interface{} {
	return c.listener
}

// Close provides an interface for a user initiated shutdown.
func (c *Client) Close() {
	c.terminal = true
	c.closeAsyncAndWait(c.parameters.ShutdownTimeout)

	timeout := make(chan bool)
	go func() {
		time.Sleep(c.parameters.ShutdownTimeout)
		close(timeout)
	}()
	select {
	case <-c.shutdown:
		return // successful cleanup
	case <-timeout:
		log.Print("shutdown timed out")
		return
	}
}

// SubscribeTicker - Ticker information includes best ask and best bid prices, 24hr volume, last trade price, volume weighted average price, etc for a given currency pair. A ticker message is published every time a trade or a group of trade happens.
func (c *Client) SubscribeTicker(pairs []string) error {
	ctx, cxl := context.WithTimeout(context.Background(), c.parameters.ContextTimeout)
	defer cxl()
	s := SubscriptionRequest{
		Event: EventSubscribe,
		Pairs: pairs,
		Subscription: Subscription{
			Name: ChanTicker,
		},
	}
	return c.asynchronous.Send(ctx, s)
}

// SubscribeCandles - Open High Low Close (Candle) feed for a currency pair and interval period.
func (c *Client) SubscribeCandles(pairs []string, interval int64) error {
	ctx, cxl := context.WithTimeout(context.Background(), c.parameters.ContextTimeout)
	defer cxl()
	s := SubscriptionRequest{
		Event: EventSubscribe,
		Pairs: pairs,
		Subscription: Subscription{
			Name:     ChanCandles,
			Interval: interval,
		},
	}
	return c.asynchronous.Send(ctx, s)
}

// SubscribeTrades - Trade feed for a currency pair.
func (c *Client) SubscribeTrades(pairs []string) error {
	ctx, cxl := context.WithTimeout(context.Background(), c.parameters.ContextTimeout)
	defer cxl()
	s := SubscriptionRequest{
		Event: EventSubscribe,
		Pairs: pairs,
		Subscription: Subscription{
			Name: ChanTrades,
		},
	}
	return c.asynchronous.Send(ctx, s)
}

// SubscribeSpread - Spread feed to show best bid and ask price for a currency pair
func (c *Client) SubscribeSpread(pairs []string) error {
	ctx, cxl := context.WithTimeout(context.Background(), c.parameters.ContextTimeout)
	defer cxl()
	s := SubscriptionRequest{
		Event: EventSubscribe,
		Pairs: pairs,
		Subscription: Subscription{
			Name: ChanSpread,
		},
	}
	return c.asynchronous.Send(ctx, s)
}

// SubscribeBook - Order book levels. On subscription, a snapshot will be published at the specified depth, following the snapshot, level updates will be published.
func (c *Client) SubscribeBook(pairs []string, depth int64) error {
	ctx, cxl := context.WithTimeout(context.Background(), c.parameters.ContextTimeout)
	defer cxl()
	s := SubscriptionRequest{
		Event: EventSubscribe,
		Pairs: pairs,
		Subscription: Subscription{
			Name:  ChanBook,
			Depth: depth,
		},
	}
	return c.asynchronous.Send(ctx, s)
}

// Unsubscribe - Unsubscribe from single subscription, can specify multiple currency pairs.
func (c *Client) Unsubscribe(channelType string, pairs []string) error {
	ctx, cxl := context.WithTimeout(context.Background(), c.parameters.ContextTimeout)
	defer cxl()
	u := UnsubscribeRequest{
		Event: EventUnsubscribe,
		Pairs: pairs,
		Subscription: Subscription{
			Name: channelType,
		},
	}
	return c.asynchronous.Send(ctx, u)
}

// Ping - Client can ping server to determine whether connection is alive, server responds with pong. This is an application level ping as opposed to default ping in WebSockets standard which is server initiated
func (c *Client) Ping() error {
	ctx, cxl := context.WithTimeout(context.Background(), c.parameters.ContextTimeout)
	defer cxl()
	ping := PingRequest{
		Event: EventPing,
	}
	return c.asynchronous.Send(ctx, ping)
}
