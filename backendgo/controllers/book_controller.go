package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"OnlineLibraryPortal/database"
	"OnlineLibraryPortal/models"

	"github.com/gin-gonic/gin"
)

func CreateBook(c *gin.Context) {

	title := c.PostForm("title")
	author := c.PostForm("author")
	publicationDate := c.PostForm("publication_date")
	genre := c.PostForm("genre")
	totalCopiesStr := c.PostForm("total_copies")
	copiesAvailableStr := c.PostForm("copies_available")

	totalCopies, err := strconv.Atoi(totalCopiesStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid total_copies"})
		return
	}
	copiesAvailable, err := strconv.Atoi(copiesAvailableStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid copies_available"})
		return
	}

	pubDate, err := time.Parse("2006-01-02", publicationDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publication_date format. Use YYYY-MM-DD"})
		return
	}

	file, err := c.FormFile("image")
	var imageUrl string
	if err == nil {

		if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create uploads folder"})
			return
		}

		filename := filepath.Base(file.Filename)
		dst := "uploads/" + strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + filename

		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		imageUrl = "/" + dst
	} else {
		imageUrl = ""
	}

	book := models.Book{
		Title:           title,
		Author:          author,
		PublicationDate: pubDate,
		Genre:           genre,
		TotalCopies:     totalCopies,
		CopiesAvailable: copiesAvailable,
		ImageURL:        imageUrl,
	}

	database.DB.Create(&book)

	c.JSON(http.StatusCreated, book)
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

func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book

	if err := database.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	title := c.PostForm("title")
	author := c.PostForm("author")
	publicationDate := c.PostForm("publication_date")
	genre := c.PostForm("genre")
	totalCopiesStr := c.PostForm("total_copies")
	copiesAvailableStr := c.PostForm("copies_available")

	if title != "" {
		book.Title = title
	}
	if author != "" {
		book.Author = author
	}
	if publicationDate != "" {
		pubDate, err := time.Parse("2006-01-02", publicationDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publication_date format. Use YYYY-MM-DD"})
			return
		}
		book.PublicationDate = pubDate
	}
	if genre != "" {
		book.Genre = genre
	}
	if totalCopiesStr != "" {
		totalCopies, err := strconv.Atoi(totalCopiesStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid total_copies"})
			return
		}
		book.TotalCopies = totalCopies
	}
	if copiesAvailableStr != "" {
		copiesAvailable, err := strconv.Atoi(copiesAvailableStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid copies_available"})
			return
		}
		book.CopiesAvailable = copiesAvailable
	}

	file, err := c.FormFile("image")
	if err == nil {
		if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create uploads folder"})
			return
		}
		filename := filepath.Base(file.Filename)
		dst := "uploads/" + strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + filename

		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		book.ImageURL = "/" + dst
	}

	database.DB.Save(&book)
	c.JSON(http.StatusOK, book)
}

func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book

	if err := database.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	database.DB.Delete(&book)
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}
