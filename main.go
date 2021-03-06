package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"

	services "github.com/engr-Eghbali/matePKG"

	structs "github.com/engr-Eghbali/matePKG/basement"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"github.com/gin-gonic/gin/"
)

///////!!! GLOBAL VARS !!!!!!/////////

var session *mgo.Session
var redisClient *redis.Client

type Image struct {
	Filename    string
	ContentType string
	Data        []byte
	Size        int
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
////////////////////////////////////////////////////////////////////////

//////////////// update users geometrical and Maps info
func GodsEye(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")
		return
	} else {

		r.ParseForm()
		VC := r.Form["vc"][0]
		ID := r.Form["id"][0] //objId
		Geo := r.Form["geo"][0]
		var Visible bool
		if r.Form["visible"][0] == "0" {
			Visible = false
		} else {
			Visible = true
		}

		user, err := services.CacheRetrieve(redisClient, ID)

		if err != nil || user[0].Vc != VC {
			fmt.Fprintln(w, "-1")
			return
		}

		//update user info
		TempUser := structs.UserCache{Geo: Geo, Vc: user[0].Vc, FriendList: user[0].FriendList, Visibility: Visible}
		flag := services.SendToCache(ID, TempUser, redisClient)
		if flag == false {
			log.Println("update cache failed due to SendToCache failur")
		}

		///iterate friends
		var friendKeys []string
		for _, id := range user[0].FriendList {
			friendKeys = append(friendKeys, id.Hex())
		}
		Friends, ferr := services.CacheRetrieve(redisClient, friendKeys...)
		if ferr != nil {

			log.Println(user)
			fmt.Fprintln(w, "0")
			return
		}

		///form the answer
		type Revelation struct {
			ID  string
			Geo string
		}
		var response []Revelation
		for i, f := range Friends {
			if f.Visibility == true {

				response = append(response, Revelation{ID: friendKeys[i], Geo: f.Geo})

			} else {

				response = append(response, Revelation{ID: friendKeys[i], Geo: "0"})

			}
		}

		b, _ := json.Marshal(response)
		fmt.Fprintln(w, string(b))

	}

}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

//////////////////////////user request to add a bill on the meeting

func AddBill(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")
		return
	}
	r.ParseForm()
	ID := r.Form["id"][0] //objId
	Vc := r.Form["vc"][0]
	Title := r.Form["title"][0]
	Geostr := strings.Split(r.Form["geo"][0], ",")
	Geo := structs.Location{X: Geostr[0], Y: Geostr[1]}
	var Bill structs.Bill
	var user structs.User
	var updateErr error
	err := json.Unmarshal([]byte(r.Form["bill"][0]), &Bill)

	if err != nil {
		fmt.Fprintln(w, "bad format request")
		return
	}

	collection := session.DB("bkbfbtpiza46rc3").C("users")
	findErr := collection.FindId(bson.ObjectIdHex(ID)).One(&user)

	if findErr != nil || user.Vc != Vc {
		fmt.Fprintln(w, "-1")
		return
	}

	for _, meet := range user.Meetings {
		if meet.Title == Title && meet.Geo == Geo {

			updateInvoice := append(meet.Invoice, Bill)
			meet.Crowd = append(meet.Crowd, ID)

			for i, personId := range meet.Crowd {

				query := bson.M{
					"_id":            bson.ObjectIdHex(personId),
					"Meetings.Title": Title,
					"Meetings.Geo":   Geo,
				}
				update := bson.M{
					"$set": bson.M{
						"Meetings.$.Invoice": updateInvoice,
					},
				}

				updateErr = collection.Update(query, update)

				if updateErr != nil {
					log.Println("invoice update failed:")
					log.Println(personId)
					log.Println("trying again...")
					updateErr = collection.Update(query, update)
					if updateErr != nil {
						log.Println("**aborted** <=END")
						log.Println("Rolling Back Transactions...")
						for j := 0; j < i; j++ {

							query = bson.M{
								"_id":            bson.ObjectIdHex(meet.Crowd[j]),
								"Meetings.Title": Title,
								"Meetings.Geo":   Geo,
							}
							update = bson.M{
								"$pop": bson.M{
									"Meetings.$.Invoice": 1,
								},
							}

							updateErr = collection.Update(query, update)
						}

						if updateErr == nil {
							log.Println("DONE...! <=END")
						} else {
							log.Println("**RollBack failur occured**")
						}
						fmt.Fprintln(w, "0")
						return
					} else {
						log.Println("DONE...! <=END")
					}

				}
			}

			fmt.Fprintln(w, "1")
			return
		}
	}

}

