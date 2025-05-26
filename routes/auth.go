package routes

import (
	"Blog/handlers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/login", handlers.Login)
	r.POST("/register", handlers.Register)
	r.POST("/refresh", handlers.RefreshToken)
}
