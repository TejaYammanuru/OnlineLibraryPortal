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
	routes.BorrowRoutes(r)

	r.Run(":8080")
}
