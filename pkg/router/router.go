package router

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	core        *gin.Engine
	servicePath string
}

func NewRouter(servicePath string) *Router {
	r := gin.Default()
	return &Router{core: r, servicePath: servicePath}
}

func (r *Router) Setup(endpoints []Endpoint) {
	for _, endpoint := range endpoints {
		if endpoint.Middleware != nil {
			r.core.Handle(endpoint.Method, r.servicePath+endpoint.Path, endpoint.Middleware.Handle(), endpoint.Handler.Handle)
		} else {
			r.core.Handle(endpoint.Method, r.servicePath+endpoint.Path, endpoint.Handler.Handle)
		}
	}
}

func (r *Router) Run(port string) {
	r.core.Run(":" + port)
}
