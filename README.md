# Episodes Tracker

«Episodes Tracker» is a Telegram Bot to track watched episodes of TV series. The bot is written based on [documentation](https://core.telegram.org/bots/api) without the use of third-party libraries and API wrappers.

![](preview.gif)

## Setup
1. Get your telegram bot token

   Create a bot from Telegram [@BotFather](https://t.me/BotFather) and obtain an access token.

2. Clone this repository and `cd` into the project source directory

3. Build using `go build`

    If you have a Go environment, you can build it with the following command:

```bash
go build -o episode-tracker cmd/main.go 
```
4. Set the environment variable and run

```bash
export BOT_TOKEN=<your_telegram_bot_token>
./episode-tracker 
```