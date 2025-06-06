package routes

import (
	"OnlineLibraryPortal/controllers"
	"OnlineLibraryPortal/middleware"

	"github.com/gin-gonic/gin"
)

func BorrowRoutes(router *gin.Engine) {
	borrow := router.Group("/borrow")
	borrow.Use(middleware.JWTAuthMiddleware())
	{
		borrow.POST("/", controllers.BorrowBook)
		borrow.POST("/return", controllers.ReturnBook)
		borrow.GET("/history", controllers.BorrowingHistory)
	}
}
