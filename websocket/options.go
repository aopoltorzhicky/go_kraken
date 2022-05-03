package websocket

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// KrakenOption - option function for `AuthClient`
type KrakenOption func(*Kraken)

// WithLogLevel - add custom log level. Default: info
func WithLogLevel(level log.Level) KrakenOption {
	return func(k *Kraken) {
		log.SetLevel(level)
	}
}

// WithReconnectTimeout - add custom reconnect timeout (time interval for next reconnecting try). Default: 5s.
func WithReconnectTimeout(timeout time.Duration) KrakenOption {
	return func(k *Kraken) {
		k.reconnectTimeout = timeout
	}
}

// WithReadTimeout - add custom read timeout. Default: 15s.
func WithReadTimeout(timeout time.Duration) KrakenOption {
	return func(k *Kraken) {
		k.readTimeout = timeout
	}
}

// WithHeartbeatTimeout - add custom heartbeat timeout (time interval for sending ping message). Default: 10s.
func WithHeartbeatTimeout(timeout time.Duration) KrakenOption {
	return func(k *Kraken) {
		k.heartbeatTimeout = timeout
	}
}
