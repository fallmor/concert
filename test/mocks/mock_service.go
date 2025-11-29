package mocks

import (
	"concert/internal/concert"
	"mime/multipart"
)

type MockConcertService struct {
	GetFanFunc          func(name string) ([]concert.Fan, error)
	GetShowFunc         func(artistName string) ([]concert.Show, error)
	SetShowFunc         func(show concert.Show) (concert.Show, error)
	SetArtistFunc       func(artist concert.Artist) (concert.Artist, error)
	ParticipateShowFunc func(fan concert.Fan) (concert.Show, error)
	ListAllShowFunc     func() ([]concert.Show, error)
	ListAllFanFunc      func() ([]concert.Fan, error)
	ListAllArtistsFunc  func() ([]concert.Artist, error)
	GetShowByIDFunc     func(id uint) (concert.Show, error)
	UploadImageFunc     func(file multipart.File, handler *multipart.FileHeader, folder string) (string, error)
	GetAllUsersFunc     func() ([]concert.User, error)
	DeleteShowFunc      func(id uint) error
	DeleteArtistFunc    func(id uint) error
	DeleteFanFunc       func(id uint) error
	UpdateShowFunc      func(show concert.Show) (concert.Show, error)
	UpdateArtistFunc    func(artist concert.Artist) (concert.Artist, error)
	UpdateFanFunc       func(fan concert.Fan) (concert.Fan, error)
}

func (m *MockConcertService) GetFan(name string) ([]concert.Fan, error) {
	return m.GetFanFunc(name)
}

func (m *MockConcertService) GetShow(artistName string) ([]concert.Show, error) {
	return m.GetShowFunc(artistName)
}

func (m *MockConcertService) SetShow(show concert.Show) (concert.Show, error) {
	return m.SetShowFunc(show)
}

func (m *MockConcertService) SetArtist(artist concert.Artist) (concert.Artist, error) {
	return m.SetArtistFunc(artist)
}

func (m *MockConcertService) ParticipateShow(fan concert.Fan) (concert.Show, error) {
	return m.ParticipateShowFunc(fan)
}

func (m *MockConcertService) ListAllShow() ([]concert.Show, error) {
	return m.ListAllShowFunc()
}

func (m *MockConcertService) ListAllFan() ([]concert.Fan, error) {
	return m.ListAllFanFunc()
}

func (m *MockConcertService) ListAllArtists() ([]concert.Artist, error) {
	return m.ListAllArtistsFunc()
}

func (m *MockConcertService) GetShowByID(id uint) (concert.Show, error) {
	return m.GetShowByIDFunc(id)
}

func (m *MockConcertService) GetAllUsers() ([]concert.User, error) {
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

func (m *MockConcertService) UpdateShow(show concert.Show) (concert.Show, error) {
	return m.UpdateShowFunc(show)
}

func (m *MockConcertService) UpdateArtist(artist concert.Artist) (concert.Artist, error) {
	return m.UpdateArtistFunc(artist)
}

func (m *MockConcertService) UpdateFan(fan concert.Fan) (concert.Fan, error) {
	return m.UpdateFanFunc(fan)
}
