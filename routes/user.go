package routes

import (
	"Blog/handlers"
	"Blog/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth())
	protected.GET("/users", handlers.GetUsers)
	protected.GET("/user/:id", handlers.GetUser)
	protected.POST("/user", handlers.CreateUser)
	protected.PUT("/user/:id", handlers.UpdateUser)
	protected.DELETE("/user/:id", handlers.DeleteUser)
	protected.GET("me", handlers.GetCurrentUser)
}
