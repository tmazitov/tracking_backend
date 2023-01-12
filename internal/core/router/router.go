package router

import "github.com/gin-gonic/gin"

type Router struct {
	core *gin.Engine
}

func NewRouter() *Router {
	r := gin.Default()
	return &Router{core: r}
}

func (r *Router) Setup(endpoints []Endpoint) {
	for _, endpoint := range endpoints {
		r.core.Handle(endpoint.Method, endpoint.Path, endpoint.Handler.Action)
	}
}

func (r *Router) Run(port string) {
	r.core.Run()
}
