package mocks

import (
	"concert/internal/models"
	"mime/multipart"
)

type MockConcertService struct {
	GetFanFunc          func(name string) ([]models.Booking, error)
	GetShowFunc         func(artistName string) ([]models.Show, error)
	SetShowFunc         func(show models.Show) (models.Show, error)
	SetArtistFunc       func(artist models.Artist) (models.Artist, error)
	ParticipateShowFunc func(fan models.Booking) (models.Show, error)
	ListAllShowFunc     func() ([]models.Show, error)
	ListAllFanFunc      func() ([]models.Booking, error)
	ListAllArtistsFunc  func() ([]models.Artist, error)
	GetArtistByIDFunc   func(id uint) (models.Artist, error)
	GetShowByIDFunc     func(id uint) (models.Show, error)
	UploadImageFunc     func(file multipart.File, handler *multipart.FileHeader, folder string) (string, error)
	GetAllUsersFunc     func() ([]models.User, error)
	DeleteShowFunc      func(id uint) error
	DeleteArtistFunc    func(id uint) error
	DeleteFanFunc       func(id uint) error
	UpdateShowFunc      func(show models.Show) (models.Show, error)
	UpdateArtistFunc    func(artist models.Artist) (models.Artist, error)
	UpdateFanFunc       func(fan models.Booking) (models.Booking, error)
}

func (m *MockConcertService) GetFan(name string) ([]models.Booking, error) {
	return m.GetFanFunc(name)
}

func (m *MockConcertService) GetShow(artistName string) ([]models.Show, error) {
	return m.GetShowFunc(artistName)
}

func (m *MockConcertService) SetShow(show models.Show) (models.Show, error) {
	return m.SetShowFunc(show)
}

func (m *MockConcertService) SetArtist(artist models.Artist) (models.Artist, error) {
	return m.SetArtistFunc(artist)
}

func (m *MockConcertService) ParticipateShow(fan models.Booking) (models.Show, error) {
	return m.ParticipateShowFunc(fan)
}

func (m *MockConcertService) ListAllShow() ([]models.Show, error) {
	return m.ListAllShowFunc()
}

func (m *MockConcertService) ListAllFan() ([]models.Booking, error) {
	return m.ListAllFanFunc()
}

func (m *MockConcertService) ListAllArtists() ([]models.Artist, error) {
	return m.ListAllArtistsFunc()
}

func (m *MockConcertService) GetShowByID(id uint) (models.Show, error) {
	return m.GetShowByIDFunc(id)
}
func (m *MockConcertService) GetArtistByID(id uint) (models.Artist, error) {
	return m.GetArtistByIDFunc(id)
}

func (m *MockConcertService) GetAllUsers() ([]models.User, error) {
	return m.GetAllUsersFunc()
}

func (m *MockConcertService) DeleteShow(id uint) error {
	return m.DeleteShowFunc(id)
}

func (m *MockConcertService) DeleteArtist(id uint) error {
	return m.DeleteArtistFunc(id)
}

func (m *MockConcertService) DeleteFan(id uint) error {
	return m.DeleteFanFunc(id)
}

func (m *MockConcertService) UpdateShow(show models.Show) (models.Show, error) {
	return m.UpdateShowFunc(show)
}

func (m *MockConcertService) UpdateArtist(artist models.Artist) (models.Artist, error) {
	return m.UpdateArtistFunc(artist)
}

func (m *MockConcertService) UpdateFan(fan models.Booking) (models.Booking, error) {
	return m.UpdateFanFunc(fan)
}
