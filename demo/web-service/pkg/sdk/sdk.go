// Code generated by coattail; DO NOT EDIT.
package sdk

import (
    "github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

type Sdk struct {
	peer *coattailtypes.Peer
}

func NewSdk(peer *coattailtypes.Peer) *Sdk {
	return &Sdk{
		peer: peer,
	}
}