package concert

import (
	"concert/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Artist{}, &models.Show{}, &models.Booking{})
	return db
}
func TestGetFan(t *testing.T) {
	db := SetupTestDB()
	service := NewConcert(db)

	artist := models.Artist{Name: "Drake", Genre: "Rock"}
	db.Create(&artist)

	show := models.Show{
		ArtistID:       artist.ID,
		Venue:          "Paris",
		Date:           time.Now(),
		Title:          "Drake Concert",
		Price:          50.0,
		TotalSeats:     100,
		AvailableSeats: 100,
	}
	db.Create(&show)

	user := models.User{
		Email:     "toto@gmail.com",
		Username:  "toto",
		FirstName: "Abdou",
		Role:      "user",
	}
	db.Create(&user)

	booking := models.Booking{
		UserID:      user.ID,
		ShowID:      show.ID,
		TicketCount: 2,
		TotalPrice:  100.0,
	}
	db.Create(&booking)

	found, err := service.GetFan("toto")

	assert.NoError(t, err)
	assert.NotEmpty(t, found)
	assert.Equal(t, 1, len(found))
	assert.Equal(t, booking.ID, found[0].ID)
	assert.Equal(t, "toto", found[0].User.Username)
}
func TestGetShow(t *testing.T) {
	db := SetupTestDB()
	service := NewConcert(db)
	artist := models.Artist{Name: "Drake", Genre: "Rock"}
	db.Create(&artist)
	show := models.Show{ArtistID: artist.ID, Venue: "Paris", Date: time.Now()}
	db.Create(&show)
	found, err := service.GetShow("Drake")
	assert.NoError(t, err)
	assert.NotEmpty(t, found)
	assert.Equal(t, show.Venue, found[0].Venue)
}
func TestGetShowByID(t *testing.T) {
	db := SetupTestDB()
	service := NewConcert(db)
	artist := models.Artist{Name: "Drake", Genre: "Rock"}
	db.Create(&artist)
	show := models.Show{ArtistID: artist.ID, Venue: "Paris", Date: time.Now()}
	db.Create(&show)
	found, err := service.GetShowByID(show.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, found)
	assert.Equal(t, show.Venue, found.Venue)
}
func TestSetShow(t *testing.T) {
	db := SetupTestDB()
	service := NewConcert(db)
	artist := models.Artist{Name: "Drake", Genre: "Rock"}
	db.Create(&artist)
	show := models.Show{ArtistID: artist.ID, Venue: "Paris", Date: time.Now()}
	db.Create(&show)
	found, err := service.SetShow(show)
	assert.NoError(t, err)
	assert.NotEmpty(t, found)
	assert.Equal(t, show.Venue, found.Venue)
}
func TestSetArtist(t *testing.T) {
	db := SetupTestDB()
	service := NewConcert(db)
	artist := models.Artist{Name: "Drake", Genre: "Rock"}
	db.Create(&artist)
	found, err := service.SetArtist(artist)
	assert.NoError(t, err)
	assert.NotEmpty(t, found)
	assert.Equal(t, artist.Name, found.Name)
}
func TestParticipateShow(t *testing.T) {
	db := SetupTestDB()
	service := NewConcert(db)
	artist := models.Artist{Name: "Youssou", Genre: "Rock"}
	db.Create(&artist)

	show := models.Show{
		ArtistID:       artist.ID,
		Venue:          "Paris",
		Date:           time.Now(),
		Title:          "Drake Concert",
		Price:          50.0,
		TotalSeats:     100,
		AvailableSeats: 100,
	}
	db.Create(&show)

	user := models.User{
		Email:     "toto1@gmail.com",
		Username:  "toto1",
		FirstName: "Abdou",
		Role:      "user",
	}
	db.Create(&user)

	booking := models.Booking{
		UserID:      user.ID,
		ShowID:      show.ID,
		TicketCount: 2,
		TotalPrice:  100.0,
	}
	db.Create(&booking)

	found, err := service.SetShow(show)
	assert.NoError(t, err)
	assert.NotEmpty(t, found)
	assert.Equal(t, show.Venue, found.Venue)
}
