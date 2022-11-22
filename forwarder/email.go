package forwarder

import (
	"bkmail/forwarder/model"
	"fmt"
	"net/smtp"
)

func SendMailSimple(message model.Message, to string) {
	auth := smtp.PlainAuth(
		"",
		"nguyenthaitan02@gmail.com",
		"vmdgzegojigqhvzh",
		"smtp.gmail.com",
	)

	msg := "Subject: " + message.UserFrom + "\n" + message.MarkdownText

	recipient := []string{to}

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"nguyenthaitan02@gmail.com",
		recipient,
		[]byte(msg),
	)

	if err != nil {
		fmt.Println(err)
	}
}
