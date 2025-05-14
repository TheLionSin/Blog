package middleware

import (
	"Blog/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.RespondError(c, http.StatusUnauthorized, "Требуется токен")
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.RespondError(c, http.StatusUnauthorized, "Неверный формат токена")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		
		userID, err := utils.ParseJWT(tokenStr)
		if err != nil {
			utils.RespondError(c, http.StatusUnauthorized, "Неверный токен")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
