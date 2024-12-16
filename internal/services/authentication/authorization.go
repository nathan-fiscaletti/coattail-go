package authentication

// AuthorizedOperation is an operation that is authorized.
type AuthorizedOperation int

const (
	Run AuthorizedOperation = iota
	Publish
	Subscribe
	Notify
)

// AuthorizationType is the type of authorization.
type AuthorizationType int

const (
	Action AuthorizationType = iota
	Receiver
)

// Authorization is a set of authorized operations.
type Authorization struct {
	// Type is the type of authorization, either action or receiver.
	Type AuthorizationType
	// Operations is the set of authorized operations.
	Operations []AuthorizedOperation
	// Name is the name of the action or receiver.
	Name string
}

// AuthorizationRequest is a request to check authorization.
type AuthorizationRequest struct {
	// Type is the type of authorization.
	Type AuthorizationType
	// Operation is the operation to check.
	Operation AuthorizedOperation
	// Name is the name of the action or receiver.
	Name string
}
