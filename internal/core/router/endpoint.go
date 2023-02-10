package router

import "github.com/tmazitov/tracking_backend.git/internal/core/middleware"

type Endpoint struct {
	Method     string
	Path       string
	Middleware middleware.Middleware
	Handler    Handler
}
