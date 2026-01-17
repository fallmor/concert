package http

import (
	"concert/internal/models"
	"concert/test/mocks"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.User{}, &models.Artist{}, &models.Show{}, &models.Booking{})
	return db
}

func TestHandler_GetFan(t *testing.T) {
	mockService := mocks.MockConcertService{
		GetFanFunc: func(name string) ([]models.Booking, error) {
			return []models.Booking{
				{UserID: 1, ShowID: 1, TicketCount: 1, TotalPrice: 50.0, Status: "confirmed"},
				{UserID: 2, ShowID: 2, TicketCount: 1, TotalPrice: 100.0, Status: "confirmed"},
				{UserID: 3, ShowID: 3, TicketCount: 1, TotalPrice: 150.0, Status: "confirmed"},
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
	request := httptest.NewRequest("GET", "/api/public/shows", nil)
	response := httptest.NewRecorder()
	handler.Route.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.True(t, strings.Contains(response.Body.String(), "Paris"))
	assert.True(t, strings.Contains(response.Body.String(), "Marseille"))
}

func TestHandler_GetMyBookings(t *testing.T) {
	user := models.User{Email: "abdou@example.com", Username: "abdou", PasswordHash: "password", Role: "user"}
	db := SetupTestDB()
	db.Create(&user)
	mockService := mocks.MockConcertService{
		GetMyBookingsFunc: func(user models.User) ([]models.Booking, error) {
			return []models.Booking{
				{UserID: 1, ShowID: 1, TicketCount: 1, TotalPrice: 50.0, Status: "confirmed"},
				{UserID: 2, ShowID: 2, TicketCount: 1, TotalPrice: 100.0, Status: "confirmed"},
				{UserID: 3, ShowID: 3, TicketCount: 1, TotalPrice: 150.0, Status: "confirmed"},
			}, nil
		},
	}
	handler, err := NewRouter(&mockService, db)
	assert.Nil(t, err)
	handler.ChiSetRoutes()
	request := httptest.NewRequest("GET", "/bookings", nil)
	response := httptest.NewRecorder()
	handler.Route.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
}

func TestHandler_SetArtist(t *testing.T) {
	admin := models.User{Email: "admin@test.com", Username: "admin", PasswordHash: "x", Role: "admin"}
	db := SetupTestDB()
	db.Create(&admin)

	mockService := mocks.MockConcertService{
		SetArtistFunc: func(artist models.Artist) (models.Artist, error) {
			return models.Artist{Name: "Drake", Genre: "Rock", ImageURL: "https://example.com/photo.jpg", AlbumURL: "https://example.com/album.jpg"}, nil
		},
	}
	handler, err := NewRouter(&mockService, db)
	assert.Nil(t, err)
	handler.ChiSetRoutes()

	request := httptest.NewRequest("POST", "/api/admin/artists", nil)
	request.AddCookie(&http.Cookie{Name: "session_concert", Value: fmt.Sprintf("user_%d", admin.ID)})
	response := httptest.NewRecorder()
	handler.Route.ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestHandler_CreateShow(t *testing.T) {
	admin := models.User{Email: "admin@test.com", Username: "admin", PasswordHash: "x", Role: "admin"}
	db := SetupTestDB()
	db.Create(&admin)

	mockService := mocks.MockConcertService{
		SetShowFunc: func(show models.Show) (models.Show, error) {
			return models.Show{Venue: "Paris", Date: time.Now()}, nil
		},
	}
	handler, err := NewRouter(&mockService, db)
	assert.Nil(t, err)
	handler.ChiSetRoutes()

	request := httptest.NewRequest("POST", "/api/admin/shows", nil)
	request.AddCookie(&http.Cookie{Name: "session_concert", Value: fmt.Sprintf("user_%d", admin.ID)})
	response := httptest.NewRecorder()
	handler.Route.ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestHandler_UpdateShow(t *testing.T) {
	admin := models.User{Email: "admin@test.com", Username: "admin", PasswordHash: "x", Role: "admin"}
	db := SetupTestDB()
	db.Create(&admin)

	mockService := mocks.MockConcertService{
		GetShowByIDFunc: func(id uint) (models.Show, error) {
			return models.Show{Venue: "Paris", Date: time.Now()}, nil
		},
		SetShowFunc: func(show models.Show) (models.Show, error) {
			return models.Show{Venue: "Paris", Date: time.Now()}, nil
		},
	}
	handler, err := NewRouter(&mockService, db)
	assert.Nil(t, err)
	handler.ChiSetRoutes()

	request := httptest.NewRequest("PUT", "/api/admin/shows/1", nil)
	request.AddCookie(&http.Cookie{Name: "session_concert", Value: fmt.Sprintf("user_%d", admin.ID)})
	response := httptest.NewRecorder()
	handler.Route.ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestHandler_DeleteShow(t *testing.T) {
	admin := models.User{Email: "admin@test.com", Username: "admin", PasswordHash: "x", Role: "admin"}
	db := SetupTestDB()
	db.Create(&admin)

	mockService := mocks.MockConcertService{
		GetShowByIDFunc: func(id uint) (models.Show, error) {
			return models.Show{Venue: "Paris", Date: time.Now()}, nil
		},
		CountConfirmSeatsFunc: func(id uint) int64 {
			return 0
		},
		DeleteShowFunc: func(id uint) error {
			return nil
		},
	}
	handler, err := NewRouter(&mockService, db)
	assert.Nil(t, err)
	handler.ChiSetRoutes()

	request := httptest.NewRequest("DELETE", "/api/admin/shows/1", nil)
	request.AddCookie(&http.Cookie{Name: "session_concert", Value: fmt.Sprintf("user_%d", admin.ID)})
	response := httptest.NewRecorder()
	handler.Route.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
}

func TestHandler_CreateArtist(t *testing.T) {
	admin := models.User{Email: "admin@test.com", Username: "admin", PasswordHash: "x", Role: "admin"}
	db := SetupTestDB()
	db.Create(&admin)

	mockService := mocks.MockConcertService{
		SetArtistFunc: func(artist models.Artist) (models.Artist, error) {
			return models.Artist{Name: "Drake", Genre: "Rock", ImageURL: "https://example.com/photo.jpg", AlbumURL: "https://example.com/album.jpg"}, nil
		},
	}
	handler, err := NewRouter(&mockService, db)
	assert.Nil(t, err)
	handler.ChiSetRoutes()

	request := httptest.NewRequest("POST", "/api/admin/artists", nil)
	request.AddCookie(&http.Cookie{Name: "session_concert", Value: fmt.Sprintf("user_%d", admin.ID)})
	response := httptest.NewRecorder()
	handler.Route.ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}
