package rest

import (
	bl "github.com/tmazitov/tracking_backend.git/internal/aaa/bl"
	"github.com/tmazitov/tracking_backend.git/internal/core/router"
)

type Router struct {
	core        *router.Router
	storage     bl.Storage
	servicePath string
}

func NewRouter(servicePath string, storage bl.Storage) Router {
	core := router.NewRouter(servicePath)

	r := Router{
		core:        core,
		storage:     storage,
		servicePath: servicePath,
	}

	r.Setup()

	return r
}

func (r *Router) Endpoints() []router.Endpoint {
	return []router.Endpoint{
		{Method: "POST", Path: "/order", Handler: &AuthUser{Storage: r.storage}},
	}
}

func (r *Router) Setup() {
	r.core.Setup(r.Endpoints())
}

func (r *Router) Run(port string) {
	r.core.Run(port)
}
