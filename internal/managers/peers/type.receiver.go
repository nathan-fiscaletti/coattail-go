package peers

type Receiver struct {
	Unit

	name string
}

func (r Receiver) Name() string {
	return r.name
}
