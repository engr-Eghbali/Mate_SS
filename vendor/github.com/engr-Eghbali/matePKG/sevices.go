package services

import (
	"encoding/base64"
	"encoding/json"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"

	structs "github.com/engr-Eghbali/matePKG/basement"
	configs "github.com/engr-Eghbali/matePKG/basement/conf"
	"github.com/go-redis/redis"

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
func InitUser(id string, vc string, session *mgo.Session, clint *redis.Client) (objid bson.ObjectId, result bool) {

	var InsertErr error
	collection := session.DB("bkbfbtpiza46rc3").C("users")

	if strings.Contains(id, "@") {
		NewUser := structs.User{ID: bson.NewObjectId(), Name: "", Phone: "", Email: id, Vc: vc, Status: 1, Avatar: "../assets/img/profileAvatar.svg", FriendList: nil, Meetings: nil, Requests: nil}
		objid = NewUser.ID
		InsertErr = collection.Insert(&NewUser)
	} else {
		NewUser := structs.User{ID: bson.NewObjectId(), Name: "", Phone: id, Email: "", Vc: vc, Status: 1, Avatar: "../assets/img/profileAvatar.svg", FriendList: nil, Meetings: nil, Requests: nil}
		objid = NewUser.ID
		InsertErr = collection.Insert(&NewUser)
	}

	if InsertErr != nil {
		log.Println("Init User failed")
		log.Println(InsertErr)
		log.Println("<=End")
		return objid, false
	} else {

		init := structs.UserCache{Geo: "0,0", Vc: vc, FriendList: nil, Visibility: true}
		if !SendToCache(objid.Hex(), init, clint) {
			log.Println("redis init failed,trying again")
			SendToCache(objid.Hex(), init, clint)
		}
		return objid, true
	}

}

////////////////////////////////////////////////////////////////
func LoginUser(userId string, vc string, session *mgo.Session) (res bool) {

	collection := session.DB("bkbfbtpiza46rc3").C("users")
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

//////////////////////////////////////////////////////
/////////////////////////////////////////////////////

////////////////// fetch user info from cache
func CacheRetrieve(client *redis.Client, keys ...string) ([]structs.UserCache, error) {

	var temp structs.UserCache
	var Results []structs.UserCache
	Users, err := client.MGet(keys...).Result()

	if err != nil {
		log.Println("cache retrieve service failur:")
		log.Println(err)
		log.Println("<=End")
		return Results, err
	}

	for _, user := range Users {
		if user != nil {
			json.Unmarshal([]byte(user.(string)), &temp)
			Results = append(Results, temp)
		}
	}
	return Results, err

}

/////////////push user info to the cache
func SendToCache(key string, data structs.UserCache, client *redis.Client) (err bool) {

	info, _ := json.Marshal(data)
	expiration, _ := time.Parse("01/01/2025", "01/01/2025")
	duration := time.Until(expiration)
	stat := client.Set(key, info, duration)
	if stat.Err() == nil {
		return true
	} else {
		log.Println(stat)
		return false
	}

}

/////////////////make special pins for every person

func PinMaker(avatar string) (markerB64 string) {

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(avatar))
	profilePic, _, err := image.Decode(reader)
	ppResized := imaging.Resize(profilePic, 40, 40, imaging.NearestNeighbor)

	randPin := configs.GetPinBase()
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(randPin))
	pin, _, err2 := image.Decode(reader)

	offset := image.Pt(12, 6)
	b := pin.Bounds()
	marker := image.NewRGBA(b)
	draw.Draw(marker, b, pin, image.ZP, draw.Src)
	draw.Draw(marker, ppResized.Bounds().Add(offset), ppResized, image.ZP, draw.Over)

	png.Encode(&buff, marker)
	// Encode the bytes in the buffer to a base64 string
	encodedString := base64.StdEncoding.EncodeToString(buff.Bytes())

	return encodedString

}
