package forwarder

import (
	"bkmail/forwarder/model"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
)

type PersonalBot struct {
	internalBot *tb.Bot
	handleUser  *tb.User
}

func NewPersonalBot(botToken string, chatId int64) PersonalBot {
	newBot := PersonalBot{}

	var err error
	newBot.internalBot, err = tb.NewBot(tb.Settings{
		Token: botToken,
	})

	if err != nil {
		log.Panicf("Can't create Telegram internalBot, err: %v", err)
	}

	newBot.handleUser = new(tb.User)
	newBot.handleUser.ID = chatId

	return newBot
}

func (bot *PersonalBot) Send(msg model.Message) {
	_, err := bot.internalBot.Send(bot.handleUser, msg.MarkdownText, tb.ModeMarkdown)
	if err != nil {
		log.Printf("Can't send message, err: %v", err)
		bot.SendError(msg.UserFrom)
	}
}

func (bot *PersonalBot) SendError(userFrom string) {
	errorMsg := fmt.Sprintf(
		"Can't send message from handleUser %s. Please check it from BKeL.", userFrom)
	_, err := bot.internalBot.Send(bot.handleUser, errorMsg, tb.ModeDefault)

	if err != nil {
		log.Printf("Can't send error message to handleUser %s, just skipping, err: %v", userFrom, err)
	}
}
