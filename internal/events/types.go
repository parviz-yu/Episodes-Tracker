package events

import "github.com/pyuldashev912/Episodes-Tracker/internal/storage"

// Fetcher
type Fetcher interface {
	Fetch(int, map[int]*State) ([]Event, error)
}

// Processor
type Processor interface {
	Process(*Event, map[int]*State) error
}

// Type represents the type of update from Telegram.
type Type int

const (
	Unknown Type = iota
	Message
	Callback
)

// Event
type Event struct {
	// Type represents the type of update from Telegram.
	//
	// This bot can process two types of update: message and callback query.
	Type Type

	// Text represent the Telegram message.
	Text string

	// ChatID represents the users ID.
	ChatID int

	// InlineMsgID represents the ID from callback's message.
	InlineMsgID int

	// CallbackID represents the ID of callback.
	CallbackID string

	Username  string
	FirstName string
}

// ActiveShow represent an active show that bot tracks in real time.
type ActiveShow struct {
	Name    string
	Season  int
	Episode int
}

// State stores objects that track user states
type State struct {
	SavedShows   []*storage.TvShow
	Prefix       string
	ActiveShow   ActiveShow
	PagBegin     int
	SelectedShow int
	IsPrefixSet  bool
}

type ReplyKeyboardMarkup struct {
	Keyboard   [][]KeyboardButton `json:"keyboard"`
	Persistent bool               `json:"is_persistent"`
	Resize     bool               `json:"resize_keyboard"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}

type InlineKeyboardMarkup struct {
	Inline [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text     string `json:"text"`
	Callback string `json:"callback_data"`
}
