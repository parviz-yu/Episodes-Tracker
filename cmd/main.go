package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pyuldashev912/tracker/internal/client"
	"github.com/pyuldashev912/tracker/internal/consumer"
	"github.com/pyuldashev912/tracker/internal/events/telegram"
	"github.com/pyuldashev912/tracker/internal/storage/sqlite"
)

func main() {
	host := "api.telegram.org"
	token := os.Getenv("BOT_TOKEN")
	storage_path := "data/sqlite"

	cli := client.New(host, token)
	storage, err := sqlite.New(storage_path, "/storage.db")
	if err != nil {
		log.Fatal(err)
	}

	storage.Init()

	eventProcessor := telegram.New(cli, storage)
	c := consumer.New(eventProcessor, eventProcessor, 100)
	fmt.Println("Service started...")

	if err := c.CheckToken(token); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Token is valid")
	c.Start()
}
