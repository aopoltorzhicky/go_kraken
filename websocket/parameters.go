package websocket

import (
	"time"
)

// Parameters defines adapter behavior.
type Parameters struct {
	AutoReconnect     bool
	ReconnectInterval time.Duration
	ReconnectAttempts int
	reconnectTry      int
	ShutdownTimeout   time.Duration

	ResubscribeOnReconnect bool

	HeartbeatTimeout time.Duration
	LogTransport     bool

	URL             string
	ManageOrderbook bool
}

func NewDefaultParameters() *Parameters {
	return &Parameters{
		AutoReconnect:          true,
		ReconnectInterval:      time.Second,
		reconnectTry:           0,
		ReconnectAttempts:      5,
		URL:                    prodBaseURL,
		ManageOrderbook:        false,
		ShutdownTimeout:        time.Second * 5,
		ResubscribeOnReconnect: true,
		HeartbeatTimeout:       time.Second * 3, // HB = 3s
		LogTransport:           false,           // log transport send/recv
	}
}

func NewDefaultSandboxParameters() *Parameters {
	return &Parameters{
		AutoReconnect:          true,
		ReconnectInterval:      time.Second,
		reconnectTry:           0,
		ReconnectAttempts:      5,
		URL:                    sandboxBaseURL,
		ManageOrderbook:        false,
		ShutdownTimeout:        time.Second * 5,
		ResubscribeOnReconnect: true,
		HeartbeatTimeout:       time.Second * 3, // HB = 3s
		LogTransport:           false,           // log transport send/recv
	}
}
