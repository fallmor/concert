package concert

import (
	"concert/internal/models"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	Db *gorm.DB
}

// type Show struct {
// 	gorm.Model
// 	Title          *string   `json:"title,omitempty"`
// 	Date           time.Time `gorm:"not null" json:"date"`
// 	Time           string    `json:"time"`
// 	ArtistID       uint      `gorm:"not null" json:"artist_id"`
// 	Artist         Artist    `gorm:"foreignKey:ArtistID;references:ID" json:"artist"`
// 	Venue          string    `gorm:"not null" json:"venue"`
// 	Price          float64   `gorm:"not null" json:"price"`
// 	TotalSeats     int       `gorm:"not null" json:"total_seats"`
// 	AvailableSeats int       `gorm:"not null" json:"available_seats"`
// 	Description    string    `json:"description,omitempty"`
// 	ImageURL       string    `json:"image_url,omitempty"`
// 	Fans           []Fan     `gorm:"foreignKey:ShowID" json:"-"`
// }

// type Fan struct {
// 	gorm.Model
// 	Name   string
// 	ShowID uint `gorm:"not null"`
// 	Show   Show `gorm:"foreignKey:ShowID;references:ID"`
// 	Price  int
// 	UserID uint
// 	//User   User `gorm:"foreignKey:UserID;references:ID"`
// }

type ArtistShow struct {
	Artist models.Artist `json:"artist"`
	Shows  []models.Show `json:"shows"`
}

func NewConcert(db *gorm.DB) *Service {
	return &Service{
		Db: db,
	}
}

func (s Service) GetFan(username string) ([]models.Booking, error) {
	var fan []models.Booking
	if result := s.Db.
		Joins("JOIN users ON bookings.user_id = users.id").
		Joins("JOIN shows ON bookings.show_id = shows.id").
		Joins("JOIN artists ON shows.artist_id = artists.id").
		Preload("Show").
		Preload("User").
		Preload("Show.Artist").
		Where("users.username = ?", username).
		Order("bookings.created_at DESC").
		Find(&fan); result.Error != nil {
		return fan, result.Error
	}
	return fan, nil
}

func (s Service) GetShow(artistName string) ([]models.Show, error) {
	var shows []models.Show
	if result := s.Db.
		Preload("Artist").
		Joins("INNER JOIN artists ON shows.artist_id = artists.id").
		Where("Name = ?", artistName).
		Find(&shows); result.Error != nil {
		return nil, result.Error
	}
	return shows, nil
}
func (s Service) GetShowByID(id uint) (models.Show, error) {
	var show models.Show
	if result := s.Db.
		Preload("Artist").
		First(&show, id); result.Error != nil {
		return show, result.Error
	}
	var fanCount int64
	s.Db.Model(&models.Booking{}).Where("show_id = ?", show.ID).Count(&fanCount)
	show.AvailableSeats = show.TotalSeats - int(fanCount)

	return show, nil
}

func (s Service) SetShow(show models.Show) (models.Show, error) {
	if result := s.Db.Save(&show); result.Error != nil {
		return models.Show{}, result.Error
	}
	return show, nil
}
func (s Service) SetArtist(artist models.Artist) (models.Artist, error) {
	if result := s.Db.Save(&artist); result.Error != nil {
		return models.Artist{}, result.Error
	}

	if result := s.Db.First(&artist, artist.ID); result.Error != nil {
		return models.Artist{}, result.Error
	}

	return artist, nil
}

func (s Service) ParticipateShow(info models.Booking) (models.Show, error) {
	info.CreatedAt = time.Now()
	if result := s.Db.Save(&info); result.Error != nil {
		return models.Show{}, result.Error
	}
	return info.Show, nil
}

func (s Service) ListAllShow() ([]models.Show, error) {
	var shows []models.Show
	if result := s.Db.Preload("Artist").Find(&shows); result.Error != nil {
		return nil, result.Error
	}
	return shows, nil
}

func (s Service) ListAllFan() ([]models.Booking, error) {
	var fans []models.Booking
	if result := s.Db.Preload("Show").Preload("Show.Artist").Find(&fans); result.Error != nil {
		return nil, result.Error
	}
	return fans, nil
}

func (s Service) ListAllArtists() ([]models.Artist, error) {
	var artists []models.Artist
	if result := s.Db.Find(&artists); result.Error != nil {
		return nil, result.Error
	}
	return artists, nil
}

func (s Service) GetUserByID(id uint) (User, error) {
	var user User
	if result := s.Db.First(&user, id); result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func (s Service) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if result := s.Db.Find(&users); result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (s Service) GetArtistByID(id uint) (ArtistShow, error) {
	var artist models.Artist
	if result := s.Db.First(&artist, id); result.Error != nil {
		return ArtistShow{}, result.Error
	}
	var shows []models.Show
	s.Db.Where("artist_id = ?", id).Find(&shows)

	return ArtistShow{Artist: artist, Shows: shows}, nil
}

func (s Service) GetFanByID(id uint) (models.Booking, error) {
	var fan models.Booking
	if result := s.Db.Preload("Show").Preload("Show.Artist").First(&fan, id); result.Error != nil {
		return models.Booking{}, result.Error
	}
	return fan, nil
}

func (s Service) DeleteShow(id uint) error {
	if result := s.Db.Delete(&models.Show{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s Service) DeleteArtist(id uint) error {
	if result := s.Db.Delete(&models.Artist{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s Service) DeleteFan(id uint) error {
	if result := s.Db.Delete(&models.Booking{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s Service) UpdateShow(show models.Show) (models.Show, error) {
	if result := s.Db.Save(&show); result.Error != nil {
		return models.Show{}, result.Error
	}
	return show, nil
}

func (s Service) UpdateArtist(artist models.Artist) (models.Artist, error) {
	if result := s.Db.Save(&artist); result.Error != nil {
		return models.Artist{}, result.Error
	}
	return artist, nil
}

func (s Service) UpdateFan(fan models.Booking) (models.Booking, error) {
	if result := s.Db.Save(&fan); result.Error != nil {
		return models.Booking{}, result.Error
	}
	return fan, nil
}

func (s Service) GetMyBookings(user models.User) ([]models.Booking, error) {
	var bookings []models.Booking
	if err := s.Db.
		Preload("Show.Artist").
		Where("user_id = ? AND status = ?", user.ID, "confirmed").
		Order("created_at DESC").
		Find(&bookings).Error; err != nil {
		return bookings, err
	}
	return bookings, nil
}

func (s Service) GetBookingById(id uint) (models.Booking, error) {
	var booking models.Booking
	if err := s.Db.Preload("Show").First(&booking, id).Error; err != nil {
		return models.Booking{}, nil
	}
	return booking, nil
}
