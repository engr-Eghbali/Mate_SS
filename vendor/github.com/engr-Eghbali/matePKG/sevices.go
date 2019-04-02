package services

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	structs "github.com/engr-Eghbali/matePKG/basement"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

////////////send mail
func SendMail(body string, recipient string, origin structs.MailOrigin) (er bool) {

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
func SendSms(txt string, recipient string, origin structs.SmsOrigin) bool {

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

///////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

//////////// generate VC and save in vc table
func CreateVcRecord(UID string, session *mgo.Session) (vcode string, stat bool) {

	//generate
	rand.Seed(time.Now().UnixNano())
	vc := strconv.Itoa(100000 + rand.Intn(999999-100000))

	VcRecord := structs.VcTable{ID: bson.NewObjectId(), UserID: UID, VC: vc}
	collection := session.DB("bkbfbtpiza46rc3").C("loginRequests")
	collection.Remove(bson.M{"userid": UID}) //remove previous VC
	InsertErr := collection.Insert(&VcRecord)

	if InsertErr != nil {
		log.Println("Creating vc record failed:")
		log.Println(InsertErr)
		log.Println("<=End")
		return "", false
	} else {
		return vc, true
	}

}

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

////////initial user record
func InitUser(id string, vc string, session *mgo.Session) (objid bson.ObjectId, result bool) {

	var InsertErr error
	collection := session.DB("bkbfbtpiza46rc3").C("users")

	if strings.Contains(id, "@") {
		NewUser := structs.User{ID: bson.NewObjectId(), Name: "", Phone: "", Email: id, Vc: vc, Status: 1, Avatar: "pic url here", FriendList: nil, Meetings: nil, Requests: nil}
		objid = NewUser.ID
		InsertErr = collection.Insert(&NewUser)
	} else {
		NewUser := structs.User{ID: bson.NewObjectId(), Name: "", Phone: id, Email: "", Vc: vc, Status: 1, Avatar: "pic url here", FriendList: nil, Meetings: nil, Requests: nil}
		objid = NewUser.ID
		InsertErr = collection.Insert(&NewUser)
	}

	if InsertErr != nil {
		log.Println("Init User failed")
		log.Println(InsertErr)
		log.Println("<=End")
		return objid, false
	} else {
		return objid, true
	}

}

////////////////////////////////////////////////////////////////
func LoginUser(userId string, vc string, session *mgo.Session) (res bool) {

	collection := session.DB("bkbfbtpiza46rc3").C("users")
	var tempUser structs.User
	var UpdateErr error

	if strings.Contains(userId, "@") {

		UpdateErr = collection.Update(bson.M{"email": userId}, bson.M{"$set": bson.M{"vc": vc, "status": 1}})
	} else {

		UpdateErr = collection.Update(bson.M{"phone": userId}, bson.M{"$set": bson.M{"vc": vc, "status": 1}})
	}

	if UpdateErr != nil {
		log.Println("login user failed due to update service failur:")
		log.Println(UpdateErr)
		log.Println("user: " + userId)
		log.Println("<=End")
		return false
	} else {
		return true
	}

}
