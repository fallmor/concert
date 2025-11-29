package concert

type ConcertService interface {
	GetFan(name string) ([]Fan, error)
	GetShow(artistName string) ([]Show, error)
	GetShowByID(id uint) (Show, error)
	SetShow(show Show) (Show, error)
	SetArtist(artist Artist) (Artist, error)
	ParticipateShow(fan Fan) (Show, error)
	ListAllShow() ([]Show, error)
	ListAllFan() ([]Fan, error)
	ListAllArtists() ([]Artist, error)
	GetAllUsers() ([]User, error)
	DeleteShow(id uint) error
	DeleteArtist(id uint) error
	DeleteFan(id uint) error
	UpdateShow(show Show) (Show, error)
	UpdateArtist(artist Artist) (Artist, error)
	UpdateFan(fan Fan) (Fan, error)
}
