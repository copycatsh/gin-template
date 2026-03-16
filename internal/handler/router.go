package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *UserHandler) *gin.Engine {
	r := gin.Default()
	_ = r.SetTrustedProxies(nil)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Gin API!"})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/users", userHandler.GetAll)
	r.POST("/users", userHandler.Create)
	r.PUT("/users/:id", userHandler.Update)

	return r
}
