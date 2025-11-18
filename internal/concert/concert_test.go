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
	
	artist := Artist{Nom: "Drake", Genre: "Rock"}
	db.Create(&artist)
	show := Show{ArtistID: artist.ID, Place: "Paris", Date: time.Now()}
	db.Create(&show)
	
	fan := Fan{Nom: "Abdou", ShowID: show.ID, Price: 100}
	db.Create(&fan)
	
	found, err := service.GetFan("Abdou")
	assert.NoError(t, err)
	assert.NotEmpty(t, found)
	assert.Equal(t, fan.Nom, found[0].Nom)
}
func TestGetShow(t *testing.T) {
	db := SetupTestDB()
	service := NewConcert(db)
	artist := Artist{Nom: "Drake", Genre: "Rock"}
	db.Create(&artist)
	show := Show{ArtistID: artist.ID, Place: "Paris", Date: time.Now()}
	db.Create(&show)
	found, err := service.GetShow("Drake")
	assert.NoError(t, err)
	assert.NotEmpty(t, found)
	assert.Equal(t, show.Place, found[0].Place)
}	
func TestGetShowByID(t *testing.T) {
	db := SetupTestDB()
	service := NewConcert(db)
	artist := Artist{Nom: "Drake", Genre: "Rock"}
	db.Create(&artist)
	show := Show{ArtistID: artist.ID, Place: "Paris", Date: time.Now()}
	db.Create(&show)
	found, err := service.GetShowByID(show.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, found)
	assert.Equal(t, show.Place, found.Place)
}	
func TestSetShow(t *testing.T) {
	db := SetupTestDB()
	service := NewConcert(db)
	artist := Artist{Nom: "Drake", Genre: "Rock"}
	db.Create(&artist)
	show := Show{ArtistID: artist.ID, Place: "Paris", Date: time.Now()}
	db.Create(&show)
	found, err := service.SetShow(show)
	assert.NoError(t, err)
	assert.NotEmpty(t, found)
	assert.Equal(t, show.Place, found.Place)
}		
func TestSetArtist(t *testing.T) {
	db := SetupTestDB()
	service := NewConcert(db)
	artist := Artist{Nom: "Drake", Genre: "Rock"}
	db.Create(&artist)
	found, err := service.SetArtist(artist)
	assert.NoError(t, err)
	assert.NotEmpty(t, found)
	assert.Equal(t, artist.Nom, found.Nom)
}	
func TestParticipateShow(t *testing.T) {
	db := SetupTestDB()
	service := NewConcert(db)
	artist := Artist{Nom: "Drake", Genre: "Rock"}
	db.Create(&artist)
	show := Show{ArtistID: artist.ID, Place: "Paris", Date: time.Now()}
	db.Create(&show)
	fan := Fan{Nom: "Abdou", ShowID: show.ID, Price: 100}
	db.Create(&fan)
	found, err := service.ParticipateShow(fan)
	assert.NoError(t, err)
	assert.NotEmpty(t, found)
	assert.Equal(t, fan.Nom, found.Fans[0].Nom)
}