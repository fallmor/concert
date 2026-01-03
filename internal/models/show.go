package models

import (
	"time"

	"gorm.io/gorm"
)

type Show struct {
	gorm.Model
	Title          string    `gorm:"not null" json:"title"`
	Date           time.Time `gorm:"not null" json:"date"`
	Time           string    `json:"time"`
	ArtistID       uint      `gorm:"not null" json:"artistId"`
	Artist         Artist    `gorm:"foreignKey:ArtistID;references:ID" json:"artist"`
	Venue          string    `gorm:"not null" json:"venue"`
	Price          float64   `gorm:"not null" json:"price"`
	TotalSeats     int       `gorm:"not null" json:"totalSeats"`
	AvailableSeats int       `gorm:"not null" json:"availableSeats"`
	Description    string    `json:"description,omitempty"`
	ImageURL       string    `json:"imageUrl,omitempty"`
	Bookings       []Booking `gorm:"foreignKey:ShowID" json:"-"`
}
