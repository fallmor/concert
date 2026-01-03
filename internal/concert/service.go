package concert

import "concert/internal/models"

type ConcertService interface {
	GetFan(name string) ([]models.Booking, error)
	GetShow(artistName string) ([]models.Show, error)
	GetShowByID(id uint) (models.Show, error)
	SetShow(show models.Show) (models.Show, error)
	SetArtist(artist models.Artist) (models.Artist, error)
	ParticipateShow(fan models.Booking) (models.Show, error)
	ListAllShow() ([]models.Show, error)
	ListAllFan() ([]models.Booking, error)
	ListAllArtists() ([]models.Artist, error)
	GetArtistByID(id uint) (ArtistShow, error)
	GetAllUsers() ([]models.User, error)
	DeleteShow(id uint) error
	DeleteArtist(id uint) error
	DeleteFan(id uint) error
	UpdateShow(show models.Show) (models.Show, error)
	UpdateArtist(artist models.Artist) (models.Artist, error)
	UpdateFan(fan models.Booking) (models.Booking, error)
	GetMyBookings(user models.User) ([]models.Booking, error)
	GetBookingById(id uint) (models.Booking, error)
}
