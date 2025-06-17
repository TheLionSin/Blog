package routes

import (
	"Blog/handlers"
	"Blog/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth())

	protected.GET("/user/:id", handlers.GetUser)
	protected.POST("/user", handlers.CreateUser)
	protected.GET("/me", handlers.GetCurrentUser)
	protected.POST("user/avatar", handlers.UploadAvatar)
	protected.PUT("/user/:id", middleware.CanEditOrAdmin(), handlers.UpdateUser)
	protected.DELETE("/user/:id", middleware.CanEditOrAdmin(), handlers.DeleteUser)

	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.RequireAuth(), middleware.RequireAdmin())
	adminRoutes.GET("/users", handlers.GetUsers)
	adminRoutes.PUT("/user/:id/restore", handlers.RestoreUser)
	adminRoutes.GET("/users/export",handlers.ExportUsersCSV)
	adminRoutes.GET("/audit-logs", handlers.GetAuditLogs)

}
