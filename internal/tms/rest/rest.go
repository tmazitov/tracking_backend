package rest

import (
	"github.com/redis/go-redis/v9"
	bl "github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/internal/tms/ws"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	"github.com/tmazitov/tracking_backend.git/pkg/router"
)

type Router struct {
	jwt         *jwt.JwtStorage
	core        *router.Router
	storage     bl.Storage
	hub         *ws.Hub
	servicePath string
}

func NewRouter(servicePath string, redis *redis.Client, storage bl.Storage, jwtStorage *jwt.JwtStorage) Router {
	core := router.NewRouter(servicePath)

	hub := ws.NewHub(storage, redis, jwtStorage)

	r := Router{
		hub:         hub,
		jwt:         jwtStorage,
		core:        core,
		storage:     storage,
		servicePath: servicePath,
	}

	r.Setup()
	go hub.Run()
	go hub.RunOrderDispatcher()

	return r
}

func (r *Router) Endpoints() []router.Endpoint {
	return []router.Endpoint{
		{Method: "PATCH", Path: "/order/:orderId/worker/update", Handler: &OrderSetWorkerHandler{Storage: r.storage, Jwt: *r.jwt, Hub: r.hub}},
		{Method: "POST", Path: "/order", Handler: &OrderCreateHandler{Storage: r.storage, Jwt: *r.jwt}},
		{Method: "GET", Path: "/order/:orderId/upgrade", Handler: &OrderStatusUpgradeHandler{Storage: r.storage, Jwt: *r.jwt}},
		{Method: "GET", Path: "/order/:orderId/start", Handler: &OrderTimeStartHandler{Storage: r.storage, Jwt: *r.jwt, Hub: r.hub}},
		{Method: "GET", Path: "/order/:orderId/end", Handler: &OrderTimeEndHandler{Storage: r.storage, Jwt: *r.jwt, Hub: r.hub}},
		{Method: "PUT", Path: "/order/:orderId", Handler: &OrderPutHandler{Storage: r.storage, Jwt: *r.jwt}},
		{Method: "GET", Path: "/order/list", Handler: &OrderListHandler{Storage: r.storage, Jwt: *r.jwt}},

		{Method: "GET", Path: "/user", Handler: &UserGetHandler{Storage: r.storage, Jwt: *r.jwt}},
		{Method: "PUT", Path: "/user", Handler: &UserPutHandler{Storage: r.storage, Jwt: *r.jwt}},

		{Method: "GET", Path: "/user/:userId/holiday", Handler: &UserHolidayCreateHandler{Storage: r.storage, Jwt: r.jwt}},
		{Method: "DELETE", Path: "/user/:userId/holiday", Handler: &UserHolidayDeleteHandler{Storage: r.storage, Jwt: r.jwt}},
		{Method: "GET", Path: "/holiday-list", Handler: &UserHolidayListByDate{Storage: r.storage, Jwt: r.jwt}},

		{Method: "GET", Path: "/staff", Handler: &StaffListHandler{Storage: r.storage, Jwt: *r.jwt}},
		{Method: "GET", Path: "/order/price-list", Handler: &OrderDefaultVariables{Storage: r.storage, Jwt: *r.jwt}},

		{Method: "GET", Path: "/order/updates", Handler: &ws.WS_OrderListHandler{Hub: r.hub, Storage: r.storage, Jwt: r.jwt}, WS: true},

		// TODO : make worker holiday CRUD
	}
}

func (r *Router) Setup() {
	r.core.Setup(r.Endpoints())
}

func (r *Router) Run(port string) {
	r.core.Run(port)
}
