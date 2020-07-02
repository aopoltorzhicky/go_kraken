package websocket

import (
	"reflect"
	"testing"
	"time"
)

func TestNewDefaultParameters(t *testing.T) {
	tests := []struct {
		name string
		want *Parameters
	}{
		{
			name: "Test creation default parameters",
			want: &Parameters{
				AutoReconnect:          true,
				ReconnectInterval:      time.Second,
				reconnectTry:           0,
				ReconnectAttempts:      5,
				URL:                    ProdBaseURL,
				ShutdownTimeout:        time.Second * 5,
				ResubscribeOnReconnect: true,
				HeartbeatTimeout:       time.Second * 3,
				HeartbeatCheckPeriod:   time.Millisecond * 100,
				LogTransport:           false,
				ContextTimeout:         time.Second * 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultParameters(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultParameters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDefaultSandboxParameters(t *testing.T) {
	tests := []struct {
		name string
		want *Parameters
	}{
		{
			name: "Test creation default parameters for sandbox",
			want: &Parameters{
				AutoReconnect:          true,
				ReconnectInterval:      time.Second,
				reconnectTry:           0,
				ReconnectAttempts:      5,
				URL:                    SandboxBaseURL,
				ShutdownTimeout:        time.Second * 5,
				ResubscribeOnReconnect: true,
				HeartbeatTimeout:       time.Second * 3,
				HeartbeatCheckPeriod:   time.Millisecond * 100,
				LogTransport:           false,
				ContextTimeout:         time.Second * 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultSandboxParameters(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultSandboxParameters() = %v, want %v", got, tt.want)
			}
		})
	}
}
