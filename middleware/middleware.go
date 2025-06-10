package middleware

import (
	"Blog/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.RespondError(c, http.StatusUnauthorized, "Требуется access токен")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := utils.ParseAccessToken(tokenStr)
		if err != nil {
			utils.RespondError(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func CanEditUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		authUserID := c.GetUint("user_id")

		idParam := c.Param("id")
		targetID, err := strconv.Atoi(idParam)
		if err != nil {
			utils.RespondError(c, http.StatusBadRequest, "Некорректный ID")
			c.Abort()
			return
		}

		if uint(targetID) != authUserID {
			utils.RespondError(c, http.StatusForbidden, "Нет прав на редактирование")
			c.Abort()
			return
		}
		c.Next()
	}
}
