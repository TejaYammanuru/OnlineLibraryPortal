package controllers

import (
	"net/http"

	"time"

	"OnlineLibraryPortal/database"
	"OnlineLibraryPortal/models"

	"github.com/gin-gonic/gin"
)

// Create a new book
func CreateBook(c *gin.Context) {
	var input struct {
		Title           string `json:"title" binding:"required"`
		Author          string `json:"author" binding:"required"`
		PublicationDate string `json:"publication_date" binding:"required"`
		Genre           string `json:"genre" binding:"required"`
		TotalCopies     int    `json:"total_copies" binding:"required"`
		CopiesAvailable int    `json:"copies_available" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse publication date
	pubDate, err := time.Parse("2006-01-02", input.PublicationDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publication_date format. Use YYYY-MM-DD"})
		return
	}

	book := models.Book{
		Title:           input.Title,
		Author:          input.Author,
		PublicationDate: pubDate,
		Genre:           input.Genre,
		TotalCopies:     input.TotalCopies,
		CopiesAvailable: input.CopiesAvailable,
	}

	database.DB.Create(&book)

	c.JSON(http.StatusCreated, book)
}

// Get all books
func GetBooks(c *gin.Context) {
	var books []models.Book
	database.DB.Find(&books)
	c.JSON(http.StatusOK, books)
}

// Get one book by ID
func GetBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book

	if err := database.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// Update book by ID
func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book

	if err := database.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	var input struct {
		Title           *string `json:"title"`
		Author          *string `json:"author"`
		PublicationDate *string `json:"publication_date"` // expect date as string (e.g. "2023-01-01")
		Genre           *string `json:"genre"`
		TotalCopies     *int    `json:"total_copies"`
		CopiesAvailable *int    `json:"copies_available"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Title != nil {
		book.Title = *input.Title
	}
	if input.Author != nil {
		book.Author = *input.Author
	}
	if input.PublicationDate != nil {
		parsedDate, err := time.Parse("2006-01-02", *input.PublicationDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publication_date format. Use YYYY-MM-DD"})
			return
		}
		book.PublicationDate = parsedDate
	}
	if input.Genre != nil {
		book.Genre = *input.Genre
	}
	if input.TotalCopies != nil {
		book.TotalCopies = *input.TotalCopies
	}
	if input.CopiesAvailable != nil {
		book.CopiesAvailable = *input.CopiesAvailable
	}

	database.DB.Save(&book)

	c.JSON(http.StatusOK, book)
}

// Delete book by ID
func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book

	if err := database.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	database.DB.Delete(&book)

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
