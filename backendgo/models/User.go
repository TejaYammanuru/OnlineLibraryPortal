package models

type User struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	JTI   string `gorm:"uniqueIndex" json:"jti"`
	Role  int    `json:"role"`
}
