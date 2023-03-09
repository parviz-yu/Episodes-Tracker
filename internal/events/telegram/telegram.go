package telegram

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/pyuldashev912/tracker/internal/client"
	"github.com/pyuldashev912/tracker/internal/events"
	"github.com/pyuldashev912/tracker/internal/storage"
	"github.com/pyuldashev912/tracker/pkg/e"
)

type Processor struct {
	tg      *client.Client
	storage storage.Storage
	offset  int
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
)

// New
func New(client *client.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
		offset:  0,
	}
}

// Fetch
func (p *Processor) Fetch(limit int, states map[int]*events.State) ([]events.Event, error) {
	params := client.Params{}
	params.AddParam("offset", p.offset)
	params.AddParam("limit", limit)

	updates, err := p.tg.Updates(params)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, upd := range updates {
		var id int
		if upd.Message != nil {
			id = upd.Message.Chat.ID
		}

		if upd.Callback != nil {
			id = upd.Callback.Message.Chat.ID
		}

		if val, ok := states[id]; ok {
			res = append(res, event(&upd, val.Prefix))
		} else {
			res = append(res, event(&upd, ""))
		}
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

// Process
func (p *Processor) Process(event *events.Event, states map[int]*events.State) error {
	switch event.Type {
	case events.Message:
		p.doCommand(event, states)
	case events.Callback:
		p.doCallback(event, states)
	default:
		return e.Wrap("can't process event", ErrUnknownEventType)
	}

	return nil
}

func event(upd *client.Update, prefix string) events.Event {
	updType := fetchType(upd)
	updText := fetchText(upd)
	res := events.Event{
		Type: updType,
		Text: fmt.Sprintf("%s %s", prefix, updText),
	}

	if updType == events.Message {
		res.ChatID = upd.Message.Chat.ID
		res.Username = upd.Message.From.Username
		res.FirstName = upd.Message.From.FirstName
	}

	if updType == events.Callback {
		res.CallbackID = upd.Callback.ID
		res.ChatID = upd.Callback.Message.Chat.ID
		res.Text = upd.Callback.Data
		res.InlineMsgID = upd.Callback.Message.ID
		res.Username = strconv.Itoa(upd.Callback.Message.Chat.ID)
	}

	return res
}

func fetchType(upd *client.Update) events.Type {
	if upd.Message == nil && upd.Callback == nil {
		return events.Unknown
	}

	if upd.Message == nil {
		return events.Callback
	}

	return events.Message
}

func fetchText(upd *client.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}
