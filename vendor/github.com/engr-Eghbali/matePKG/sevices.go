package services

import (
	"log"
	"net/smtp"
)

////////////send mail
func SendMail(body string, recipient string) (er bool) {
	from := "whereismymate.app@gmail.com"
	pass := "Wakeuptrane2sfc$"
	to := recipient

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Your verification code is : \n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return false
	}

	log.Print("verification code sent to : " + recipient)
	return true
}

////////////////////////////////////////////////////////
