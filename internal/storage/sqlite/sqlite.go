package sqlite

import (
	"database/sql"
	"os"

	"github.com/pyuldashev912/tracker/internal/storage"
	"github.com/pyuldashev912/tracker/pkg/e"

	_ "github.com/mattn/go-sqlite3"
)

// Storage
type Storage struct {
	db *sql.DB
}

// New
func New(path string, name string) (*Storage, error) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, e.Wrap("can't create folders for storage", err)
	}

	db, err := sql.Open("sqlite3", path+name)
	if err != nil {
		return nil, e.Wrap("can't open database", err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap("can't connect to database", err)
	}

	return &Storage{
		db: db,
	}, nil
}

// CreateUser
func (s *Storage) CreateUser(user *storage.User) error {

	// If user exists we didn't save him again
	if s.isUserExists(user) {
		return nil
	}

	query := `INSERT INTO users(username, telegram_id) VALUES (?, ?)`
	_, err := s.db.Exec(query, user.Username, user.TelegramID)
	if err != nil {
		return e.Wrap("can't insert new user to database", err)
	}

	return nil
}

// SaveTvShow
func (s *Storage) SaveTvShow(tvShow *storage.TvShow) error {
	query := `
	INSERT INTO tv_shows(name, season, episode, users_telegram_id)
	VALUES (?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		tvShow.Name,
		tvShow.Season,
		tvShow.Episode,
		tvShow.UsersTelegramID,
	)

	if err != nil {
		return e.Wrap("can't save new TV Show", err)
	}

	return nil
}

// UpdateLastWatchedEpisode
func (s *Storage) UpdateLastWatchedEpisode(tvShow *storage.TvShow) error {
	query := `
	UPDATE tv_shows
	SET episode=?
	WHERE name=? AND season=? AND users_telegram_id=?
	`

	_, err := s.db.Exec(
		query,
		tvShow.Episode,
		tvShow.Name,
		tvShow.Season,
		tvShow.UsersTelegramID,
	)

	if err != nil {
		return e.Wrap("can't update watched episode", err)
	}

	return nil
}

func (s *Storage) IsTvShowExists(tvShow *storage.TvShow) (bool, error) {
	query := `
	SELECT COUNT(*) FROM tv_shows WHERE name=? AND users_telegram_id=?
	`

	var count int
	if err := s.db.QueryRow(
		query,
		tvShow.Name,
		tvShow.UsersTelegramID,
	).Scan(&count); err != nil {
		return false, e.Wrap("can't check if TV Show exists", err)
	}

	return count > 0, nil
}

// ListAllTvShows
func (s *Storage) ListAllTvShows(userTelegramID int) ([]*storage.TvShow, error) {
	errMsg := "can't list TV Shows"

	query := `SELECT * FROM tv_shows WHERE users_telegram_id=?`
	rows, err := s.db.Query(query, userTelegramID)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}
	defer rows.Close()

	var tempID int
	res := make([]*storage.TvShow, 0, 1)
	for rows.Next() {
		tvShow := new(storage.TvShow)
		err := rows.Scan(
			&tempID,
			&tvShow.Name,
			&tvShow.Season,
			&tvShow.Episode,
			&tvShow.UsersTelegramID,
		)
		if err != nil {
			return nil, e.Wrap(errMsg, err)
		}

		res = append(res, tvShow)
	}

	if err = rows.Err(); err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	return res, nil
}

// RemoveTvShow
func (s *Storage) RemoveTvShow(tvShow *storage.TvShow) error {
	query := `DELETE FROM tv_shows WHERE name=? AND users_telegram_id=?`

	_, err := s.db.Exec(query, tvShow.Name, tvShow.UsersTelegramID)
	if err != nil {
		return e.Wrap("can't remove Tv Show", err)
	}

	return nil
}

// Init
func (s *Storage) Init() error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			username TEXT NOT NULL,
			telegram_id INTEGER NOT NULL
		);

		CREATE TABLE IF NOT EXISTS tv_shows (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			season INTEGER NOT NULL,
			episode INTEGER NOT NULL,
			users_telegram_id INTEGER NOT NULL,
			FOREIGN KEY (users_telegram_id) REFERENCES users(telegram_id)
		);
		`
	_, err := s.db.Exec(query)
	if err != nil {
		return e.Wrap("can't create tables", err)
	}

	return nil
}

func (s *Storage) isUserExists(user *storage.User) bool {
	query := `SELECT COUNT(*) FROM users WHERE telegram_id=?`

	var count int
	s.db.QueryRow(query, user.Username, user.TelegramID).Scan(&count)

	return count > 0
}
