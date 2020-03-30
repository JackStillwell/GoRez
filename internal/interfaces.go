package gorez_internal

// httpGetter is an interface for mocking
type httpGetter interface {
	get(url string) ([]byte, error)
}
