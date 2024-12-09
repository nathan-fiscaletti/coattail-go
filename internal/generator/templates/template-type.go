package templates

type Template interface {
	Fill(dir string) error
}
