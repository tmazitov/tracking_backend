package router

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine      *gin.Engine
	serviceName string
}

func NewRouter(serviceName string) *Router {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	r.Use(cors.New(config))

	return &Router{engine: r, serviceName: serviceName}
}

func (r *Router) Setup(endpoints []Endpoint) {

	var prefix string

	for _, endpoint := range endpoints {

		if endpoint.WS {
			prefix = fmt.Sprintf("/%s/ws", r.serviceName)
		} else {
			prefix = fmt.Sprintf("/%s/api", r.serviceName)
		}

		if endpoint.Middleware != nil {
			r.engine.Handle(endpoint.Method, prefix+endpoint.Path, endpoint.Middleware.Handle(), endpoint.Handler.Handle)
		} else {
			r.engine.Handle(endpoint.Method, prefix+endpoint.Path, endpoint.Handler.Handle)
		}
	}
}

func (r *Router) Run(port string) {
	r.engine.Run(":" + port)
}
