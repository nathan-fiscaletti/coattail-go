package types

import "encoding/gob"

func init() {
	gob.Register(Request{})
}

type Request struct {
	Password string `json:"password"`
}
