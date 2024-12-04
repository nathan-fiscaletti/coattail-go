package keys

type ContextKey int

const (
	AuthenticationServiceKey ContextKey = iota
	HostKey
	LoggerKey
	DatabaseKey
	ConnectionKey
	AuthenticationKey
)
