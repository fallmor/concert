package concert

import (
	"concert/internal/models"

	"gorm.io/gorm"
)

// to be deleted after
type User struct {
	gorm.Model
	Email        string           `gorm:"uniqueIndex;not null" json:"email"`
	Username     string           `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string           `gorm:"not null" json:"-"`
	FirstName    string           `json:"firstName,omitempty"`
	LastName     string           `json:"lastName,omitempty"`
	Role         string           `gorm:"default:user"` // user, admin, moderator
	Fans         []models.Booking `gorm:"foreignKey:UserID"`
}

type UserRole string

const (
	RoleUser      UserRole = "user"
	RoleModerator UserRole = "moderator"
	RoleAdmin     UserRole = "admin"
)
