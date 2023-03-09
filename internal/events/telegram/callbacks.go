package telegram

import (
	"fmt"
	"log"
	"strconv"

	"github.com/pyuldashev912/tracker/internal/client"
	"github.com/pyuldashev912/tracker/internal/events"
	"github.com/pyuldashev912/tracker/internal/storage"
	"github.com/pyuldashev912/tracker/pkg/e"
)

// TODO: add return
func (p *Processor) doCallback(event *events.Event, states map[int]*events.State) {
	defer panicHandler()

	log.Printf("[INFO] got new callback query '%s' from '%s' \n", event.Text, event.Username)

	p.sendAnswer(event)

	params := client.Params{}
	params.AddParam("message_id", event.InlineMsgID)
	params.AddParam("chat_id", event.ChatID)

	switch event.Text {
	case "Back":

		if states[event.ChatID].PagBegin == 0 {
			states[event.ChatID].PagBegin = 0
		} else {
			states[event.ChatID].PagBegin -= paginationLimit
		}
		answer, markup := buildInlineList(states[event.ChatID].SavedShows, states[event.ChatID].PagBegin)
		params.AddParam("text", answer)
		params.AddParam("reply_markup", *markup)

	case "Forward":
		states[event.ChatID].PagBegin += paginationLimit
		answer, markup := buildInlineList(states[event.ChatID].SavedShows, states[event.ChatID].PagBegin)
		params.AddParam("text", answer)
		params.AddParam("reply_markup", *markup)

	case "List":
		answer, markup := buildInlineList(states[event.ChatID].SavedShows, states[event.ChatID].PagBegin)
		params.AddParam("text", answer)
		params.AddParam("reply_markup", *markup)

	case "Select":
		show := states[event.ChatID].SavedShows[states[event.ChatID].SelectedShow]
		newActiveShow := events.ActiveShow{
			Name:    show.Name,
			Season:  show.Season,
			Episode: show.Episode,
		}
		states[event.ChatID].ActiveShow = newActiveShow
		params.AddParam("text", fmt.Sprintf(msgSelected, newActiveShow.Name))

	case "Remove":
		index := states[event.ChatID].SelectedShow
		tvShow := states[event.ChatID].SavedShows[index]
		if err := p.removeTvShow(tvShow); err != nil {
			return
		}

		params.AddParam(
			"text",
			fmt.Sprintf(msgRemoved, tvShow.Name),
		)

		// Removing show from active
		if tvShow.Name == states[event.ChatID].ActiveShow.Name {
			states[event.ChatID].ActiveShow = events.ActiveShow{}
		}

		// Removing show from cache
		states[event.ChatID].SavedShows = append(
			states[event.ChatID].SavedShows[:index],
			states[event.ChatID].SavedShows[index+1:]...,
		)

	default:
		num, _ := strconv.Atoi(event.Text)
		states[event.ChatID].SelectedShow = num - 1
		answer, markup := buildInlineItem(states[event.ChatID].SavedShows, num)

		params.AddParam("text", answer)
		params.AddParam("reply_markup", markup)
		params.AddParam("parse_mode", "Markdown")
	}

	if err := p.tg.EditMessageText(params); err != nil {
		log.Println(err)
	}

}

func (p *Processor) sendAnswer(event *events.Event) {
	params := client.Params{}
	params.AddParam("callback_query_id", event.CallbackID)
	p.tg.AnswerCallbackQuery(params)
}

func (p *Processor) removeTvShow(tvShow *storage.TvShow) error {
	defer panicHandler()
	if err := p.storage.RemoveTvShow(tvShow); err != nil {
		return e.Wrap("can't remove tv show", err)
	}

	return nil
}

func buildInlineItem(shows []*storage.TvShow, number int) (string, *events.InlineKeyboardMarkup) {
	defer panicHandler()
	answer := fmt.Sprintf(
		"*%s*\nSeason:%d\nLast watched episode:%d",
		shows[number-1].Name,
		shows[number-1].Season,
		shows[number-1].Episode,
	)

	baseInline := baseInlineMarkup()
	baseInline.Inline[0] = append(
		baseInline.Inline[0],
		events.InlineKeyboardButton{
			Text:     "Select",
			Callback: "Select",
		},
	)

	baseInline.Inline[0] = append(
		baseInline.Inline[0],
		events.InlineKeyboardButton{
			Text:     "Remove",
			Callback: "Remove",
		},
	)

	baseInline.Inline = append(
		baseInline.Inline,
		[]events.InlineKeyboardButton{
			{
				Text:     "Back to list",
				Callback: "List",
			},
		},
	)

	return answer, &baseInline
}

func panicHandler() {
	// params := client.Params

	rec := recover()
	if rec != nil {
		log.Println("[RECOVER]", rec)
	}
}
