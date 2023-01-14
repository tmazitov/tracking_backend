package router

import (
	"fmt"

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
		fmt.Printf("Setup endpoint: %s | %s \n", endpoint.Method, r.servicePath+endpoint.Path)
		r.core.Handle(endpoint.Method, r.servicePath+endpoint.Path, endpoint.Handler.Handle)
	}
}

func (r *Router) Run(port string) {
	r.core.Run(":" + port)
}
