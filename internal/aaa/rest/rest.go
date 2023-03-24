package rest

import (
	bl "github.com/tmazitov/tracking_backend.git/internal/aaa/bl"
	"github.com/tmazitov/tracking_backend.git/internal/core/conductor"
	"github.com/tmazitov/tracking_backend.git/internal/core/jwt"
	"github.com/tmazitov/tracking_backend.git/internal/core/router"
)

type Router struct {
	jwt         *jwt.JwtStorage
	core        *router.Router
	storage     bl.Storage
	conductor   conductor.Conductor
	servicePath string
}

func NewRouter(servicePath string, storage bl.Storage, conductor conductor.Conductor, jwtStorage *jwt.JwtStorage) Router {
	core := router.NewRouter(servicePath)

	r := Router{
		jwt:         jwtStorage,
		core:        core,
		storage:     storage,
		conductor:   conductor,
		servicePath: servicePath,
	}

	r.Setup()

	return r
}

func (r *Router) Endpoints() []router.Endpoint {

	var (
	// authMiddleware = &auth.Middleware{Jwt: r.jwt}
	)

	return []router.Endpoint{
		{Method: "POST", Path: "/auth", Handler: &AuthUserSendCode{Storage: r.storage, Conductor: r.conductor}},
		{Method: "POST", Path: "/auth/code", Handler: &AuthUserTakeCode{Storage: r.storage, Conductor: r.conductor}},

		{Method: "POST", Path: "/refresh", Handler: &RefreshHandler{Storage: r.storage, Jwt: r.jwt}},
	}
}

func (r *Router) Setup() {
	r.core.Setup(r.Endpoints())
}

func (r *Router) Run(port string) {
	r.core.Run(port)
}
