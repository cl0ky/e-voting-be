package middleware

import (
	"github/com/cl0ky/e-voting-be/env"
	"strings"

	"github.com/gin-gonic/gin"
)

func getAllowOrigins() string {
	if len(env.AllowOrigins) > 0 {
		return strings.Join(env.AllowOrigins, ", ")
	}
	return "*"
}

func getAllowHeaders() string {
	var h []string = []string{
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"Authorization",
		"accept",
		"origin",
		"Cache-Control",
		"X-Token",
		"X-Auth-Token",
		"clientpath",
		// "X-CSRF-Token",
		// "X-Requested-With",
	}
	return strings.Join(h, ", ")
}

func corsMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", getAllowOrigins())
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", getAllowHeaders())
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}

func CORSMiddleware() gin.HandlerFunc {
	return corsMiddleware
}
