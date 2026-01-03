package models

import "gorm.io/gorm"

type Artist struct {
	gorm.Model
	Name     string `gorm:"not null" json:"name"`
	Genre    string `json:"genre"`
	Bio      string `json:"bio,omitempty"`
	ImageURL string `json:"imageUrl,omitempty"`
	AlbumURL string `json:"albumUrl,omitempty"`
	Shows    []Show `gorm:"foreignKey:ArtistID" json:"-"`
}
