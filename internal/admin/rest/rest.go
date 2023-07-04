package rest

import (
	"github.com/gin-gonic/gin"
	bl "github.com/tmazitov/tracking_backend.git/internal/admin/bl"
	"github.com/tmazitov/tracking_backend.git/internal/admin/middleware"
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

	adminRoleMiddleware := middleware.AdminRoleMiddleware{Jwt: jwtStorage}
	var middleware []gin.HandlerFunc = []gin.HandlerFunc{
		adminRoleMiddleware.Handle,
	}

	r.AddMiddleware(middleware)
	r.Setup()

	return r
}

func (r *Router) Endpoints() []router.Endpoint {
	return []router.Endpoint{
		{Method: "GET", Path: "/offer/list", Handler: &OfferListHandler{Storage: r.storage}},
		{Method: "GET", Path: "/offer/:offerId/accept", Handler: &OfferAcceptHandler{Storage: r.storage}},
		{Method: "GET", Path: "/offer/:offerId/reject", Handler: &OfferRejectHandler{Storage: r.storage}},
		{Method: "POST", Path: "/staff/remove", Handler: &StaffRemoveHandler{Storage: r.storage}},
		{Method: "PUT", Path: "/staff/work-time", Handler: &StaffWorkTimePut{Storage: r.storage}},
		{Method: "PUT", Path: "/order/price-list", Handler: &OrderPriceListPutHandler{Storage: r.storage}},
	}
}

func (r *Router) Setup() {
	r.core.Setup(r.Endpoints())
}

func (r *Router) AddMiddleware(middleware []gin.HandlerFunc) {
	r.core.AddMiddleware(middleware)
}

func (r *Router) Run(port string) {
	r.core.Run(port)
}
