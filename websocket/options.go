package websocket

// AuthOption - option function for `AuthClient`
type AuthOption func(*AuthClient)

// WithParams - add custom params to `AuthClient`
func WithParams(params *Parameters) AuthOption {
	return func(auth *AuthClient) {
		if params != nil {
			auth.parameters = params
		}
	}
}
