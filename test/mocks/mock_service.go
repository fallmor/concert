package mocks

import (
	"concert/internal/concert"
	"mime/multipart"
)


type MockConcertService struct {
	GetFanFunc func(name string) ([]concert.Fan, error)
	GetShowFunc func(artistName string) ([]concert.Show, error)
	SetShowFunc func(show concert.Show) (concert.Show, error)
	SetArtistFunc func(artist concert.Artist) (concert.Artist, error)
	ParticipateShowFunc func(fan concert.Fan) (concert.Show, error)
	ListAllShowFunc func() ([]concert.Show, error)
	ListAllFanFunc func() ([]concert.Fan, error)
	ListAllArtistsFunc func() ([]concert.Artist, error)
	GetShowByIDFunc func(id uint) (concert.Show, error)
	UploadImageFunc func(file multipart.File, handler *multipart.FileHeader, folder string) (string, error)
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