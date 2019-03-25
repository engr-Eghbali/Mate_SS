package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"github.com/gin-gonic/gin/"
)

///////!!! GLOBAL VARS !!!!!!/////////

var session *mgo.Session

type User struct {
	ID    bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name  string        `json:name`
	Phone string        `json:phone`
	Email string        `json:email`
}

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

//func Redis(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "text/javascript")
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//
//	if r.Method == "POST" {
//
//		r.ParseForm()
//		key := r.Form["key"][0]
//		val := r.Form["val"][0]
//
//		client := redis.NewClient(&redis.Options{
//			Addr:     "redis-13657.c135.eu-central-1-1.ec2.cloud.redislabs.com:13657",
//			Password: "tlqTsgjgzDOqZb2bYjHAMCcC4uh9U49o", // no password set
//			DB:       0,                                  // use default DB
//		})
//
//		_, err := client.Ping().Result()
//
//		if err != nil {
//			fmt.Fprintln(w, err)
//		} else {
//
//			err := client.Set(key, val, 0).Err()
//			if err != nil {
//				panic(err)
//			}
//
//			value, err := client.Get(key).Result()
//			if err != nil {
//				panic(err)
//			}
//			fmt.Fprintln(w, value)
//
//		}
//
//	} else {
//		fmt.Fprintln(w, "not post method")
//	}
//}

/////////////////////////////////////////////////////////////////////////////////////////

///////////////postqresql simple sample

//func MongoTest(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "text/javascript")
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//
//	type user struct {
//		ID    bson.ObjectId `json:"id" bson:"_id,omitempty"`
//		Name  string        `json:"name"`
//		Phone string        `json:"Phone"`
//	}
//
//	var newuser user
//
//	if r.Method == "POST" {
//
//		r.ParseForm()
//		key := r.Form["key"][0]
//		val := r.Form["val"][0]
//
//		newuser.ID = bson.NewObjectId()
//		newuser.Name = key
//		newuser.Phone = val
//
//		session, err := mgo.Dial("mongodb://udlt7amzwwc3lav9copw:QWqGAUmRLERX081CYU4k@bkbfbtpiza46rc3-mongodb.services.clever-cloud.com:27017/bkbfbtpiza46rc3")
//		defer session.Close()
//		session.SetMode(mgo.Monotonic, true)
//
//		if err != nil {
//
//			fmt.Fprintln(w, err)
//
//		} else {
//
//			///**check if any inorder cart ,remove it befor start new one
//			c := session.DB("bkbfbtpiza46rc3").C("users")
//			err = c.Insert(&newuser)
//
//			if err != nil {
//
//				fmt.Fprintln(w, "query failed")
//
//			} else {
//
//				fmt.Fprintln(w, "query done")
//
//			}
//
//		}
//
//	} else {
//		fmt.Fprintln(w, "not post method")
//	}
//}
//
/////////////////////////////////////////////////////////////////////////

///// Generate/save/send verification code to client Email
func SendVerificationMail(mail string) (err bool) {

	//generate
	rand.Seed(time.Now().UnixNano())
	vc := 100000 + rand.Intn(999999-100000)

	//save

}

////////// Handle submit function
func SubmitReq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.ParseForm()
	mORp := r.Form["type"][0]
	data := r.Form["data"][0]

	if mORp == "m" {

		if strings.Contains(data, "@") {

			err := SendVerificationMail(data)

			if err != true {
				fmt.Fprintf(w, "0")
			} else {
				fmt.Fprintf(w, "1")
			}
		}

	}

	///// handle confirm code by SMS
	//	if mORp =="p"{
	//	}

	if mORp != "m" && mORp != "p" {
		fmt.Fprintf(w, "bad request")
	}
}

///////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////

func main() {

	///main vars
	var DBerr error

	/// Call determine listen address
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}

	/// Mongo Dialling
	session, DBerr = mgo.Dial("mongodb://udlt7amzwwc3lav9copw:QWqGAUmRLERX081CYU4k@bkbfbtpiza46rc3-mongodb.services.clever-cloud.com:27017/bkbfbtpiza46rc3")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	if DBerr != nil {
		log.Fatal(DBerr)
	}

	//Routing
	http.HandleFunc("/submit", SubmitReq)

	if Port_err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}

}
