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
}
