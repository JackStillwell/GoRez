package internal_gorez

// httpGetter is an interface for mocking
type httpGetter interface {
	get(url string) ([]byte, error)
}
