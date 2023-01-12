package router

type Endpoint struct {
	Method  string
	Path    string
	Handler Handler
}
