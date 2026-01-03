package http

import (
	"concert/internal/concert"
	"concert/internal/models"

	"github.com/go-chi/chi/v5"
	"go.temporal.io/sdk/client"
	"gorm.io/gorm"
)

type Handler struct {
	Route          chi.Router
	Service        concert.ConcertService
	Db             *gorm.DB
	TemporalClient client.Client
}

type PageData struct {
	Title string
	Fans  []models.Booking
}

type ShowInfo struct {
	Shows []models.Show
}

type FanInfo struct {
	Fans []models.Booking
}

type ArtistInfo struct {
	Artists []models.Artist
}

type UserInfo struct {
	Users []models.User
}

type TemporalUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
