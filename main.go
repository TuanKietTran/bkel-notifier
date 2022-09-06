package main

import (
	"bkel-fetching/utils"
	"bkel-fetching/utils/bot"
	"bkel-fetching/utils/env"
	"bkel-fetching/utils/moodle"
	"log"
	"os"
)

const LogFileName = "bkel-fetching.log"

const newMoodleHostname = "e-learning.hcmut.edu.vn"
const oldMoodleHostname = "e-learning-old.hcmut.edu.vn"

func main() {
	env.LoadDotEnv()
	logFile := openFileForLogging(LogFileName)

	botToken, chatID, _, newMoodleToken := env.GetVars()

	telegramBot := bot.NewPersonalBot(botToken, chatID)
	log.Println("Bot created successfully")

	// New Moodle
	newMoodleClient := moodle.NewClient(newMoodleHostname, newMoodleToken)
	newMessages := newMoodleClient.FetchMessages()

	for _, msg := range newMessages {
		newMoodleClient.MarkChatAsRead(msg.Id)
		parsedMsg := utils.RenderMessage(msg)

		telegramBot.Send(msg.UserFrom, parsedMsg)
		log.Printf("Handled message from user: %v", msg.Id)
	}

	newMoodleClient.MarkAllNotificationsAsRead()

	// Old Moodle
	oldMoodleClient := moodle.NewClient(oldMoodleHostname, newMoodleToken)
	newMessages = oldMoodleClient.FetchMessages()

	for _, msg := range newMessages {
		newMoodleClient.MarkChatAsRead(msg.Id)
		parsedMsg := utils.RenderMessage(msg)

		telegramBot.Send(msg.UserFrom, parsedMsg)
		log.Printf("Handled message from user: %v", msg.Id)
	}

	oldMoodleClient.MarkAllNotificationsAsRead()

	_ = logFile.Close()
}

func openFileForLogging(filename string) *os.File {
	newLogFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Can't open log file, err: %v", err)
	}

	log.SetOutput(newLogFile)
	return newLogFile
}