/////////////setup a meeting
func SetMeeting(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")
		return
	}
	r.ParseForm()
	ID := r.Form["id"][0] //objId
	Vc := r.Form["vc"][0]
	Title := r.Form["title"][0]
	Time, _ := time.Parse(time.RFC3339, r.Form["time"][0])

	var Crowd []string
	if strings.Contains(r.Form["crowd"][0], ",") {
		Crowd = strings.Split(r.Form["crowd"][0], ",")
	} else {
		Crowd = append(Crowd, string(r.Form["crowd"][0]))
	}

	Geo := structs.Location{X: strings.Split(r.Form["geo"][0], ",")[0], Y: strings.Split(r.Form["geo"][0], ",")[1]}

	var user, temp structs.User
	var updateErr error
	collection := session.DB("bkbfbtpiza46rc3").C("users")
	findErr := collection.FindId(bson.ObjectIdHex(ID)).One(&user)

	if findErr != nil || user.Vc != Vc {
		fmt.Fprintln(w, "-1")
		return
	}

	newMeeting := structs.Meet{Title: Title, Time: Time, Host: user.Name, Crowd: Crowd, Geo: Geo}

	updateErr = collection.UpdateId(bson.ObjectIdHex(ID), bson.M{"$set": bson.M{"meetings": append(user.Meetings, newMeeting)}})
	if updateErr != nil {
		fmt.Fprintln(w, "0")
		log.Println("user set meeting failed: ")
		log.Println(updateErr)
		return
	}

	if Crowd[0] == "" {
		b, _ := json.Marshal(append(user.Meetings, newMeeting))
		fmt.Fprintln(w, string(b))
		return
	}

	for _, personID := range Crowd {

		log.Println(personID)
		findErr = collection.FindId(bson.ObjectIdHex(personID)).One(&temp)

		if findErr == nil {
			for _, id := range temp.FriendList {
				if id == bson.ObjectIdHex(ID) {
					updateErr = collection.UpdateId(temp.ID, bson.M{"$set": bson.M{"meetings": append(temp.Meetings, newMeeting)}})
					if updateErr != nil {
						log.Println("invite member to meeting failed:")
						log.Println(updateErr)
						log.Println("trying again:")
						updateErr = collection.UpdateId(temp.ID, bson.M{"$set": bson.M{"meetings": append(temp.Meetings, newMeeting)}})
						log.Println(updateErr)
					}
					break
				}
			}
		}

	}

	b, _ := json.Marshal(append(user.Meetings, newMeeting))

	fmt.Fprintln(w, string(b))

}

///////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////

////////////leave a meeting
func LeaveMeeting(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")
		return
	}

	r.ParseForm()
	ID := r.Form["id"][0] //objId
	Vc := r.Form["vc"][0]
	Title := r.Form["title"][0]
	geo := structs.Location{X: strings.Split(r.Form["geo"][0], ",")[0], Y: strings.Split(r.Form["geo"][0], ",")[1]}

	var user structs.User
	var newMeetingList []structs.Meet
	collection := session.DB("bkbfbtpiza46rc3").C("users")
	findErr := collection.FindId(bson.ObjectIdHex(ID)).One(&user)

	if findErr != nil || user.Vc != Vc {
		fmt.Fprintln(w, "-1")
		return
	}

	for i, meet := range user.Meetings {
		if meet.Title == Title && meet.Geo.X == geo.X && meet.Geo.Y == geo.Y {
			newMeetingList = append(user.Meetings[:i], user.Meetings[i+1:]...)
			break
		}
	}

	updateErr := collection.UpdateId(bson.ObjectIdHex(ID), bson.M{"$set": bson.M{"meetings": newMeetingList}})
	if updateErr != nil {
		fmt.Fprintln(w, "0")
		log.Println("user leave meeting method failed:")
		log.Println(updateErr)
		return
	}

	b, Merr := json.Marshal(newMeetingList)

	if Merr != nil {
		log.Println(Merr)
	}

	fmt.Fprintln(w, string(b))

}

