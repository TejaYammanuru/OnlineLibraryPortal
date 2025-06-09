package controllers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"OnlineLibraryPortal/database"
	"OnlineLibraryPortal/models"

	"github.com/gin-gonic/gin"
)

func isAuthorizedToModify(role int) bool {
	return role == 1 || role == 2
}

type BookRequest struct {
	Title           string `json:"title" binding:"required"`
	Author          string `json:"author" binding:"required"`
	Genre           string `json:"genre" binding:"required"`
	PublicationDate string `json:"publication_date"`
	TotalCopies     int    `json:"total_copies"`
	CopiesAvailable int    `json:"copies_available"`
	ImageURL        string `json:"image_url"`
}

func CreateBook(c *gin.Context) {
	userRole := c.GetInt("userRole")
	fmt.Println("User role from context:", userRole)

	if !isAuthorizedToModify(userRole) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var req BookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	if req.TotalCopies < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Total copies cannot be negative"})
		return
	}
	if req.CopiesAvailable < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Available copies cannot be negative"})
		return
	}
	if req.CopiesAvailable > req.TotalCopies {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Available copies cannot exceed total copies"})
		return
	}

	var pubDate time.Time
	var err error
	if req.PublicationDate != "" {

		dateFormats := []string{
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05",
			"2006-01-02T15:04Z",
			"2006-01-02",
		}

		for _, format := range dateFormats {
			pubDate, err = time.Parse(format, req.PublicationDate)
			if err == nil {
				break
			}
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publication_date format"})
			return
		}
	} else {
		pubDate = time.Now()
	}

	var imageUrl string
	if req.ImageURL != "" {
		if strings.HasPrefix(req.ImageURL, "data:image/") {

			imageUrl, err = saveBase64Image(req.ImageURL)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image: " + err.Error()})
				return
			}
		} else {

			imageUrl = req.ImageURL
		}
	}

	book := models.Book{
		Title:           req.Title,
		Author:          req.Author,
		PublicationDate: pubDate,
		Genre:           req.Genre,
		TotalCopies:     req.TotalCopies,
		CopiesAvailable: req.CopiesAvailable,
		ImageURL:        imageUrl,
	}

	if err := database.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, book)
}

func UpdateBook(c *gin.Context) {
	userRole := c.GetInt("userRole")
	if !isAuthorizedToModify(userRole) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	bookID := c.Param("id")

	var req BookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// Validate data
	if req.TotalCopies < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Total copies cannot be negative"})
		return
	}
	if req.CopiesAvailable < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Available copies cannot be negative"})
		return
	}
	if req.CopiesAvailable > req.TotalCopies {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Available copies cannot exceed total copies"})
		return
	}

	// Parse publication date
	var pubDate time.Time
	var err error
	if req.PublicationDate != "" {
		dateFormats := []string{
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05",
			"2006-01-02T15:04Z",
			"2006-01-02",
		}

		for _, format := range dateFormats {
			pubDate, err = time.Parse(format, req.PublicationDate)
			if err == nil {
				break
			}
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publication_date format"})
			return
		}
	}

	// Find existing book
	var book models.Book
	if err := database.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Handle image upload
	var imageUrl string = book.ImageURL // Keep existing image by default
	if req.ImageURL != "" {
		if strings.HasPrefix(req.ImageURL, "data:image/") {
			// Handle base64 image
			imageUrl, err = saveBase64Image(req.ImageURL)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image: " + err.Error()})
				return
			}
		} else if req.ImageURL != book.ImageURL {
			// Handle new URL
			imageUrl = req.ImageURL
		}
	}

	// Update book
	book.Title = req.Title
	book.Author = req.Author
	book.Genre = req.Genre
	book.PublicationDate = pubDate
	book.TotalCopies = req.TotalCopies
	book.CopiesAvailable = req.CopiesAvailable
	book.ImageURL = imageUrl

	if err := database.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func saveBase64Image(base64String string) (string, error) {
	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		return "", err
	}

	// Parse the base64 string
	parts := strings.Split(base64String, ",")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid base64 format")
	}

	// Get the file extension from the data URL
	var ext string
	if strings.Contains(parts[0], "jpeg") {
		ext = ".jpg"
	} else if strings.Contains(parts[0], "png") {
		ext = ".png"
	} else if strings.Contains(parts[0], "gif") {
		ext = ".gif"
	} else {
		ext = ".jpg" // default
	}

	// Decode base64 data
	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	// Generate unique filename
	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
	filepath := "uploads/" + filename

	// Write file
	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return "", err
	}

	return "/" + filepath, nil
}

func GetBooks(c *gin.Context) {

	var books []models.Book
	database.DB.Find(&books)
	c.JSON(http.StatusOK, books)
}

func GetBook(c *gin.Context) {

	id := c.Param("id")
	var book models.Book
	if err := database.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

// func UpdateBook(c *gin.Context) {
// 	userRole := c.GetInt("userRole")
// 	fmt.Println("User role from context:", userRole)

// 	if !isAuthorizedToModify(userRole) {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
// 		return
// 	}

// 	id := c.Param("id")
// 	var book models.Book

// 	if err := database.DB.First(&book, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
// 		return
// 	}

// 	var input struct {
// 		Title           string `json:"title"`
// 		Author          string `json:"author"`
// 		Genre           string `json:"genre"`
// 		PublicationDate string `json:"publication_date"`
// 		TotalCopies     int    `json:"total_copies"`
// 		CopiesAvailable int    `json:"copies_available"`
// 		ImageURL        string `json:"image_url"` // Only relevant if base64 or image URL is sent from frontend
// 	}

// 	if err := c.BindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
// 		return
// 	}

// 	// Apply updates
// 	if input.Title != "" {
// 		book.Title = input.Title
// 	}
// 	if input.Author != "" {
// 		book.Author = input.Author
// 	}
// 	if input.Genre != "" {
// 		book.Genre = input.Genre
// 	}
// 	if input.PublicationDate != "" {
// 		pubDate, err := time.Parse("2006-01-02T15:04", input.PublicationDate)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publication_date format. Use YYYY-MM-DDTHH:MM"})
// 			return
// 		}
// 		book.PublicationDate = pubDate
// 	}
// 	book.TotalCopies = input.TotalCopies
// 	book.CopiesAvailable = input.CopiesAvailable

// 	// ImageURL field is optional unless you're processing base64 data
// 	if input.ImageURL != "" {
// 		book.ImageURL = input.ImageURL
// 	}

// 	if err := database.DB.Save(&book).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, book)
// }

func DeleteBook(c *gin.Context) {
	userRole := c.GetInt("userRole")
	fmt.Println("User role from context:", userRole)

	if !isAuthorizedToModify(userRole) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	id := c.Param("id")
	var book models.Book

	if err := database.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	database.DB.Delete(&book)
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}
