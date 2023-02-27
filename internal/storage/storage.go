package storage

type Storage interface {
	CreateUser(*User) error
	SaveTvShow(*TvShow) error
	UpdateLastWatchedEpisode(*TvShow) error
	IsTvShowExists(*TvShow) (bool, error)
	ListAllTvShows(int) ([]*TvShow, error)
}

type TvShow struct {
	Name            string
	Season          int
	Episode         int
	UsersTelegramID int
}

type User struct {
	TelegramID int
	Username   string
}
