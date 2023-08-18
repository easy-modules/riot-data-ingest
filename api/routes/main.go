package routes

import (
	"github.com/gin-gonic/gin"
)

func AppRoute(route *gin.Engine) *gin.RouterGroup {
	v1 := route.Group("/v1")

	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}
	return v1
}