///////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////

////handle and deliver friend request to another user
func SendFriendReq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")
		return
	}

	r.ParseForm()
	ID := r.Form["id"][0] //objId
	Vc := r.Form["vc"][0]
	Frequest := r.Form["frequest"][0]
	var friend, user structs.User

	collection := session.DB("bkbfbtpiza46rc3").C("users")
	findErr := collection.Find(bson.M{"name": Frequest}).One(&friend)

	if findErr == mgo.ErrNotFound {
		fmt.Fprintln(w, "0")
		return
	}

	findErr = collection.FindId(bson.ObjectIdHex(ID)).One(&user)
	if findErr != nil { // +user.Name == Frequest
		fmt.Fprintln(w, "0")
		log.Println("friend request failur due to query error:")
		log.Println(findErr)
		log.Println("<=End")
		return
	}

	if user.Vc != Vc {
		fmt.Fprintln(w, "-1")
		return
	}

	for _, r := range friend.Requests {
		if r.SenderName == user.Name {
			fmt.Fprintln(w, "2")
			return
		}
	}

	Reqlist := append(friend.Requests, structs.Request{SenderName: user.Name, SenderPic: user.Avatar})
	updateErr := collection.UpdateId(friend.ID, bson.M{"$set": bson.M{"requests": Reqlist}})

	if updateErr != nil {
		fmt.Fprintln(w, "0")
		log.Println("user F-Request failed due to query update error:")
		log.Println(updateErr)
		log.Println("<=End")
		return
	}

	fmt.Fprintln(w, "1")

}

////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////

////handle accepting friend request and make both friends
func AcceptFrequest(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")
		return
	}

	r.ParseForm()
	ID := r.Form["id"][0] //objId
	VC := r.Form["vc"][0]
	Fname := r.Form["fname"][0]

	var user, friend structs.User
	collection := session.DB("bkbfbtpiza46rc3").C("users")

	Err1 := collection.FindId(bson.ObjectIdHex(ID)).One(&user)
	Err2 := collection.Find(bson.M{"name": Fname}).One(&friend)

	if Err1 != nil || Err2 != nil || user.Vc != VC {
		fmt.Fprintln(w, "-1")
		return
	}

	var ReqList []structs.Request
	var Flist []bson.ObjectId
	for i, r := range user.Requests {
		if r.SenderName == friend.Name {
			ReqList = append(user.Requests[:i], user.Requests[i+1:]...)
			break
		}
	}
	Flist = append(user.FriendList, friend.ID)
	Err1 = collection.UpdateId(user.ID, bson.M{"$set": bson.M{"friendlist": Flist, "requests": ReqList}})
	if Err1 != nil {
		fmt.Fprintln(w, "0")
		log.Println("accepting request failed due to update query error:")
		log.Println(Err1)
		log.Println("<=End")
		return
	}

	Flist = append(friend.FriendList, user.ID)
	Err2 = collection.UpdateId(friend.ID, bson.M{"$set": bson.M{"friendlist": Flist}})

	if Err2 != nil {
		fmt.Fprintln(w, "0")
		log.Println("accepting request failed due to update query error:")
		log.Println(Err1)
		log.Println("<=End")
		collection.UpdateId(user.ID, bson.M{"$set": bson.M{"friendlist": user.FriendList, "requests": user.Requests}})
		return
	}

	/////then update redis cache
	var userCache, friendCache structs.UserCache
	userTMP, cacheErr1 := services.CacheRetrieve(redisClient, user.ID.Hex())
	friendTMP, cacheErr2 := services.CacheRetrieve(redisClient, friend.ID.Hex())

	if cacheErr1 != nil || cacheErr2 != nil {
		fmt.Fprintln(w, "0")
		log.Println("user accept frequest set failed due to cache retrieve service error:")
		log.Println(cacheErr1)
		log.Println(cacheErr2)
		log.Println("<=End")
		return
	}

	userCache = structs.UserCache{Geo: userTMP[0].Geo, Vc: userTMP[0].Vc, FriendList: append(userTMP[0].FriendList, friend.ID), Visibility: userTMP[0].Visibility}
	friendCache = structs.UserCache{Geo: friendTMP[0].Geo, Vc: friendTMP[0].Vc, FriendList: append(friendTMP[0].FriendList, user.ID), Visibility: friendTMP[0].Visibility}

	setCacheErr1 := services.SendToCache(user.ID.Hex(), userCache, redisClient)
	setCacheErr2 := services.SendToCache(friend.ID.Hex(), friendCache, redisClient)

	if setCacheErr1 != true || setCacheErr2 != true {
		fmt.Fprintln(w, "0")
		log.Println("user accept f-request set failed due to cache set service error:")
		log.Println(setCacheErr1)
		log.Println(setCacheErr2)
		log.Println("<=End")
		return
	}
	fmt.Fprintln(w, "1")

}

