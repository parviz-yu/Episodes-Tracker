package telegram

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/pyuldashev912/Episodes-Tracker/internal/client"
	"github.com/pyuldashev912/Episodes-Tracker/internal/events"
	"github.com/pyuldashev912/Episodes-Tracker/internal/storage"
	"github.com/pyuldashev912/Episodes-Tracker/pkg/e"
)

const (
	addCmdA    = "/add"
	addCmdB    = "➕ Add"
	updCmdA    = "/upd"
	updCmdB    = "🔄 Update"
	listCmdA   = "/list"
	listCmdB   = "📝 List"
	cancelCmdA = "/cancel"
	cancelCmdB = "❌ Cancel"

	startCmd = "/start"
	helpCmd  = "/help"
)

// Represents
const paginationLimit = 5

var (
	ErrInvalidInput = errors.New("invalid input")
)

func (p *Processor) doCommand(event *events.Event, states map[int]*events.State) error {
	defer panicHandler()

	text := strings.TrimSpace(event.Text)

	log.Printf("[INFO] got new command '%s' from '%s' \n", text, event.Username)

	params := client.Params{}
	params.AddParam("chat_id", event.ChatID)

	switch {
	case strings.HasSuffix(text, cancelCmdA) || strings.HasSuffix(text, cancelCmdB):
		states[event.ChatID].Prefix = ""
		states[event.ChatID].IsPrefixSet = false
		params.AddParam("reply_markup", mainKeyboard)
		if show := states[event.ChatID].ActiveShow; show.Name != "" {
			msg := fmt.Sprintf(msgCancel, show.Name)
			params.AddParam("text", msg)
		} else {
			params.AddParam("text", msgCancelNew)
		}

		return p.tg.SendMessage(params)

	case strings.HasPrefix(text, addCmdA) || strings.HasPrefix(text, addCmdB):

		if !states[event.ChatID].IsPrefixSet {
			params.AddParam("reply_markup", cancelKeyboard)
			params.AddParam("text", msgAddInfo)
			states[event.ChatID].Prefix = "/add"
			states[event.ChatID].IsPrefixSet = true
			return p.tg.SendMessage(params)
		}

		return p.addNewTvShow(event, states)

	case strings.HasPrefix(text, updCmdA) || strings.HasPrefix(text, updCmdB):

		if states[event.ChatID].ActiveShow.Name == "" {
			params.AddParam("text", msgNoAddedShows)
			return p.tg.SendMessage(params)
		}

		if !states[event.ChatID].IsPrefixSet {
			params.AddParam("reply_markup", cancelKeyboard)
			msg := fmt.Sprintf(msgUpdateInfo, states[event.ChatID].ActiveShow.Name)
			params.AddParam("text", msg)
			states[event.ChatID].Prefix = "/upd"
			states[event.ChatID].IsPrefixSet = true
			return p.tg.SendMessage(params)
		}

		return p.updateTvShow(event, states)

	case strings.HasSuffix(text, listCmdA) || strings.HasSuffix(text, listCmdB):
		states[event.ChatID].PagBegin = 0
		return p.listTvShows(event, states)

	case text == startCmd:
		states[event.ChatID] = &events.State{}
		return p.addNewUser(event)
	case text == helpCmd:
		params.AddParam("text", msgHelp)
		params.AddParam("parse_mode", "Markdown")
		params.AddParam("disable_web_page_preview", "True")
		return p.tg.SendMessage(params)
	default:
		params.AddParam("text", msgUnknownCommand)
		return p.tg.SendMessage(params)
	}
}

func (p *Processor) listTvShows(event *events.Event, states map[int]*events.State) error {
	shows, err := p.storage.ListAllTvShows(event.ChatID)
	if err != nil {
		return e.Wrap("can't list tv shows", err)
	}

	params := client.Params{}
	params.AddParam("chat_id", event.ChatID)

	if len(shows) == 0 {
		params.AddParam("text", msgNoAddedShows)
		return p.tg.SendMessage(params)
	}

	states[event.ChatID].SavedShows = make([]*storage.TvShow, len(shows))
	copy(states[event.ChatID].SavedShows, shows)

	answer, markup := buildInlineList(shows, states[event.ChatID].PagBegin)
	params.AddParam("text", answer)
	params.AddParam("reply_markup", *markup)

	if err := p.tg.SendMessage(params); err != nil {
		return err
	}

	return nil
}

func (p *Processor) updateTvShow(event *events.Event, states map[int]*events.State) error {
	params := client.Params{}
	params.AddParam("chat_id", event.ChatID)
	params.AddParam("reply_markup", cancelKeyboard)

	episode, err := strconv.Atoi(strings.SplitN(event.Text, " ", 2)[1])
	if err != nil {
		params.AddParam("text", msgInvalidEpisode)
		return p.tg.SendMessage(params)
	}

	updatedTvShow := &storage.TvShow{
		Name:            states[event.ChatID].ActiveShow.Name,
		Season:          states[event.ChatID].ActiveShow.Season,
		Episode:         episode,
		UsersTelegramID: event.ChatID,
	}

	if err := p.storage.UpdateLastWatchedEpisode(updatedTvShow); err != nil {
		return e.Wrap("can't update last watched episode", err)
	}

	states[event.ChatID].ActiveShow = events.ActiveShow{
		Name:    states[event.ChatID].ActiveShow.Name,
		Season:  states[event.ChatID].ActiveShow.Season,
		Episode: episode,
	}

	msg := fmt.Sprintf(msgUpdated, updatedTvShow.Name)
	params.AddParam("text", msg)
	params.AddParam("reply_markup", mainKeyboard)
	if err := p.tg.SendMessage(params); err != nil {
		return e.Wrap("can't update last watched episode", err)
	}

	states[event.ChatID].Prefix = ""
	states[event.ChatID].IsPrefixSet = false

	return nil
}

