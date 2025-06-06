package routes

import (
	"OnlineLibraryPortal/controllers"

	"github.com/gin-gonic/gin"
)

func BookRoutes(router *gin.Engine) {
	books := router.Group("/books")
	{
		books.GET("/", controllers.GetBooks)
		books.GET("/:id", controllers.GetBook)
		books.POST("/", controllers.CreateBook)
		books.PUT("/:id", controllers.UpdateBook)
		books.DELETE("/:id", controllers.DeleteBook)
	}
}
