package coattailmodels

import (
	"fmt"

	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	// The Address to notify when the action is published.
	Address string `json:"address"`

	// Action represents the action on the local instance
	// that the subscriber is subscribed to.
	Action string `json:"action"`

	// Receiver represents the receiver on the remote
	// instance that will be notified when the action
	// is published.
	Receiver string `json:"receiver"`
}

func (s Subscription) String() string {
	return fmt.Sprintf("%s -> %s//%s", s.Action, s.Address, s.Receiver)
}