///////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////

//////////unfriend function
func Unfriend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")
		return
	}
	r.ParseForm()
	ID := r.Form["id"][0]
	VC := r.Form["vc"][0]
	Target := r.Form["target"][0] //username for unfriend

	var user, friend structs.User
	collection := session.DB("bkbfbtpiza46rc3").C("users")

	findErr1 := collection.FindId(bson.ObjectIdHex(ID)).One(&user)
	findErr2 := collection.Find(bson.M{"name": Target}).One(&friend)

	if findErr1 != nil || findErr2 != nil || user.Vc != VC {
		fmt.Fprintln(w, "-1")
		return
	}

	var NewFriendList1, NewFriendList2 []bson.ObjectId
	for i, f := range user.FriendList {
		if f == friend.ID {
			NewFriendList1 = append(user.FriendList[:i], user.FriendList[i+1:]...)
			break
		}
	}
	for i, f := range friend.FriendList {
		if f == user.ID {
			NewFriendList2 = append(friend.FriendList[:i], friend.FriendList[i+1:]...)
			break
		}
	}

	updateErr1 := collection.UpdateId(user.ID, bson.M{"$set": bson.M{"friendlist": NewFriendList1}})
	updateErr2 := collection.UpdateId(friend.ID, bson.M{"$set": bson.M{"friendlist": NewFriendList2}})

	if updateErr1 != nil || updateErr2 != nil {
		fmt.Fprintln(w, "0")
		log.Println("user unfriend failed due to update query failur:")
		log.Println(updateErr1)
		log.Println(updateErr2)
		log.Println("<=End")
		return
	}

	fmt.Fprintln(w, "1")

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////

