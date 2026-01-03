package models

import "gorm.io/gorm"

type Booking struct {
	gorm.Model
	UserID      uint    `gorm:"not null" json:"userId"`
	User        User    `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	ShowID      uint    `gorm:"not null" json:"showId"`
	Show        Show    `gorm:"foreignKey:ShowID;references:ID" json:"show,omitempty"`
	TicketCount int     `gorm:"not null" json:"ticketCount"`
	TotalPrice  float64 `gorm:"not null" json:"totalPrice"`
	Status      string  `gorm:"default:'confirmed'" json:"status"`
}
