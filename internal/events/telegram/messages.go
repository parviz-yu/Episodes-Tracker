package telegram

import "github.com/pyuldashev912/Episodes-Tracker/internal/events"

const msgHelp = `I can help you track and manage the episodes of the TV series you watch.

Just click the '➕ Add' button to add new TV series. Everything is explained underway.

To inform the bot about the last episode you watched, please click the '🔄 Update' button to add the episode number.

Forgot which episode you were on? Click '📝 List' to get the last viewed episode for each added TV series.

This bot is open-source on [GitHub](https://github.com/pyuldashev912/Episodes-Tracker).
`

const msgHello = "Hi _%s_! 👋\n\n" + msgHelp

const (
	msgAdded         = "'%s' added! ✅"
	msgSelected      = "'%s' selected for tracking 🔍"
	msgRemoved       = "'%s' successfully removed 🗑 \n\nPlease select a TV series for tracking."
	msgUpdated       = "Last watched episode of '%s' updated 👌"
	msgAlreadyExists = "You have already added '%s' 🤗"

	msgAddInfo    = "Let's add a new TV Show 🎬 \n\n Please use this format: \n SeriesName/Season/LastWatchedEpisode \n e.g. Silicon Valley/1/3"
	msgUpdateInfo = "Input last watched episode of '%s' 🍿"

	msgCancel         = "Which episode of '%s' did you watch? 📺"
	msgCancelNew      = "Feel free to add TV series 😉"
	msgUnknownCommand = "Unknown command 🙈"
	msgNoAddedShows   = "You don't have any TV series added 🤷‍♂️"
	msgInvalidInput   = "Invalid input ⚠️\n\n Please use this format: \n SeriesName/Season/LastWatchedEpisode \n e.g. Silicon Valley/1/3"
	msgInvalidEpisode = "Last watched episode should be integer 🤔 \n Try one more time"
)

var mainKeyboard = events.ReplyKeyboardMarkup{
	Keyboard: [][]events.KeyboardButton{
		{
			{
				Text: "➕ Add",
			},
			{
				Text: "🔄 Update",
			},
			{
				Text: "📝 List",
			},
		},
	},
	Persistent: true,
	Resize:     true,
}

var cancelKeyboard = events.ReplyKeyboardMarkup{
	Keyboard: [][]events.KeyboardButton{
		{
			{
				Text: "❌ Cancel",
			},
		},
	},
	Persistent: true,
	Resize:     true,
}
