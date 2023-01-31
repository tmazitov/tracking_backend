package rest

import (
	bl "github.com/tmazitov/tracking_backend.git/internal/aaa/bl"
	"github.com/tmazitov/tracking_backend.git/internal/core/conductor"
	"github.com/tmazitov/tracking_backend.git/internal/core/router"
)

type Router struct {
	core        *router.Router
	storage     bl.Storage
	conductor   conductor.Conductor
	servicePath string
}

func NewRouter(servicePath string, storage bl.Storage, conductor conductor.Conductor) Router {
	core := router.NewRouter(servicePath)

	r := Router{
		core:        core,
		storage:     storage,
		conductor:   conductor,
		servicePath: servicePath,
	}

	r.Setup()

	return r
}

func (r *Router) Endpoints() []router.Endpoint {
	return []router.Endpoint{
		{Method: "POST", Path: "/auth", Handler: &AuthUserSendCode{Storage: r.storage, Conductor: r.conductor}},
		{Method: "POST", Path: "/auth/code", Handler: &AuthUserTakeCode{Storage: r.storage, Conductor: r.conductor}},
	}
}

func (r *Router) Setup() {
	r.core.Setup(r.Endpoints())
}

func (r *Router) Run(port string) {
	r.core.Run(port)
}
