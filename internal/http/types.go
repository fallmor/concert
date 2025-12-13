package http

import (
	"concert/internal/concert"

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
	Fans  []concert.Fan
}

type ShowInfo struct {
	Shows []concert.Show
}

type FanInfo struct {
	Fans []concert.Fan
}

type ArtistInfo struct {
	Artists []concert.Artist
}

type UserInfo struct {
	Users []concert.User
}
