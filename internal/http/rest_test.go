package http

import (
	"concert/internal/concert"
	"concert/internal/models"
	"concert/test/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandler_GetFan(t *testing.T) {
	mockService := mocks.MockConcertService{
		GetFanFunc: func(name string) ([]models.Booking, error) {
			return []models.Booking{
				{Name: "Abdou", Price: 50},
				{Name: "khady", Price: 100},
				{Name: "Aminata", Price: 150},
			}, nil
		},
	}

	handler, err := NewRouter(&mockService, nil)
	assert.Nil(t, err)
	handler.ChiSetRoutes()

	request := httptest.NewRequest("GET", "/fan/Abdou", nil)
	response := httptest.NewRecorder()
	handler.Route.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
}

func TestHandler_ListAllShow(t *testing.T) {
	mockService := mocks.MockConcertService{
		ListAllShowFunc: func() ([]models.Show, error) {
			return []models.Show{
				{Venue: "Paris", Date: time.Now()},
				{Venue: "Marseille", Date: time.Now()},
			}, nil
		},
	}
	handler, err := NewRouter(&mockService, nil)
	assert.Nil(t, err)
	handler.ChiSetRoutes()
	request := httptest.NewRequest("GET", "/shows", nil)
	response := httptest.NewRecorder()
	handler.Route.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.True(t, strings.Contains(response.Body.String(), "Paris"))
	assert.True(t, strings.Contains(response.Body.String(), "Marseille"))
	assert.Equal(t, "text/html; charset=UTF-8", response.Header().Get("Content-Type"))
}

func TestHandler_ListAllFan(t *testing.T) {
	mockService := mocks.MockConcertService{
		ListAllFanFunc: func() ([]models.Booking, error) {
			return []models.Booking{
				{Name: "Abdou", Price: 50},
				{Name: "khady", Price: 100},
				{Name: "Aminata", Price: 150},
			}, nil
		},
	}
	handler, err := NewRouter(&mockService, nil)
	assert.Nil(t, err)
	handler.ChiSetRoutes()
	request := httptest.NewRequest("GET", "/list", nil)
	response := httptest.NewRecorder()
	handler.Route.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.True(t, strings.Contains(response.Body.String(), "Abdou"))
	assert.True(t, strings.Contains(response.Body.String(), "khady"))
	assert.True(t, strings.Contains(response.Body.String(), "Aminata"))
	assert.Equal(t, "text/html; charset=UTF-8", response.Header().Get("Content-Type"))
}

func TestHandler_SetArtist(t *testing.T) {
	mockService := mocks.MockConcertService{
		SetArtistFunc: func(artist concert.Artist) (concert.Artist, error) {
			return models.Artist{Name: "Drake", Genre: "Rock", ImageURL: "https://example.com/photo.jpg", AlbumURL: "https://example.com/album.jpg"}, nil
		},
	}
	handler, err := NewRouter(&mockService, nil)
	assert.Nil(t, err)
	handler.ChiSetRoutes()
	request := httptest.NewRequest("POST", "/artist", nil)
	response := httptest.NewRecorder()
	handler.Route.ServeHTTP(response, request)
	// statusSeeOther because we are redirecting to the new artist page
	// fake test we are not testing the multipart form data
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

// Test other functions
//........
//........
// Test other functions
