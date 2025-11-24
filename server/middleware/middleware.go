package middleware

import (
	"net/http"
	"strings"

	env "github/com/cl0ky/e-voting-be/env"
	"github/com/cl0ky/e-voting-be/internal/auth"
	jwtutil "github/com/cl0ky/e-voting-be/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		tokenStr := ""
		if header != "" && strings.HasPrefix(header, "Bearer ") {
			tokenStr = strings.TrimPrefix(header, "Bearer ")
		} else {
			cookie, err := c.Cookie("token")
			if err == nil {
				tokenStr = cookie
			}
		}

		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		jwtManager := jwtutil.NewJWTManager(env.JWTSecret, 0)
		claims, err := jwtManager.Verify(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		dbVal, ok := c.MustGet("db").(*gorm.DB)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "DB not found in context"})
			return
		}
		userRepo := auth.NewRepository(dbVal)
		user, err := userRepo.GetUserByID(claims.UserId)
		if err != nil || user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Set("user_id", claims.UserId)
		c.Set("role", claims.Role)
		c.Set("user", user)
		c.Next()
	}
}
