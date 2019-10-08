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
	isConnected  bool
	terminal     bool
	init         bool
	heartbeat    time.Time
	hbChannel    chan error

	// connection & operational behavior
	parameters *Parameters

	// close signal sent to user on shutdown
	shutdown chan bool

	// downstream listener channel to deliver API objects
	listener chan interface{}

	// race management
	waitGroup sync.WaitGroup

	subscriptions map[int64]*SubscriptionStatus
	factories     map[string]ParseFactory
}

// Create returns a new websocket transport.
func (w *websocketAsynchronousFactory) Create() asynchronous {
	return newWs(w.parameters.URL, w.parameters.LogTransport)
}

// New creates a default client.
func New(sandbox bool) *Client {
	var params *Parameters
	if sandbox {
		params = NewDefaultSandboxParameters()
	} else {
		params = NewDefaultParameters()
	}
	c := &Client{
		asyncFactory: &websocketAsynchronousFactory{
			parameters: params,
		},
		isConnected:   false,
		parameters:    params,
		listener:      make(chan interface{}),
		terminal:      false,
		shutdown:      nil,
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
		if time.Now().After(c.heartbeat) {
			c.hbChannel <- fmt.Errorf("Heartbeat disconnect on client")
			return
		}
	}
}

func (c *Client) listenHeartbeat() <-chan error {
	return c.hbChannel
}

func (c *Client) updateHeartbeat() {
	c.heartbeat = time.Now().Add(c.parameters.HeartbeatTimeout)
}

func (c *Client) listenDisconnect() {
	select {
	case e := <-c.asynchronous.Done():
		log.Printf("socket disconnect")
		if e != nil {
			log.Printf("socket disconnect: %s", e.Error())
		}
		c.isConnected = false
		err := c.reconnect(e)
		if err != nil {
			log.Printf("socket disconnect: %s", err.Error())
		}
	case <-c.shutdown:
		log.Printf("Shutdown listen disconnect")
		c.isConnected = false
		return

	case e := <-c.listenHeartbeat():
		log.Printf("Heartbeat")
		if e != nil {
			log.Println(e.Error())
			c.closeAsyncAndWait(c.parameters.ShutdownTimeout)
			err := c.reconnect(nil)
			if err != nil {
				log.Printf("socket disconnect: %s", err.Error())
			}
		}
		c.isConnected = false
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
	if c.shutdown != nil {
		close(c.shutdown)
	}

	c.init = true
	c.asynchronous = c.asyncFactory.Create()
	c.shutdown = make(chan bool)
	c.updateHeartbeat()

	go c.listenDisconnect()

	go c.listenUpstream()

	go c.controlHeartbeat()
}

func (c *Client) connect() error {
	err := c.asynchronous.Connect()
	if err == nil {
		c.isConnected = true
	}
	return err
}

func (c *Client) resubscribe() error {
	ctx := context.Background()
	for _, sub := range c.subscriptions {
		s := SubscriptionRequest{
			Event:        EventSubscribe,
			Pairs:        []string{sub.Pair},
			Subscription: sub.Subscription,
		}

		err := c.asynchronous.Send(ctx, s)
		if err != nil {
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
		err = c.connect()
		if err == nil {
			err := c.resubscribe()
			if err == nil {
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
				log.Printf("[DEBUG]: %s\n", msg)
				err := c.handleMessage(msg)
				if err != nil {
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
	close(c.shutdown)
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

func (c *Client) handleMessage(msg []byte) (err error) {
	t := bytes.TrimLeftFunc(msg, unicode.IsSpace)
	c.updateHeartbeat()
	if bytes.HasPrefix(t, []byte("[")) {
		err = c.handleChannel(msg)
	} else if bytes.HasPrefix(t, []byte("{")) {
		err = c.handleEvent(msg)
	} else {
		return fmt.Errorf("unexpected message: %s", msg)
	}
	return err
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
	event := &EventType{}
	err := json.Unmarshal(msg, event)
	if err != nil {
		return err
	}
	switch event.Event {

	case EventPong:
		pong := PongResponse{}
		err = json.Unmarshal(msg, &pong)
		if err != nil {
			return err
		}
		log.Print("Pong received")

	case EventSystemStatus:
		systemStatus := SystemStatus{}
		err = json.Unmarshal(msg, &systemStatus)
		if err != nil {
			return err
		}
		log.Printf("Status: %s", systemStatus.Status)
		log.Printf("Connection ID: %s", systemStatus.ConnectionID.String())
		log.Printf("Version: %s", systemStatus.Version)

	case EventSubscriptionStatus:
		status := SubscriptionStatus{}
		err = json.Unmarshal(msg, &status)
		if err != nil {
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
		cancelOrderResponse := CancelOrderResponse{}
		err = json.Unmarshal(msg, &cancelOrderResponse)
		if err != nil {
			return err
		}

		if cancelOrderResponse.Status == StatusError {
			log.Printf("[ERROR] %s", cancelOrderResponse.ErrorMessage)
		} else if cancelOrderResponse.Status == StatusOK {
			log.Print("[INFO] Order successfully cancelled")
			c.listener <- DataUpdate{
				ChannelName: EventCancelOrder,
				Data:        cancelOrderResponse,
			}
		} else {
			log.Printf("[ERROR] Unknown status: %s", cancelOrderResponse.Status)
		}
	case EventAddOrderStatus:
		addOrderResponse := AddOrderResponse{}
		err = json.Unmarshal(msg, &addOrderResponse)
		if err != nil {
			return err
		}

		if addOrderResponse.Status == StatusError {
			log.Printf("[ERROR] %s", addOrderResponse.ErrorMessage)
		} else if addOrderResponse.Status == StatusOK {
			log.Print("[INFO] Order successfully sent")
			c.listener <- DataUpdate{
				ChannelName: EventAddOrder,
				Data:        addOrderResponse,
			}
		} else {
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
	err := json.Unmarshal(msg, &data)
	if err != nil {
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
	return c.isConnected
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