///////////change avatar request handler
func AvatarChange(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")
		return
	}

	r.ParseMultipartForm(0)
	defer r.MultipartForm.RemoveAll()
	avatar, info, Ferr := r.FormFile("avatar")
	VC := r.FormValue("vc")
	ID := r.FormValue("id")

	contentType := info.Header.Get("Content-Type")
	log.Println(info)
	log.Println(Ferr)

	var temp = new(structs.User)
	collection := session.DB("bkbfbtpiza46rc3").C("users")

	FindErr := collection.FindId(bson.ObjectIdHex(ID)).One(&temp)

	if FindErr != nil || temp.Vc != VC || Ferr != nil /*|| contentType != "image/png"*/ {
		fmt.Fprintln(w, "-1")
		return
	}

	//
	///////////read image
	imageBytes, err := ioutil.ReadAll(avatar)

	if err != nil {
		fmt.Fprintln(w, "0")
		log.Println("reading image file failed:")
		log.Println(err)
		log.Println("<=END")
		return
	}

	_, _, err = image.Decode(bytes.NewReader(imageBytes))

	// ***** DO SOMTHING ABOUT
	if err != nil {
		//fmt.Fprintln(w, "0")
		log.Println("decode image file failed:")
		log.Println(err)
		log.Println("<=END")
	}

	i := Image{
		Filename:    info.Filename,
		ContentType: contentType,
		Data:        imageBytes,
		Size:        len(imageBytes),
	}
	///
	/////convert to base64
	b64Avatar := fmt.Sprintf("data:%s;base64,%s", i.ContentType, base64.StdEncoding.EncodeToString(i.Data))

	////////////////////////////////
	////////////////////////////////
	updateErr := collection.UpdateId(temp.ID, bson.M{"$set": bson.M{"avatar": b64Avatar}})

	if updateErr != nil {
		fmt.Fprintln(w, 0)
		log.Println("user change avatar failed:")
		log.Println(updateErr)
		log.Println("<=END")
		return
	}
	fmt.Fprintln(w, "1")
	return

}

/////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////

///////////username submition/changing handling
func UserNameChange(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")
		return
	}

	r.ParseForm()
	VC := r.Form["vc"][0]
	ID := r.Form["id"][0]
	UserName := r.Form["username"][0]

	var temp = new(structs.User)
	collection := session.DB("bkbfbtpiza46rc3").C("users")

	FindErr := collection.Find(bson.M{"name": UserName}).One(&temp)

	if FindErr == nil {
		fmt.Fprintln(w, "reserved")
		return
	}

	if FindErr == mgo.ErrNotFound {

		FindErr = collection.FindId(bson.ObjectIdHex(ID)).One(&temp)

		if temp.Vc == VC {

			UpdateErr := collection.UpdateId(temp.ID, bson.M{"$set": bson.M{"name": UserName}})

			if UpdateErr != nil {
				fmt.Fprintln(w, "0")
				return
			} else {
				fmt.Fprintln(w, "1")
				return
			}

		} else {
			fmt.Fprintln(w, "-1")
			return
		}

	}

}

////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////

////////////////change user mail on request
func MailChange(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")
		return
	}

	r.ParseForm()
	VC := r.Form["vc"][0]
	ID := r.Form["id"][0]
	Mail := r.Form["mail"][0]

	var temp = new(structs.User)
	collection := session.DB("bkbfbtpiza46rc3").C("users")

	FindErr := collection.Find(bson.M{"email": Mail}).One(&temp)

	if FindErr == nil {
		fmt.Fprintln(w, "reserved")
		return
	}

	if FindErr == mgo.ErrNotFound {

		FindErr = collection.FindId(bson.ObjectIdHex(ID)).One(&temp)

		if temp.Vc == VC && len(temp.Phone) > 7 {

			UpdateErr := collection.UpdateId(temp.ID, bson.M{"$set": bson.M{"email": Mail}})

			if UpdateErr != nil {
				fmt.Fprintln(w, "0")
				return
			} else {
				fmt.Fprintln(w, "1")
				return
			}

		} else {
			fmt.Fprintln(w, "-1")
			return
		}

	}

}

///// Generate/save/send verification code to client Email
func SendVerificationMail(mail string) (err bool) {

	//generate/save to vc table
	vc, result := services.CreateVcRecord(mail, session)
	if result == false {
		log.Println("sending verification failed cause of VC service failur")
		return false
	}

	//send
	origin := structs.MailOrigin{From: "whereismymate.app@gmail.com", Password: "Wakeuptrane2sfc$"}
	MailErr := services.SendMail(vc, mail, origin)
	if MailErr != true {
		log.Println("User Submition Failed Cause Of SMTP Error:002")
		log.Println(MailErr)
		log.Println("End <=002")
		return false
	}

	return true

}

///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////

