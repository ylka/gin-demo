package handlers

import "github.com/gin-gonic/gin"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-API-KEY") != "xxxx" {
			c.AbortWithStatus(401)
		}

		c.Next()
	}
}
