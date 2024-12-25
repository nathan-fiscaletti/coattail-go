package types

import "encoding/gob"

func init() {
	gob.Register(Response{})
}

type Response struct {
	Authenticated bool `json:"authenticated"`
}