//////generate/save/send verification code to client Phone No
func SendVerificationSMS(phone string) bool {

	//generate/save to vc table
	vc, result := services.CreateVcRecord(phone, session)
	if result == false {
		log.Println("sending verification failed cause of VC service failur")
		return false
	}

	//send
	origin := structs.SmsOrigin{From: "10001398", ApiKey: "ED09D0D7-5FBA-43A2-8B9D-F0AE79666B52"}
	SmsErr := services.SendSms(vc, phone, origin)
	if SmsErr != true {
		log.Println("User Submition Failed Cause Of SMS service Error:008")
		log.Println(SmsErr)
		log.Println("End <=008")
		return false
	}

	return true

}

///////verify user by verification code and then login or init user methods will call
func UserVerify(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")

	} else {

		r.ParseForm()
		data := r.Form["ID"][0] //mail or phone
		vc := r.Form["vc"][0]
		collection := session.DB("bkbfbtpiza46rc3").C("loginRequests")
		recordTemp := new(structs.VcTable)
		var FindErr error
		var result bool
		var objid bson.ObjectId

		FindErr = collection.Find(bson.M{"userid": data}).One(&recordTemp)

		if FindErr == nil {
			if recordTemp.VC == vc {

				collection = session.DB("bkbfbtpiza46rc3").C("users")
				var usrTemp structs.User

				if strings.Contains(data, "@") {
					FindErr = collection.Find(bson.M{"email": data}).One(&usrTemp)
				} else {
					FindErr = collection.Find(bson.M{"phone": data}).One(&usrTemp)
				}

				///if user doesnt exist then init new one
				if FindErr == mgo.ErrNotFound {
					objid, result = services.InitUser(data, vc, session, redisClient)
				}

				//if exist then login it
				if FindErr == nil {
					result = services.LoginUser(data, vc, session)
					objid = usrTemp.ID
					init := structs.UserCache{Geo: "0,0", Vc: vc, FriendList: usrTemp.FriendList, Visibility: true}
					if !services.SendToCache(objid.Hex(), init, redisClient) {
						log.Println("redis init failed,trying again")
						services.SendToCache(objid.Hex(), init, redisClient)
					}

				}

				// if error
				if FindErr != nil && FindErr != mgo.ErrNotFound {
					log.Println("user verification failed due to DB query failur")
					log.Println(FindErr)
					log.Println("user ID:")
					log.Println(data)
					log.Println("<=End")
					fmt.Fprintln(w, "0")
				}

				if result == true {
					fmt.Fprintln(w, objid.Hex()+"<>"+usrTemp.Name+"<>"+usrTemp.Avatar)
				} else {
					log.Println("user verification failed due to inituser/loginuser service failur:")
					log.Println(data)
					log.Println("<=End")
					fmt.Fprintln(w, "0")
				}

			} else {
				fmt.Fprintln(w, "-1")
			}
		} else {
			log.Println("user verification failed due to VC table failur")
			log.Println(FindErr)
			log.Println("<=End")
			fmt.Fprintln(w, "0")
		}

	}

}

/////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

//// call verification methods first (+ip limitation must add later+)
func Authenticator(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")
		return
	}

	r.ParseForm()
	data := r.Form["data"][0]

	if strings.Contains(data, "@") {

		err := SendVerificationMail(data)
		if err != true {
			fmt.Fprintf(w, "0")
			return
		} else {
			fmt.Fprintf(w, "1")
			return
		}

	} else {
		err := SendVerificationSMS(data)
		if err != true {
			fmt.Fprintf(w, "0")
			return
		} else {
			fmt.Fprintf(w, "1")
			return
		}
	}

}

/////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////
//// hand shaking func
func HandShake(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")
		return
	}

	r.ParseForm()
	ID := r.Form["id"][0]
	VC := r.Form["vc"][0]

	userCache, err := services.CacheRetrieve(redisClient, ID)

	if err != nil || userCache == nil {
		log.Println("user handashaking failed:")
		log.Println(err)
		log.Println("<=END")
		fmt.Fprintln(w, "0")
		return
	}

	if userCache[0].Vc == VC {
		fmt.Fprintln(w, "1")
		return
	} else {
		fmt.Fprintln(w, "-1")
		return
	}

}

