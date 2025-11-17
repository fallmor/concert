package concert

import (
	"time"

	"gorm.io/gorm"
)

type Service struct {
	Db *gorm.DB
}

type Artist struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Nom       string
	Genre     string
	PhotoURL  string `gorm:"type:varchar(500)"`
	AlbumURL  string `gorm:"type:varchar(500)"`
}

type Show struct {
	gorm.Model
	Date     time.Time
	ArtistID uint   `gorm:"not null"`
	Artist   Artist `gorm:"foreignKey:ArtistID;references:ID"`
	Place    string
	Fans     []Fan `gorm:"foreignKey:ShowID"`
}

type Fan struct {
	gorm.Model
	Nom    string
	ShowID uint `gorm:"not null"`
	Show   Show `gorm:"foreignKey:ShowID;references:ID"`
	Price  int
	UserID uint
	//User   User `gorm:"foreignKey:UserID;references:ID"`
}

func NewConcert(db *gorm.DB) *Service {
	return &Service{
		Db: db,
	}
}

func (s Service) GetFan(name string) ([]Fan, error) {
	var fan []Fan
	if result := s.Db.
		Joins("LEFT JOIN shows ON fans.show_id = shows.id").
		Joins("LEFT JOIN artists ON shows.artist_id = artists.id").
		Preload("Show").
		Preload("Show.Artist").
		Where("fans.nom = ?", name).
		Find(&fan); result.Error != nil {
		return fan, result.Error
	}
	return fan, nil
}

func (s Service) GetShow(artistName string) ([]Show, error) {
	var shows []Show
	if result := s.Db.
		Preload("Artist").
		Joins("INNER JOIN artists ON shows.artist_id = artists.id").
		Where("nom = ?", artistName).
		Find(&shows); result.Error != nil {
		return nil, result.Error
	}
	return shows, nil
}
func (s Service) GetShowByID(id uint) (Show, error) {
	var show Show
	if result := s.Db.
		Preload("Artist").
		Preload("Fans", func(db *gorm.DB) *gorm.DB {
			return db.Order("fans.created_at DESC")
		}).
		First(&show, id); result.Error != nil {
		return show, result.Error
	}
	return show, nil
}

func (s Service) SetShow(show Show) (Show, error) {
	if result := s.Db.Save(&show); result.Error != nil {
		return Show{}, result.Error
	}
	return show, nil
}
func (s Service) SetArtist(artist Artist) (Artist, error) {
	if result := s.Db.Save(&artist); result.Error != nil {
		return Artist{}, result.Error
	}
	
	if result := s.Db.First(&artist, artist.ID); result.Error != nil {
		return Artist{}, result.Error
	}
	
	return artist, nil
}

func (s Service) ParticipateShow(info Fan) (Show, error) {
	info.CreatedAt = time.Now()
	if result := s.Db.Save(&info); result.Error != nil {
		return Show{}, result.Error
	}
	return info.Show, nil
}

func (s Service) ListAllShow() ([]Show, error) {
	var shows []Show
	if result := s.Db.Preload("Artist").Find(&shows); result.Error != nil {
		return nil, result.Error
	}
	return shows, nil
}

func (s Service) ListAllFan() ([]Fan, error) {
	var fans []Fan
	if result := s.Db.Preload("Show").Preload("Show.Artist").Find(&fans); result.Error != nil {
		return nil, result.Error
	}
	return fans, nil
}

func (s Service) ListAllArtists() ([]Artist, error) {
	var artists []Artist
	if result := s.Db.Find(&artists); result.Error != nil {
		return nil, result.Error
	}
	return artists, nil
}
