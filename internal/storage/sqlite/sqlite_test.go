package sqlite

import (
	"fmt"
	"os"
	"testing"

	"github.com/pyuldashev912/tracker/internal/storage"
	"github.com/stretchr/testify/assert"
)

// Truncate
func (s *Storage) Truncate(t *testing.T, tableName string) {
	if _, err := s.db.Exec(fmt.Sprintf("DELETE FROM %s;", tableName)); err != nil {
		t.Fatal(err)
	}

	s.db.Close()

	if err := os.RemoveAll("data"); err != nil {
		t.Fatal(err)
	}
}

func DBHelper(t *testing.T) *Storage {
	path := "data/storage"
	t.Helper()

	storage, err := New(path, "/test-storage.db")
	if err != nil {
		t.Fatal(err)
	}

	storage.Init()

	return storage
}

func TestStorage_CreateUser(t *testing.T) {
	user := storage.User{
		Username:   "pyuldashev",
		TelegramID: 123354654684564,
	}

	s := DBHelper(t)
	err := s.CreateUser(&user)
	assert.NoError(t, err)

	err = s.CreateUser(&user)
	assert.NoError(t, err)

	s.Truncate(t, "users")
}

func TestStorage_SaveTvShow(t *testing.T) {
	tvShow := storage.TvShow{
		Name:            "Silicon Valley",
		Season:          1,
		Episode:         1,
		UsersTelegramID: 123354654684564,
	}

	s := DBHelper(t)
	err := s.SaveTvShow(&tvShow)
	assert.NoError(t, err)

	s.Truncate(t, "tv_shows")
}

func TestStorage_UpdateLastWatchedEpisode(t *testing.T) {
	tvShow := storage.TvShow{
		Name:            "Silicon Valley",
		Season:          1,
		Episode:         1,
		UsersTelegramID: 123354654684564,
	}

	newTvShow := storage.TvShow{
		Name:            "Silicon Valley",
		Season:          1,
		Episode:         3,
		UsersTelegramID: 123354654684564,
	}

	s := DBHelper(t)
	s.SaveTvShow(&tvShow)
	err := s.UpdateLastWatchedEpisode(&newTvShow)
	assert.NoError(t, err)
	s.Truncate(t, "tv_shows")
}

func TestStorage_IsTvShowExists(t *testing.T) {
	tvShow := storage.TvShow{
		Name:            "Silicon Valley",
		Season:          1,
		Episode:         1,
		UsersTelegramID: 123354654684564,
	}

	s := DBHelper(t)
	s.SaveTvShow(&tvShow)
	result, err := s.IsTvShowExists(&tvShow)
	assert.Equal(t, true, result)
	assert.NoError(t, err)

	tvShow.Name = "Friends"
	result, err = s.IsTvShowExists(&tvShow)
	assert.Equal(t, false, result)
	assert.NoError(t, err)
	s.Truncate(t, "tv_shows")
}

func TestStorage_ListAllTvShows(t *testing.T) {
	tvShow := storage.TvShow{
		Name:            "Silicon Valley",
		Season:          1,
		Episode:         1,
		UsersTelegramID: 123354654684564,
	}

	s := DBHelper(t)
	s.SaveTvShow(&tvShow)

	result, err := s.ListAllTvShows(123354654684564)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	s.Truncate(t, "tv_shows")
}

func TestStorage_RemoveTvShows(t *testing.T) {
	tvShow := storage.TvShow{
		Name:            "Silicon Valley",
		Season:          1,
		Episode:         1,
		UsersTelegramID: 123354654684564,
	}

	s := DBHelper(t)
	s.SaveTvShow(&tvShow)

	err := s.RemoveTvShow(&tvShow)
	assert.NoError(t, err)
	s.Truncate(t, "tv_shows")
}
