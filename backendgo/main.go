package main

import (
	"OnlineLibraryPortal/database"
	"OnlineLibraryPortal/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	database.Connect()

	r := gin.Default()
	r.Static("/uploads", "./uploads")

	routes.BookRoutes(r)

	r.Run(":8080")
}
