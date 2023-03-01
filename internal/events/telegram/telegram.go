package telegram

import (
	"errors"
	"fmt"

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

func New(client *client.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
		offset:  0,
	}
}

func (p *Processor) Fetch(limit int, meta *events.Meta) ([]events.Event, error) {
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
		res = append(res, event(&upd, meta))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func event(upd *client.Update, meta *events.Meta) events.Event {
	updType := fetchType(upd)
	updText := fetchText(upd)
	res := events.Event{
		Type: updType,
		Text: fmt.Sprintf("%s %s", meta.Prefix, updText),
	}

	if updType == events.Message {
		res.ChatID = upd.Message.Chat.ID
		res.Username = upd.Message.From.Username
		res.FirstName = upd.Message.From.FirstName
	}

	return res
}

func fetchType(upd *client.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}

func fetchText(upd *client.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}
