package concert

import (
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
	db.AutoMigrate(&Artist{}, &Show{}, &Fan{})
	return db
}
func TestGetFan(t *testing.T) {
	db := SetupTestDB()
	service := NewConcert(db)

	artist := Artist{Name: "Drake", Genre: "Rock"}
	db.Create(&artist)
	show := Show{ArtistID: artist.ID, Venue: "Paris", Date: time.Now()}
	db.Create(&show)

	fan := Fan{Name: "Abdou", ShowID: show.ID, Price: 100}
	db.Create(&fan)

	found, err := service.GetFan("Abdou")
	assert.NoError(t, err)
	assert.NotEmpty(t, found)
	assert.Equal(t, fan.Name, found[0].Name)
}
func TestGetShow(t *testing.T) {
	db := SetupTestDB()
	service := NewConcert(db)
	artist := Artist{Name: "Drake", Genre: "Rock"}
	db.Create(&artist)
	show := Show{ArtistID: artist.ID, Venue: "Paris", Date: time.Now()}
	db.Create(&show)
	found, err := service.GetShow("Drake")
	assert.NoError(t, err)
	assert.NotEmpty(t, found)
	assert.Equal(t, show.Venue, found[0].Venue)
}
func TestGetShowByID(t *testing.T) {
	db := SetupTestDB()
	service := NewConcert(db)
	artist := Artist{Name: "Drake", Genre: "Rock"}
	db.Create(&artist)
	show := Show{ArtistID: artist.ID, Venue: "Paris", Date: time.Now()}
	db.Create(&show)
	found, err := service.GetShowByID(show.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, found)
	assert.Equal(t, show.Venue, found.Venue)
}
func TestSetShow(t *testing.T) {
	db := SetupTestDB()
	service := NewConcert(db)
	artist := Artist{Name: "Drake", Genre: "Rock"}
	db.Create(&artist)
	show := Show{ArtistID: artist.ID, Venue: "Paris", Date: time.Now()}
	db.Create(&show)
	found, err := service.SetShow(show)
	assert.NoError(t, err)
	assert.NotEmpty(t, found)
	assert.Equal(t, show.Venue, found.Venue)
}
func TestSetArtist(t *testing.T) {
	db := SetupTestDB()
	service := NewConcert(db)
	artist := Artist{Name: "Drake", Genre: "Rock"}
	db.Create(&artist)
	found, err := service.SetArtist(artist)
	assert.NoError(t, err)
	assert.NotEmpty(t, found)
	assert.Equal(t, artist.Name, found.Name)
}
func TestParticipateShow(t *testing.T) {
	db := SetupTestDB()
	service := NewConcert(db)
	artist := Artist{Name: "Drake", Genre: "Rock"}
	db.Create(&artist)
	show := Show{ArtistID: artist.ID, Venue: "Paris", Date: time.Now()}
	db.Create(&show)
	fan := Fan{Name: "Abdou", ShowID: show.ID, Price: 100}
	db.Create(&fan)
	found, err := service.ParticipateShow(fan)
	assert.NoError(t, err)
	assert.NotEmpty(t, found)
	assert.Equal(t, fan.Name, found.Fans[0].Name)
}
