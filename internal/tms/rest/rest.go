package rest

import (
	bl "github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	"github.com/tmazitov/tracking_backend.git/pkg/router"
)

type Router struct {
	jwt         *jwt.JwtStorage
	core        *router.Router
	storage     bl.Storage
	servicePath string
}

func NewRouter(servicePath string, storage bl.Storage, jwtStorage *jwt.JwtStorage) Router {
	core := router.NewRouter(servicePath)

	r := Router{
		jwt:         jwtStorage,
		core:        core,
		storage:     storage,
		servicePath: servicePath,
	}

	r.Setup()

	return r
}

func (r *Router) Endpoints() []router.Endpoint {
	return []router.Endpoint{
		{Method: "POST", Path: "/order", Handler: &OrderCreateHandler{Storage: r.storage, Jwt: *r.jwt}},
		{Method: "GET", Path: "/order/:orderId/upgrade", Handler: &OrderStatusUpgradeHandler{Storage: r.storage, Jwt: *r.jwt}},
		{Method: "PUT", Path: "/order/:orderId", Handler: &OrderPutHandler{Storage: r.storage, Jwt: *r.jwt}},
		{Method: "GET", Path: "/order/list", Handler: &OrderListHandler{Storage: r.storage, Jwt: *r.jwt}},
		{Method: "GET", Path: "/user", Handler: &UserGetHandler{Storage: r.storage, Jwt: *r.jwt}},
		{Method: "PUT", Path: "/user", Handler: &UserPutHandler{Storage: r.storage, Jwt: *r.jwt}},
	}
}

func (r *Router) Setup() {
	r.core.Setup(r.Endpoints())
}

func (r *Router) Run(port string) {
	r.core.Run(port)
}
