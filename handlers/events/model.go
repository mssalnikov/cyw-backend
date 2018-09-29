package events

import (
	"github.com/simplereach/timeutils"
	"database/sql"
)

// Event struct
type Event struct {
	Id          uint64         `db:"id" json:"id"`
	UserId      uint64         `db:"user_id" json:"user_id"`
	Name        string         `db:"name" json:"name"`
	Description string         `db:"description" json:"description"`
	Start       timeutils.Time `db:"start" json:"startt"`
	Finish      timeutils.Time `db:"finish" json:"finish"`
}

type PointOfEvent struct {
	Id          uint64        `db:"id" json:"id"`
	EventId     uint64        `db:"event_id" json:"event_id"`
	Container   string        `db:"container" json:"container"`
	Naviaddress string        `db:"naviaddress" json:"naviaddress"`
	Question    string        `db:"question" json:"question"`
	Answer      string        `db:"answer" json:"answer"`
	Token       int64         `db:"token" json:"token"`
	PrevPointId sql.NullInt64 `db:"prev_point_id" json:"prev_point_id"`
}

type UserPoint struct {
	Id       uint64 `db:"id" json:"id"`
	UserId   uint64 `db:"user_id" json:"user_id"`
	PointId  uint64 `db:"point_id" json:"point_id"`
	IsSolved bool   `db:"is_solved" json:"is_solved"`
}

type NaviaddressPost struct {
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	DefaultLang string  `json:"default_lang"`
	AddressType string  `json:"address_type"`
}

//{"result":{"address_type":"free","container":"7","count_favorites":0,"count_likes":0,"count_views":0,"default_lang":"ru","expires_on":"2018-09-29T12:46:47.0804142Z","id":1537159,"is_external":false,"is_html":false,"langs":["ru"],"locale":"ru","map_visibility":true,"naviaddress":"703893","owner_id":17366,"point":{"lat":55.424,"lng":36.3123},"postal_address":"","priority":4}}

type NaviaddressPostResponse struct {
	Result NaviaddressPostResponseResult `json:"result"`
}

type NaviaddressPostResponseResult struct {
	Naviaddress string `json:"naviaddress"`
	Container   string `json:"container"`
	AddressType string `json:"address_type"`
}

type NaviaddressAccept struct {
	Container   string `json:"container"`
	Naviaddress string `json:"naviaddress"`
}

type NaviaddressAcceptResult struct {
	Container   string `json:"container"`
	Naviaddress string `json:"naviaddress"`
	AddressType string `json:"address_type"`
}

type NaviaddressPut struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DefaultLang string  `json:"default_lang"`
	Lang string  `json:"lang"`
}

type AcceptEvent struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type NewEvent struct {
	Points      []EventPoint   `json:"points"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Start       timeutils.Time `json:"start"`
	Finish      timeutils.Time `json:"finish"`
}

type EventPoint struct {
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Question string  `json:"question"`
	Answer   string  `json:"answer"`
}
