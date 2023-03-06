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
func (p *Processor) doCallback(event *events.Event, meta *events.Meta) {

	p.sendAnswer(event)

	params := client.Params{}
	params.AddParam("message_id", event.InlineMsgID)
	params.AddParam("chat_id", event.ChatID)

	switch event.Text {
	case "Back":

		if meta.PagBegin == 0 {
			meta.PagBegin = 0
		} else {
			meta.PagBegin -= paginationLimit
		}
		answer, markup := buildInlineList(meta.SavedShows[event.ChatID], meta.PagBegin)
		params.AddParam("text", answer)
		params.AddParam("reply_markup", *markup)

	case "Forward":
		meta.PagBegin += paginationLimit
		answer, markup := buildInlineList(meta.SavedShows[event.ChatID], meta.PagBegin)
		params.AddParam("text", answer)
		params.AddParam("reply_markup", *markup)

	case "List":
		answer, markup := buildInlineList(meta.SavedShows[event.ChatID], meta.PagBegin)
		params.AddParam("text", answer)
		params.AddParam("reply_markup", *markup)

	case "Select":
		show := meta.SavedShows[event.ChatID][meta.SelectedShow]
		newActiveShow := events.ActiveShow{
			Name:    show.Name,
			Season:  show.Season,
			Episode: show.Episode,
		}
		meta.ActiveShows[event.ChatID] = newActiveShow
		params.AddParam("text", fmt.Sprintf(msgSelected, newActiveShow.Name))
		p.tg.SendMessage(params)
		return

	case "Remove":
		if err := p.removeTvShow(meta.SavedShows[event.ChatID][meta.SelectedShow]); err != nil {
			return
		}

		params.AddParam(
			"text",
			fmt.Sprintf(msgRemoved, meta.SavedShows[event.ChatID][meta.SelectedShow].Name),
		)

		// Removing show from active
		show := meta.SavedShows[event.ChatID][meta.SelectedShow]
		if show.Name == meta.ActiveShows[event.ChatID].Name {
			meta.ActiveShows[event.ChatID] = events.ActiveShow{}
		}

		// Removing show from cache
		meta.SavedShows[event.ChatID] = append(
			meta.SavedShows[event.ChatID][:meta.SelectedShow],
			meta.SavedShows[event.ChatID][meta.SelectedShow+1:]...,
		)

	default:
		num, _ := strconv.Atoi(event.Text)
		meta.SelectedShow = num - 1
		answer, markup := buildInlineItem(meta.SavedShows[event.ChatID], num)

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
	rec := recover()
	if rec != nil {
		log.Println("[RECOVER]", rec)
	}
}
