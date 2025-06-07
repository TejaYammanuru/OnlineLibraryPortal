package controllers

import (
	"OnlineLibraryPortal/database"
	"OnlineLibraryPortal/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func BorrowBook(c *gin.Context) {
	userRole := c.MustGet("userRole").(int)
	if userRole != 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only members can borrow books"})
		return
	}

	var req struct {
		BookID uint `json:"book_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID := c.MustGet("userID").(uint)

	var book models.Book
	if err := database.DB.First(&book, req.BookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if book.CopiesAvailable < 1 {
		c.JSON(http.StatusConflict, gin.H{"error": "No copies available"})
		return
	}

	borrowRecord := models.BorrowRecord{
		UserID:     userID,
		BookID:     req.BookID,
		BorrowedAt: time.Now(),
	}

	tx := database.DB.Begin()

	if err := tx.Create(&borrowRecord).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create borrow record"})
		return
	}

	book.CopiesAvailable -= 1
	if err := tx.Save(&book).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book availability"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Book borrowed successfully"})
}

func ReturnBook(c *gin.Context) {
	userRole := c.MustGet("userRole").(int)
	if userRole != 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only members can return books"})
		return
	}

	var req struct {
		BookID uint `json:"book_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID := c.MustGet("userID").(uint)

	var borrowRecord models.BorrowRecord
	err := database.DB.
		Where("user_id = ? AND book_id = ? AND returned_at IS NULL", userID, req.BookID).
		Order("borrowed_at desc").
		First(&borrowRecord).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active borrow record found"})
		return
	}

	now := time.Now()
	borrowRecord.ReturnedAt = &now

	tx := database.DB.Begin()

	if err := tx.Save(&borrowRecord).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update borrow record"})
		return
	}

	if err := tx.Model(&models.Book{}).Where("id = ?", req.BookID).Update("copies_available", gorm.Expr("copies_available + ?", 1)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update book availability"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Book returned successfully"})
}

func BorrowingHistory(c *gin.Context) {
	fmt.Println("BorrowingHistory handler called")

	userID := c.MustGet("userID").(uint)
	userRole := c.MustGet("userRole").(int)

	fmt.Println("UserID:", userID)
	fmt.Println("UserRole:", userRole)

	var records []models.BorrowRecord
	var err error

	if userRole == 1 {
		fmt.Println("Role is librarian, fetching all records")

		err = database.DB.Preload("Book").
			Preload("User", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name", "email", "jti", "role")
			}).
			Find(&records).Error
	} else {
		fmt.Println("Role is member, fetching their own history")
		err = database.DB.Preload("Book").
			Preload("User", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name", "email", "jti", "role")
			}).
			Where("user_id = ?", userID).
			Find(&records).Error
	}

	if err != nil {
		fmt.Println("Error fetching borrowing history:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch borrowing history"})
		return
	}

	fmt.Println("Fetched borrow records count:", len(records))
	c.JSON(http.StatusOK, records)
}
