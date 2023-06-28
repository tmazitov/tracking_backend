package router

import "github.com/gin-gonic/gin"

type Endpoint struct {
	Method     string
	WS         bool
	Path       string
	Middleware []gin.HandlerFunc
	Handler    Handler
}
