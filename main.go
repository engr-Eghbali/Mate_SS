package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	//"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
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
// func send(body string) (er bool) {
//	from := "trane2sfc@gmail.com"
//	pass := "Wakeuptrane2sfc$"
//	to := "engr.eghbali@gmail.com"
//
//	msg := "From: " + from + "\n" +
//		"To: " + to + "\n" +
//		"Subject: Hello there\n\n" +
//		body
//
//	err := smtp.SendMail("smtp.gmail.com:587",
//		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
//		from, []string{to}, []byte(msg))
//
//	if err != nil {
//		log.Printf("smtp error: %s", err)
//		return false
//	}
//
//	log.Print("sent, visit http://foobarbazz.mailinator.com")
//	return true
//}

////////////////////////////////////////////////////////

///////////ajax response and parse
//func hello(w http.ResponseWriter, r *http.Request) {
//
//	w.Header().Set("Content-Type", "text/javascript")
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//
//	if r.Method == "POST" {
//		r.ParseForm()
//
//		id := r.Form["id"][0]
//		mail := send(id)
//		if mail {
//
//			fmt.Fprintln(w, "yeeeaaaaah !")
//
//		} else {
//			fmt.Fprintln(w, "not sent")
//		}
//
//	} else {
//		fmt.Fprintln(w, "not post")
//	}
//
//}

////////////////////////////////////////////////////////////////////////////

//#########################REDIS SAMPLE######################################################
//                                                                           ################
//create pool                                                                ################

func Redis(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "POST" {

		r.ParseForm()
		key := r.Form["key"][0]
		val := r.Form["val"][0]

		client := redis.NewClient(&redis.Options{
			Addr:     "redis-13657.c135.eu-central-1-1.ec2.cloud.redislabs.com:13657",
			Password: "tlqTsgjgzDOqZb2bYjHAMCcC4uh9U49o", // no password set
			DB:       0,                                  // use default DB
		})

		pong, err := client.Ping().Result()

		if err != nil {
			fmt.Fprintln(w, err)
		} else {

			err := client.Set(key, val, 0).Err()
			if err != nil {
				panic(err)
			}

			value, err := client.Get(key).Result()
			if err != nil {
				panic(err)
			}
			fmt.Fprintln(w, value)

		}

	} else {
		fmt.Fprintln(w, "not post method")
	}
}

func main() {

	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", Redis)

	log.Printf("Listening on %s...\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}

}
