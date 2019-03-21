package websocket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
	"unicode"
)

// Asynchronous interface decouples the underlying transport from API logic.
type Asynchronous interface {
	Connect() error
	Send(ctx context.Context, msg interface{}) error
	Listen() <-chan []byte
	Close()
	Done() <-chan error
}

// AsynchronousFactory provides an interface to re-create asynchronous transports during reconnect events.
type AsynchronousFactory interface {
	Create() Asynchronous
}

// WebsocketAsynchronousFactory creates a websocket-based asynchronous transport.
type WebsocketAsynchronousFactory struct {
	parameters *Parameters
}

// NewWebsocketAsynchronousFactory creates a new websocket factory with a given URL.
func NewWebsocketAsynchronousFactory(parameters *Parameters) AsynchronousFactory {
	return &WebsocketAsynchronousFactory{
		parameters: parameters,
	}
}

// Create returns a new websocket transport.
func (w *WebsocketAsynchronousFactory) Create() Asynchronous {
	return newWs(w.parameters.URL, w.parameters.LogTransport)
}

// Client provides a unified interface for users to interact with the Kraken Websocket API.
// nolint:megacheck,structcheck
type Client struct {
	asyncFactory AsynchronousFactory // for re-creating transport during reconnects

	timeout            int64 // read timeout
	cancelOnDisconnect bool
	asynchronous       Asynchronous
	isConnected        bool
	terminal           bool
	init               bool
	heartbeat          time.Time

	// connection & operational behavior
	parameters *Parameters

	// close signal sent to user on shutdown
	shutdown chan bool

	// downstream listener channel to deliver API objects
	listener chan interface{}

	// race management
	lock      sync.Mutex
	waitGroup sync.WaitGroup

	subscriptions map[int64]*SubscriptionStatus
}

// New creates a default client.
func New() *Client {
	params := NewDefaultParameters()
	c := &Client{
		asyncFactory: &WebsocketAsynchronousFactory{
			parameters: params,
		},
		isConnected:   false,
		parameters:    params,
		listener:      make(chan interface{}),
		terminal:      false,
		shutdown:      nil,
		asynchronous:  nil,
		heartbeat:     time.Now().Add(params.HeartbeatTimeout),
		subscriptions: make(map[int64]*SubscriptionStatus),
	}
	c.registerPublicFactories()
	return c
}

// NewSandbox creates a default sandbox client.
func NewSandbox() *Client {
	params := NewDefaultSandboxParameters()
	c := &Client{
		asyncFactory: &WebsocketAsynchronousFactory{
			parameters: params,
		},
		isConnected:   false,
		parameters:    params,
		listener:      make(chan interface{}),
		terminal:      false,
		shutdown:      nil,
		asynchronous:  nil,
		heartbeat:     time.Now().Add(params.HeartbeatTimeout),
		subscriptions: make(map[int64]*SubscriptionStatus),
	}
	c.registerPublicFactories()
	return c
}

// CancelOnDisconnect ensures all orders will be canceled if this API session is disconnected.
func (c *Client) CancelOnDisconnect(cxl bool) *Client {
	c.cancelOnDisconnect = cxl
	return c
}

func (c *Client) registerPublicFactories() {

}

// IsConnected returns true if the underlying asynchronous transport is connected to an endpoint.
func (c *Client) IsConnected() bool {
	return c.isConnected
}

