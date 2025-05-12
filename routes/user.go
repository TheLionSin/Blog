package routes

import (
	"Blog/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	r.GET("/users", handlers.GetUsers)
	r.GET("/user/:id", handlers.GetUser)
	r.POST("/user", handlers.CreateUser)
	r.PUT("/user/:id", handlers.UpdateUser)
	r.DELETE("/user/:id", handlers.DeleteUser)
}
