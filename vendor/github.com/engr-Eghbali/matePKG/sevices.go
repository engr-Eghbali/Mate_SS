package services

import (
	"log"
	"net/smtp"
)

type MailOrigin struct {
	From     string
	Password string
}

////////////send mail
func SendMail(body string, recipient string, origin MailOrigin) (er bool) {

	msg := "From: " + origin.From + "\n" +
		"To: " + recipient + "\n" +
		"Subject: Your verification code is : \n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", origin.From, origin.Password, "smtp.gmail.com"),
		origin.From, []string{recipient}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return false
	}

	log.Print("verification code sent to : " + recipient)
	return true
}

////////////////////////////////////////////////////////
