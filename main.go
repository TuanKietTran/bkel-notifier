package main

import (
	"bkel-fetching/utils/bot"
	"bkel-fetching/utils/env"
	"bkel-fetching/utils/moodle"
	"fmt"
	"html"
	"jaytaylor.com/html2text"
	"log"
	"os"
)

const LineSplit = "--------------------"
const LogFileName = "bkel-fetching.log"

func main() {
	env.LoadDotEnv()
	logFile := openFileForLogging(LogFileName)

	botToken, moodleToken, chatId := env.GetVars()

	myBot := bot.NewPersonalBot(botToken, chatId)
	log.Println("Bot created successfully")

	moodleId := moodle.GetMoodleID(moodleToken)
	newMessages := moodle.FetchMessages(moodleToken, moodleId)

	for _, msg := range newMessages {
		moodle.MarkChatAsRead(moodleToken, msg.Id)

		msgInMarkdown, err := html2text.FromString(
			msg.Text,
			html2text.Options{
				PrettyTables: true,
				OmitLinks:    false,
			})
		if err != nil {
			log.Printf("Can't parse message to Markdown, err: %v", err)
			msgInMarkdown = html.EscapeString(msg.Text)
		}

		msgWithHeader := fmt.Sprintf("%s\nFrom: %s\n%s\n%s", LineSplit, msg.UserFrom, LineSplit, msgInMarkdown)
		myBot.Send(msg.UserFrom, msgWithHeader)
		log.Printf("Handled message from user: %v", msg.Id)
	}

	moodle.MarkAllNotificationsAsRead(moodleToken)

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
