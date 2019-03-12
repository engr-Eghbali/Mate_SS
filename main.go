package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"

	//"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

//////////////port detemine
func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

///////////////////////////////////////////////////////////////////

////////////send mail
func send(body string) (bool er) {
	from := "trane2sfc@gmail.com"
	pass := "Wakeuptrane2sfc$"
	to := "engr.eghbali@gmail.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return false
	}

	log.Print("sent, visit http://foobarbazz.mailinator.com")
	return true
}

////////////////////////////////////////////////////////

///////////ajax response and parse
func hello(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "POST" {
		r.ParseForm()

		id := r.Form["id"][0]
		mail := send(id)
		if mail {

			fmt.Fprintln(w, "yeeeaaaaah !")

		} else {
			fmt.Fprintln(w, "not sent")
		}

	} else {
		fmt.Fprintln(w, "not post")
	}

}

////////////////////////////////////////////////////////////////////////////

func main() {

	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", hello)

	log.Printf("Listening on %s...\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}

}
