package gorez

// HTTPGetter is an interface for mocking
type HTTPGetter interface {
	Get(url string) ([]byte, error)
}
