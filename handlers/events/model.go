package events

import (
	"github.com/simplereach/timeutils"
)

// Event struct
type Event struct {
	Id          uint64         `db:"id" json:"id"`
	UserId      uint64         `db:"user_id" json:"user_id"`
	Name        string         `db:"name" json:"name"`
	Description string         `db:"description" json:"description"`
	Start       timeutils.Time `db:"start" json:"startt"`
	Finish      timeutils.Time `db:"finish" json:"finish"`
	IsPrivate   bool           `db:"is_private" json:"is_private"`
}

type PointOfEvent struct {
	Id          uint64 `db:"id" json:"id"`
	EventId     uint64 `db:"event_id" json:"event_id"`
	Naviaddress string `db:"naviaddress" json:"naviaddress"`
	Question    string `db:"question" json:"question"`
	Answer      string `db:"answer" json:"answer"`
	Token       int64  `db:"token" json:"token"`
	PrevPointId int64  `db:"prev_point_id" json:"prev_point_id"`
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

type NaviaddressAccept struct {
	Container   string `json:"container"`
	Naviaddress string `json:"naviaddress"`
}

type NaviaddressPut struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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
	IsPrivate   bool           `json:"is_private"`
}

type EventPoint struct {
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Question string  `json:"question"`
	Answer   string  `json:"answer"`
}
