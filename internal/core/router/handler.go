package router

import "github.com/gin-gonic/gin"

type Handler interface {
	Input() interface{}
	Result() interface{}
	Action(c *gin.Context)
}
