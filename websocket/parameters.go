package websocket

import (
	"time"
)

// Parameters defines adapter behavior.
type Parameters struct {
	URL string

	ReconnectInterval    time.Duration
	ShutdownTimeout      time.Duration
	ContextTimeout       time.Duration
	HeartbeatCheckPeriod time.Duration
	HeartbeatTimeout     time.Duration

	ReconnectAttempts      int
	reconnectTry           int
	LogTransport           bool
	AutoReconnect          bool
	ResubscribeOnReconnect bool
}

// NewDefaultParameters - create default Parameters object for prod
func NewDefaultParameters() *Parameters {
	return &Parameters{
		AutoReconnect:          true,
		ReconnectInterval:      time.Second,
		reconnectTry:           0,
		ReconnectAttempts:      5,
		URL:                    ProdBaseURL,
		ShutdownTimeout:        time.Second * 5,
		ResubscribeOnReconnect: true,
		HeartbeatCheckPeriod:   time.Millisecond * 100,
		HeartbeatTimeout:       time.Second * 3, // HB = 3s
		LogTransport:           false,           // log transport send/recv,
		ContextTimeout:         time.Second * 5,
	}
}

// NewDefaultSandboxParameters - create default Parameters object for sandbox
func NewDefaultSandboxParameters() *Parameters {
	return &Parameters{
		AutoReconnect:          true,
		ReconnectInterval:      time.Second,
		reconnectTry:           0,
		ReconnectAttempts:      5,
		URL:                    SandboxBaseURL,
		ShutdownTimeout:        time.Second * 5,
		ResubscribeOnReconnect: true,
		HeartbeatTimeout:       time.Second * 3, // HB = 3s
		HeartbeatCheckPeriod:   time.Millisecond * 100,
		LogTransport:           false, // log transport send/recv
		ContextTimeout:         time.Second * 5,
	}
}

// NewDefaultAuthParameters - create default Parameters object for auth socket
func NewDefaultAuthParameters() *Parameters {
	return &Parameters{
		AutoReconnect:          true,
		ReconnectInterval:      time.Second,
		reconnectTry:           0,
		ReconnectAttempts:      5,
		URL:                    AuthBaseURL,
		ShutdownTimeout:        time.Second * 5,
		ResubscribeOnReconnect: true,
		HeartbeatTimeout:       time.Second * 3, // HB = 3s
		HeartbeatCheckPeriod:   time.Millisecond * 100,
		LogTransport:           false, // log transport send/recv,
		ContextTimeout:         time.Second * 5,
	}
}

// NewDefaultSandboxAuthParameters - create default Parameters object for sanbox auth
func NewDefaultSandboxAuthParameters() *Parameters {
	return &Parameters{
		AutoReconnect:          true,
		ReconnectInterval:      time.Second,
		reconnectTry:           0,
		ReconnectAttempts:      5,
		URL:                    AuthSandboxBaseURL,
		ShutdownTimeout:        time.Second * 5,
		ResubscribeOnReconnect: true,
		HeartbeatTimeout:       time.Second * 3, // HB = 3s
		HeartbeatCheckPeriod:   time.Millisecond * 100,
		LogTransport:           false, // log transport send/recv,
		ContextTimeout:         time.Second * 5,
	}
}
