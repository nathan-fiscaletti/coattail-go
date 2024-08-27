package peers

type Action struct {
	Unit

	name string
}

func (a Action) Name() string {
	return a.name
}
