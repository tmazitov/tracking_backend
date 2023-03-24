package core

import (
	"log"

	"github.com/gin-gonic/gin"
)

func ErrorLog(status int, message string, err error, c *gin.Context) {
	log.Println(err)
	c.JSON(status, gin.H{
		"status": message,
	})
}

func SendResponse(status int, responseData interface{}, c *gin.Context) {
	c.JSON(status, responseData)
}