func (c *Client) listenDisconnect() {
	select {
	case e := <-c.asynchronous.Done(): // transport shutdown
		if e != nil {
			log.Printf("socket disconnect: %s", e.Error())
		}
		c.isConnected = false
		err := c.reconnect(e)
		if err != nil {
			log.Printf("socket disconnect: %s", err.Error())
		}
	case <-c.shutdown: // normal shutdown
		c.isConnected = false
		return // exit routine
	default:
		if time.Now().After(c.heartbeat) { // subscription heartbeat timeout
			log.Printf("heartbeat disconnect")
			c.isConnected = false
			c.closeAsyncAndWait(c.parameters.ShutdownTimeout)
			err := c.reconnect(nil)
			if err != nil {
				log.Printf("socket disconnect: %s", err.Error())
			}
		}
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
	log.Printf("ManageOrderbook=%t", c.parameters.ManageOrderbook)
}

// Connect to the Bitfinex API, this should only be called once.
func (c *Client) Connect() error {
	c.dumpParams()
	c.reset()
	return c.connect()
}

// reset assumes transport has already died or been closed
func (c *Client) reset() {
	// subs := c.subscriptions.Reset()
	// shutown if existing websocket connected
	if c.shutdown != nil {
		close(c.shutdown)
	}

	// if subs != nil {
	// 	c.resetSubscriptions = subs
	// }
	c.init = true
	c.asynchronous = c.asyncFactory.Create()
	c.shutdown = make(chan bool)

	// wait for shutdown signals from child & caller
	go c.listenDisconnect()
	// listen to data from async
	go c.listenUpstream()
}

func (c *Client) connect() error {
	err := c.asynchronous.Connect()
	if err == nil {
		c.isConnected = true
	}
	if c.parameters.ManageOrderbook {
		_, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
	}
	return err
}

func (c *Client) reconnect(err error) error {
	if c.terminal {
		err_exit := c.exit(err)
		if err_exit != nil {
			return err_exit
		}
		return err
	}
	if !c.parameters.AutoReconnect {
		err := fmt.Errorf("AutoReconnect setting is disabled, do not reconnect: %s", err.Error())
		err_exit := c.exit(err)
		if err_exit != nil {
			return err_exit
		}
		return err
	}
	reconnectTry := 0
	for ; reconnectTry < c.parameters.ReconnectAttempts; reconnectTry++ {
		log.Printf("waiting %s until reconnect...", c.parameters.ReconnectInterval)
		time.Sleep(c.parameters.ReconnectInterval)
		log.Printf("reconnect attempt %d/%d", reconnectTry+1, c.parameters.ReconnectAttempts)
		c.reset()
		err = c.connect()
		if err == nil {
			log.Print("reconnect OK")
			reconnectTry = 0
			return nil
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

// start this goroutine before connecting, but this should die during a connection failure
func (c *Client) listenUpstream() {
	for {
		select {
		case <-c.shutdown:
			return // only exit point
		case msg := <-c.asynchronous.Listen():
			if msg != nil {
				// Errors here should be non critical so we just log them.
				log.Printf("[DEBUG]: %s\n", msg)
				err := c.handleMessage(msg)
				if err != nil {
					log.Printf("[WARN]: %s\n", err)
				}
			}
		}
	}
}

// terminal, unrecoverable state. called after async is closed.
func (c *Client) close(e error) {
	if c.listener != nil {
		if e != nil {
			c.listener <- e
		}
		close(c.listener)
	}
	// shutdowns goroutines
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

// Listen provides an atomic interface for receiving API messages.
// When a websocket connection is terminated, the publisher channel will close.
func (c *Client) Listen() <-chan interface{} {
	return c.listener
}

// Close provides an interface for a user initiated shutdown.
// Close will close the Done() channel.
func (c *Client) Close() {
	c.terminal = true
	c.closeAsyncAndWait(c.parameters.ShutdownTimeout)
	// c.subscriptions.Close()

	// clean shutdown waits on shutdown channel, which is triggered by cascading resource
	// cleanups after a closed asynchronous transport
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

func (c *Client) handleMessage(msg []byte) error {
	t := bytes.TrimLeftFunc(msg, unicode.IsSpace)
	err := error(nil)
	// either a channel data array or an event object, raw json encoding
	if bytes.HasPrefix(t, []byte("[")) {
		err = c.handleChannel(msg)
	} else if bytes.HasPrefix(t, []byte("{")) {
		err = c.handleEvent(msg)
	} else {
		return fmt.Errorf("unexpected message: %s", msg)
	}
	return err
}

func (c *Client) handleEvent(msg []byte) error {
	event := &EventType{}
	err := json.Unmarshal(msg, event)
	if err != nil {
		return err
	}
	switch event.Event {

	case EventPong:
		pong := Pong{}
		err = json.Unmarshal(msg, &pong)
		if err != nil {
			return err
		}
		log.Printf("Pong received. Req ID: %s", pong.ReqID)

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
	case EventHeartbeat:

	default:
		fmt.Printf("unknown event: %s", msg)
	}
	return nil
}

func (c *Client) newSubscription(ctx context.Context, s SubscriptionRequest) error {
	return c.asynchronous.Send(ctx, s)
}

func (c *Client) SubscribeTicker(ctx context.Context, pairs []string) error {
	s := SubscriptionRequest{
		Event: EventSubscribe,
		Pairs: pairs,
		Subscription: Subscription{
			Name: ChanTicker,
		},
	}
	return c.newSubscription(ctx, s)
}

func (c *Client) lookupByChannelID(chanID int64) (*SubscriptionStatus, error) {
	sub, ok := c.subscriptions[chanID]
	if ok {
		return sub, nil
	}
	return nil, fmt.Errorf("Unknown channel ID: %s", chanID)
}

func (c *Client) handleChannel(msg []byte) error {
	var raw []interface{}
	err := json.Unmarshal(msg, &raw)
	if err != nil {
		return err
	} else if len(raw) < 2 {
		return nil
	}

	chID, ok := raw[0].(float64)
	if !ok {
		return fmt.Errorf("expected message to start with a channel id but got %#v instead", raw[0])
	}
	chanID := int64(chID)
	sub, err := c.lookupByChannelID(chanID)
	if err != nil {
		// no subscribed channel for message
		return err
	}

	switch sub.Subscription.Name {
	case ChanBook:
	case ChanCandles:
	case ChanSpread:
	case ChanTicker:
	case ChanTrades:
	default:
		return fmt.Errorf("Unknown message type: %s", sub.Subscription.Name)
	}
	return nil
}
