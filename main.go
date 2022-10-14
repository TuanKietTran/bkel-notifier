package main

import (
	"bkel-fetching/env"
	"bkel-fetching/forwarder"
	"bkel-fetching/forwarder/adapter"
	"bkel-fetching/moodle"
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

	telegramBot := forwarder.NewPersonalBot(botToken, chatID)
	log.Println("Bot created successfully")

	// New Moodle
	newMoodleClient := moodle.NewClient(newMoodleHostname, newMoodleToken)

	newMessages := newMoodleClient.FetchMessages()

	for _, msg := range newMessages {
		if msg.IsNotification == 1 {
			newMoodleClient.MarkNotificationAsRead(msg.Id)
		} else {
			newMoodleClient.MarkChatAsRead(msg.Id)
		}
		parsedMsg := adapter.RenderMessage(msg)

		log.Printf("Received message from calendar_forward %s", msg.UserFrom)
		telegramBot.Send(parsedMsg)
		log.Printf("Handled message from calendar_forward: %v", msg.Id)
	}

	// Old Moodle
	oldMoodleClient := moodle.NewClient(oldMoodleHostname, newMoodleToken)
	newMessages = oldMoodleClient.FetchMessages()

	for _, msg := range newMessages {
		newMoodleClient.MarkChatAsRead(msg.Id)
		parsedMsg := adapter.RenderMessage(msg)

		telegramBot.Send(parsedMsg)
		log.Printf("Handled message from calendar_forward: %v", msg.Id)
	}

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
