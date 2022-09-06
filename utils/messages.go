package utils

import (
	"bkel-fetching/model"
	"fmt"
	"html"
	"jaytaylor.com/html2text"
	"log"
)

const LineSplit = "--------------------"

func RenderMessage(msg model.Message) (parsedMsg string) {
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
	return msgWithHeader
}