func (p *Processor) addNewTvShow(event *events.Event, states map[int]*events.State) error {
	errMsg := "can't add new Tv Show"

	params := client.Params{}
	params.AddParam("chat_id", event.ChatID)
	params.AddParam("reply_markup", cancelKeyboard)

	// Get inputs after second split
	inputs := strings.Split(strings.SplitN(event.Text, " ", 2)[1], "/")
	if len(inputs) != 3 {
		params.AddParam("text", msgInvalidInput)
		return p.tg.SendMessage(params)
	}

	season, err := strconv.Atoi(inputs[1])
	if err != nil {
		params.AddParam("text", msgInvalidInput)
		return p.tg.SendMessage(params)
	}

	episode, err := strconv.Atoi(inputs[2])
	if err != nil {
		params.AddParam("text", msgInvalidInput)
		return p.tg.SendMessage(params)
	}

	show := &storage.TvShow{
		Name:            inputs[0],
		Season:          season,
		Episode:         episode,
		UsersTelegramID: event.ChatID,
	}

	exists, err := p.storage.IsTvShowExists(show)
	if err != nil {
		return e.Wrap(errMsg, err)
	}

	if exists {
		msg := fmt.Sprintf(msgAlreadyExists, show.Name)
		params.AddParam("text", msg)
		return p.tg.SendMessage(params)
	}

	if err := p.storage.SaveTvShow(show); err != nil {
		return e.Wrap(errMsg, err)
	}

	states[event.ChatID].ActiveShow = events.ActiveShow{
		Name:    inputs[0],
		Season:  season,
		Episode: episode,
	}

	msg := fmt.Sprintf(msgAdded, show.Name)
	params.AddParam("text", msg)
	params.AddParam("reply_markup", mainKeyboard)
	if err := p.tg.SendMessage(params); err != nil {
		return e.Wrap(errMsg, err)
	}

	states[event.ChatID].Prefix = ""
	states[event.ChatID].IsPrefixSet = false

	return nil
}

func (p *Processor) addNewUser(event *events.Event) error {
	user := &storage.User{
		TelegramID: event.ChatID,
		Username:   event.Username,
	}

	if err := p.storage.CreateUser(user); err != nil {
		return e.Wrap("can't add new user", err)
	}

	params := client.Params{}
	params.AddParam("chat_id", event.ChatID)
	params.AddParam("text", fmt.Sprintf(msgHello, event.FirstName))
	params.AddParam("parse_mode", "Markdown")
	params.AddParam("disable_web_page_preview", "True")
	params.AddParam("reply_markup", mainKeyboard)

	if err := p.tg.SendMessage(params); err != nil {
		e.Wrap("can't add new user", err)
	}

	return nil
}

// slice, pagBegin
func buildInlineList(shows []*storage.TvShow, begin int) (string, *events.InlineKeyboardMarkup) {
	var answer strings.Builder

	baseInline := baseInlineMarkup()

	end := 0
	if paginationLimit+begin < len(shows) {
		end = paginationLimit + begin
	} else {
		end = len(shows)
	}

	for i := begin; i < end; i++ {
		num := strconv.Itoa(i + 1)
		answer.WriteString(fmt.Sprintf("%s. %s\n", num, shows[i].Name))

		baseInline.Inline[0] = append(
			baseInline.Inline[0],
			events.InlineKeyboardButton{
				Text:     num,
				Callback: num,
			},
		)
	}

	if begin+paginationLimit < len(shows) {
		if begin != 0 {
			baseInline.Inline = append(
				baseInline.Inline,
				[]events.InlineKeyboardButton{
					{
						Text:     "< Previous",
						Callback: "Back",
					},
					{
						Text:     "Next >",
						Callback: "Forward",
					},
				},
			)
		} else {
			baseInline.Inline = append(
				baseInline.Inline,
				[]events.InlineKeyboardButton{
					{
						Text:     "Next >",
						Callback: "Forward",
					},
				},
			)
		}
	}

	if begin != 0 && begin+paginationLimit >= len(shows) {
		baseInline.Inline = append(
			baseInline.Inline,
			[]events.InlineKeyboardButton{
				{
					Text:     "< Previous",
					Callback: "Back",
				},
			},
		)
	}

	return answer.String(), &baseInline
}

func baseInlineMarkup() events.InlineKeyboardMarkup {
	return events.InlineKeyboardMarkup{
		Inline: make([][]events.InlineKeyboardButton, 1),
	}
}
