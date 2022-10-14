package adapter

import (
	forwarderModel "bkel-fetching/forwarder/model"
	moodleModel "bkel-fetching/moodle/model"
	"fmt"
	"html"
	"jaytaylor.com/html2text"
	"log"
)

var lineSplit = "--------------------"

func RenderMessage(msg moodleModel.Message) (renderedMsg forwarderModel.Message) {
	renderedMsg.Id = msg.Id
	renderedMsg.UserFrom = msg.UserFrom
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

	renderedMsg.MarkdownText = fmt.Sprintf("%s\nFrom: %s\n%s\n%s", lineSplit, msg.UserFrom, lineSplit, msgInMarkdown)
	return renderedMsg
}
