package coattailmodels

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	ID uuid.UUID

	// TODO: Subscriber ID transmission / identification
	// SubscriberID represents the ID of the peer that
	// is subscribed to the action.
	SubscriberID uuid.UUID

	// Action represents the action on the local instance
	// that the subscriber is subscribed to.
	Action string

	// Receiver represents the receiver on the remote
	// instance that will be notified when the action
	// is published.
	Receiver string
}
