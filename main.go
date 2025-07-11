package main

import (
	"Blog/migrate"
	"Blog/routes"
	"Blog/storage"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.Static("/uploads", "./uploads")
	storage.ConnectDB()
	migrate.RunMigrations()

	routes.RegisterUserRoutes(r)
	routes.AuthRoutes(r)

	r.Run(":8080")

}
