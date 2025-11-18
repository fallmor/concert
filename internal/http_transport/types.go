package httptransport

import (
	"concert/internal/concert"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Handler struct {
	Route   chi.Router
	Service concert.ConcertService
	Db      *gorm.DB
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