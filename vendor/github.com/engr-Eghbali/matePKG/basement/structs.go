package structs

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type VcTable struct {
	ID     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserID string        `json:userid`
	VC     string        `json:vc`
}

type Location struct {
	X   string `json:x`
	Y   string `json:y`
	Add string `json:add`
}

type Request struct {
	SenderID   bson.ObjectId `json:senderid`
	SenderName string        `json:sendername`
	SenderPic  string        `json:sendername`
}

type Meet struct {
	Title string    `json:title`
	Time  time.Time `json:time`
	Host  string    `json:host`
	Crowd string    `json:crowd`
	Geo   Location  `json:location`
}

type User struct {
	ID         bson.ObjectId   `json:"id" bson:"_id,omitempty"`
	Name       string          `json:name`
	Phone      string          `json:phone`
	Email      string          `json:email`
	Vc         string          `json:vc`
	Status     int8            `json:status`
	Avatar     string          `json:avatar`
	FriendList []bson.ObjectId `json:friendlist`
	Meetings   []Meet          `json:meetings`
	Requests   []Request       `json:request`
}

type MailOrigin struct {
	From     string
	Password string
}

type SmsOrigin struct {
	From   string
	ApiKey string
}
