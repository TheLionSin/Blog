package main

import (
	"Blog/models"
	"Blog/routes"
	"Blog/storage"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	storage.ConnectDB()
	storage.DB.AutoMigrate(&models.User{})

	routes.RegisterUserRoutes(r)

	r.Run(":8080")

}
