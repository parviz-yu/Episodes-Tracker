package events

import "github.com/pyuldashev912/tracker/internal/storage"

// Fetcher
type Fetcher interface {
	Fetch(int, *Meta) ([]Event, error)
}

// Processor
type Processor interface {
	Process(*Event, *Meta) error
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

// Meta stores objects that track user states
type Meta struct {
	IsPrefixSet  bool
	Prefix       string
	ActiveShows  map[int]ActiveShow
	SavedShows   map[int][]*storage.TvShow
	PagBegin     int
	SelectedShow int
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
