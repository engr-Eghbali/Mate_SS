package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	services "pkg"
	structs "pkg/basement"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"github.com/gin-gonic/gin/"
)

///////!!! GLOBAL VARS !!!!!!/////////

var session *mgo.Session

//////////////port detemine
func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

///////////////////////////////////////////////////////////////////

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
	vc := strconv.Itoa(100000 + rand.Intn(999999-100000))

	//save
	NewUser := structs.User{ID: bson.NewObjectId(), Name: "none", Phone: "none", Email: mail, Vc: vc, Status: 0}
	collection := session.DB("bkbfbtpiza46rc3").C("users")
	InsertErr := collection.Insert(&NewUser)
	if InsertErr != nil {
		log.Println("=>User Submition Failed Cause Of DB Insert Error:001")
		log.Println(InsertErr)
		log.Println("End <=001")
		return false
	}

	//send
	MailErr := services.SendMail(vc, mail)
	if MailErr != true {
		log.Println("User Submition Failed Cause Of SMTP Error:002")
		log.Println(MailErr)
		log.Println("End <=002")
		return false
	}

	return true

}

////////// Handle submit function
func SubmitReq(w http.ResponseWriter, data string) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if strings.Contains(data, "@") {

		err := SendVerificationMail(data)

		if err != true {
			fmt.Fprintf(w, "0")
			return
		} else {
			fmt.Fprintf(w, "1")
			return
		}
	}

	///// handle confirm code by SMS
	//	if mORp =="p"{
	//	}

}

///////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////

///////verify user by verification code
func UserVerify(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		r.ParseForm()
		data := r.Form["ID"][0]
		vc := r.Form["vc"][0]
		collection := session.DB("bkbfbtpiza46rc3").C("users")
		temp := new(structs.User)
		var FindErr error

		if strings.Contains(data, "@") {
			FindErr = collection.Find(bson.M{"email": data}).One(&temp)

		} else {
			FindErr = collection.Find(bson.M{"phone": data}).One(&temp)

		}

		if FindErr == mgo.ErrNotFound {
			fmt.Fprintln(w, "-1")
			log.Println("VC Canceled by Not Found Error due to bad request maybe:004")
			log.Println(FindErr)
			log.Println("<=End004")
			return
		}

		if temp.Vc == vc {
			err := collection.Update(bson.M{"_id": temp.ID}, bson.M{"$set": bson.M{"status": 1}})

			if err != nil {
				fmt.Fprintln(w, "0")
				log.Println("VC Canceled due to DB error:005")
				log.Println(err)
				log.Println("<=End005")
			} else {
				fmt.Fprintln(w, "1")
			}

		} else {
			fmt.Fprintln(w, "0")
		}

	} else {
		fmt.Fprintln(w, "bad request")
	}

}

/////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

//// login or submit? make sure....
func Authenticator(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")
		return
	}

	r.ParseForm()
	data := r.Form["data"][0]
	collection := session.DB("bkbfbtpiza46rc3").C("users")
	temp := new(structs.User)
	var FindErr error

	if strings.Contains(data, "@") {
		FindErr = collection.Find(bson.M{"email": data}).One(&temp)
	} else {
		FindErr = collection.Find(bson.M{"phone": data}).One(&temp)
	}

	if FindErr == mgo.ErrNotFound {
		SubmitReq(w, data)
		return
	}

	if FindErr != nil {

		log.Println("=>User Submition/Login Canceled Cause of DB Find Query Err:003")
		log.Println(FindErr)
		log.Println("End<=003")
		fmt.Fprintln(w, "0")
		return
	}

	//if temp.Status == 0 { LoginUser() }
	//if temp.Status == 1 { LogedinUserHandler()}

}

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
	http.HandleFunc("/Auth", Authenticator)
	http.HandleFunc("Verify", UserVerify)
	if Porterr := http.ListenAndServe(addr, nil); Porterr != nil {
		panic(err)
	}

}