/////////////////////////////////////////////////
////////////////////////////////////////////////////

////client asks for meeting list

func RetrieveMeetings(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")
		return
	}

	r.ParseForm()
	ID := r.Form["id"][0]
	VC := r.Form["vc"][0]
	var user structs.User
	var temp structs.User
	var crowdsName []string

	collection := session.DB("bkbfbtpiza46rc3").C("users")

	findErr := collection.FindId(bson.ObjectIdHex(ID)).One(&user)

	if findErr == mgo.ErrNotFound || VC != user.Vc {
		fmt.Fprintln(w, "-1")
		return
	}

	if findErr != nil {
		fmt.Fprintln(w, "0")
		return
	}

	for i, meet := range user.Meetings {

		if len(meet.Crowd) > 0 {

			for _, person := range meet.Crowd {

				if person != "" {
					findErr = collection.FindId(bson.ObjectIdHex(person)).One(&temp)
					if findErr == nil {
						crowdsName = append(crowdsName, temp.Name)
					}
				}
			}

		}

		user.Meetings[i].Crowd = crowdsName
		crowdsName = nil
	}

	b, _ := json.Marshal(user.Meetings)
	resp := string(b)

	fmt.Fprintln(w, resp)
	return

}

////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////

////////// client asks for his friendList
func RetrieveFriends(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")
		return
	}

	type Profile struct {
		ID     string
		Name   string
		Avatar string
	}
	r.ParseForm()
	ID := r.Form["id"][0]
	VC := r.Form["vc"][0]
	var user, friend structs.User

	var friendsProfile []Profile

	collection := session.DB("bkbfbtpiza46rc3").C("users")

	findErr := collection.FindId(bson.ObjectIdHex(ID)).One(&user)

	if findErr == mgo.ErrNotFound || VC != user.Vc {
		fmt.Fprintln(w, "-1")
		return
	}

	if findErr != nil {
		fmt.Fprintln(w, "0")
		return
	}

	for _, f := range user.FriendList {

		Err := collection.FindId(f).One(&friend)

		if Err == nil {
			friendsProfile = append(friendsProfile, Profile{ID: friend.ID.Hex(), Name: friend.Name, Avatar: friend.Avatar})
		}

	}

	b, _ := json.Marshal(friendsProfile)
	resp := string(b)

	fmt.Fprintln(w, resp)
	return

}

/////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////

///check if user exist(by username) then return avatar+name else 0
func WhoisUser(w http.ResponseWriter, r *http.Request) {

	///**** append rate limitation func in this function is necesserlly ****/

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")
		return
	}

	type Profile struct {
		Name   string
		Avatar string
	}
	r.ParseForm()
	ID := r.Form["id"][0]
	VC := r.Form["vc"][0]
	Query := r.Form["query"][0]

	var friend structs.User
	var queryResult Profile

	if len(ID) < 20 || len(VC) < 6 || len(Query) < 3 {
		fmt.Fprintln(w, "-1")
		return
	}

	collection := session.DB("bkbfbtpiza46rc3").C("users")

	findErr := collection.Find(bson.M{"name": Query}).One(&friend)

	if findErr == mgo.ErrNotFound {
		fmt.Fprintln(w, "0")
		return
	}

	if findErr != nil {
		fmt.Fprintln(w, "0")
		return
	}

	queryResult.Avatar = friend.Avatar
	queryResult.Name = friend.Name

	b, _ := json.Marshal(queryResult)
	fmt.Fprintln(w, string(b))
	return

}

/////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////

