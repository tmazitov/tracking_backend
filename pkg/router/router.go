package router

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine      *gin.Engine
	middleware  []gin.HandlerFunc
	serviceName string
}

func NewRouter(serviceName string) *Router {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	r.Use(cors.New(config))

	return &Router{engine: r, serviceName: serviceName, middleware: []gin.HandlerFunc{}}
}

func (r *Router) AddMiddleware(middleware []gin.HandlerFunc) {
	r.engine.Use(middleware...)
}

func (r *Router) Setup(endpoints []Endpoint) {

	var prefix string

	for _, endpoint := range endpoints {

		if endpoint.WS {
			prefix = fmt.Sprintf("/%s/ws", r.serviceName)
		} else {
			prefix = fmt.Sprintf("/%s/api", r.serviceName)
		}

		var endpointProcess []gin.HandlerFunc = endpoint.Middleware
		endpointProcess = append(endpointProcess, endpoint.Handler.Handle)

		r.engine.Handle(endpoint.Method, prefix+endpoint.Path, endpointProcess...)
	}
}

func (r *Router) Run(port string) {
	r.engine.Run(":" + port)
}
