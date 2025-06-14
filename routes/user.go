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
	protected.GET("/me", handlers.GetCurrentUser)

	protected.PUT("/user/:id", middleware.CanEditOrAdmin(), handlers.UpdateUser)
	protected.DELETE("/user/:id", middleware.CanEditOrAdmin(), handlers.DeleteUser)
}
