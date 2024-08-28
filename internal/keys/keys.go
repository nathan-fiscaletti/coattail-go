package keys

type ContextKey int

const (
	AuthenticationServiceKey ContextKey = iota
	PeerManagerKey
)