func RetrievePendigReqs(w http.ResponseWriter, r *http.Request) {

	///**** append rate limitation func in this function is necesserlly ****/

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")
		return
	}

	type Profile struct {
		Name   string
		Avatar string
	}
	r.ParseForm()
	ID := r.Form["id"][0]
	VC := r.Form["vc"][0]

	var user structs.User

	collection := session.DB("bkbfbtpiza46rc3").C("users")

	findErr := collection.FindId(bson.ObjectIdHex(ID)).One(&user)

	if findErr != nil || VC != user.Vc {
		fmt.Fprintln(w, "0")
		return
	}

	b, _ := json.Marshal(user.Requests)
	fmt.Fprintln(w, string(b))
	return

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/////////deny a recieved frequest
func DenyFrequest(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")
		return
	}

	r.ParseForm()
	ID := r.Form["id"][0]
	VC := r.Form["vc"][0]
	Target := r.Form["target"][0]
	var user structs.User

	collection := session.DB("bkbfbtpiza46rc3").C("users")

	findErr := collection.FindId(bson.ObjectIdHex(ID)).One(&user)

	if findErr != nil || VC != user.Vc {
		fmt.Fprintln(w, "-1")
		return
	}

	for i, req := range user.Requests {
		if req.SenderName == Target {
			user.Requests = append(user.Requests[:i], user.Requests[i+1:]...)
			updateErr := collection.UpdateId(user.ID, bson.M{"$set": bson.M{"request": user.Requests}})
			if updateErr != nil {
				fmt.Fprintln(w, "0")
				return
			} else {
				fmt.Fprintln(w, "1")
				return
			}
		}
	}
	fmt.Fprintln(w, "0")

}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/////////make pin table for client
func PinMap(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		fmt.Fprintln(w, "bad request")
		return
	}

	r.ParseForm()
	ID := r.Form["id"][0]
	VC := r.Form["vc"][0]

	var user structs.User
	var pinMap []structs.PinMap
	var fTemp structs.User

	collection := session.DB("bkbfbtpiza46rc3").C("users")

	findErr := collection.FindId(bson.ObjectIdHex(ID)).One(&user)
	if findErr != nil || VC != user.Vc {
		fmt.Fprintln(w, "-1")
		return
	}

	for _, fid := range user.FriendList {

		findErr = collection.FindId(fid).One(&fTemp)

		if findErr != nil {
			pinMap = append(pinMap, structs.PinMap{ID: fid.Hex(), Pin: "0"})
		} else {

			marker := services.PinMaker(strings.Replace(fTemp.Avatar, "data:image/png;base64,", "", 1))
			pinMap = append(pinMap, structs.PinMap{ID: fid.Hex(), Pin: marker})
		}

	}

	b, _ := json.Marshal(pinMap)
	fmt.Fprintln(w, string(b))

}

////////////////////////////////////////////////
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
	//redis client init
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis-13657.c135.eu-central-1-1.ec2.cloud.redislabs.com:13657",
		Password: "tlqTsgjgzDOqZb2bYjHAMCcC4uh9U49o", // no password set
		DB:       0,                                  // use default DB
	})

	//Routing
	http.HandleFunc("/Auth", Authenticator)
	http.HandleFunc("/Verify", UserVerify)
	http.HandleFunc("/UserName", UserNameChange)
	http.HandleFunc("/Email", MailChange)
	http.HandleFunc("/Avatar", AvatarChange)
	http.HandleFunc("/EyeOfProvidence", GodsEye)
	http.HandleFunc("/Frequest", SendFriendReq)
	http.HandleFunc("/AccFrequest", AcceptFrequest)
	http.HandleFunc("/DenyFrequest", DenyFrequest)
	http.HandleFunc("/Unfriend", Unfriend)
	http.HandleFunc("/SetMeeting", SetMeeting)
	http.HandleFunc("/LeaveMeeting", LeaveMeeting)
	http.HandleFunc("/HandShake", HandShake)
	http.HandleFunc("/ReqMeetingList", RetrieveMeetings)
	http.HandleFunc("/ReqFriendList", RetrieveFriends)
	http.HandleFunc("/Whois", WhoisUser)
	http.HandleFunc("/ReqPendingReqs", RetrievePendigReqs)
	http.HandleFunc("/Bill", AddBill)
	http.HandleFunc("/PinMap", PinMap)

	if Porterr := http.ListenAndServe(addr, nil); Porterr != nil {
		panic(err)
	}

}
