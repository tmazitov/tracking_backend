package router

import "github.com/tmazitov/tracking_backend.git/pkg/middleware"

type Endpoint struct {
	Method     string
	WS         bool
	Path       string
	Middleware middleware.Middleware
	Handler    Handler
}
