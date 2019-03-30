package services

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
)

type MailOrigin struct {
	From     string
	Password string
}

type SmsOrigin struct {
	From   string
	ApiKey string
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

////////////sned sms
func SendSms(txt string, recipient string, origin SmsOrigin) bool {

	resp, err := http.Get("https://login.parsgreen.com/UrlService/sendSMS.ashx?from=" + origin.From + "&to=" + recipient + "&&text=" + txt + "&signature=" + origin.ApiKey)
	defer resp.Body.Close()
	body, err2 := ioutil.ReadAll(resp.Body)
	content := string(body)
	log.Println(content)
	log.Println(err2)
	if err != nil {
		// handle error
		log.Println("Sending SMS to: " + recipient + " Failed 006 <=End")
		return false
	} else {
		log.Println("SMS Sent to: " + recipient)
		return true
	}

}
