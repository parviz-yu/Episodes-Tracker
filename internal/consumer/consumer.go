package consumer

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pyuldashev912/tracker/internal/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) CheckToken(token string) error {
	URL := fmt.Sprintf("https://api.telegram.org/bot%s/getMe", token)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("invalid Token")
	}

	return nil
}

func (c *Consumer) Start() error {
	states := make(map[int]*events.State)

	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize, states)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := c.handleEvents(gotEvents, states); err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}

	}
}

func (c *Consumer) handleEvents(events []events.Event, states map[int]*events.State) error {
	for _, event := range events {
		if err := c.processor.Process(&event, states); err != nil {
			log.Printf("[ERR] can't handle event: %s", err.Error())

			continue
		}
		fmt.Printf("%#v\n\n", states[event.ChatID])
	}

	return nil
}
