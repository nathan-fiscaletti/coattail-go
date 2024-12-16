package authentication

import (
	"net"
	"time"

	"github.com/nathan-fiscaletti/coattail-go/internal/services/permission"
)

// Claims is a set of claims that can be used to issue a token.
type Claims struct {
	AuthorizedNetwork net.IPNet
	Permitted         int32
	Authorizations    []Authorization
	Expiry            time.Time
}

// Permissions returns the permissions of the claims.
func (c Claims) Permissions() permission.Permissions {
	return permission.GetPermissions(c.Permitted)
}

// IsAuthorized checks if the claims are authorized for the provided request.
func (c Claims) IsAuthorized(req AuthorizationRequest) bool {
	for _, a := range c.Authorizations {
		if a.Type == req.Type && (a.Name == req.Name || a.Name == "") {
			for _, operation := range a.Operations {
				if operation == req.Operation {
					return true
				}
			}
		}
	}
	return false
}
