package events

import (
	"github.com/simplereach/timeutils"
	"database/sql"
	"time"
)

// Event struct
type Event struct {
	Id          uint64         `db:"id" json:"id"`
	UserId      uint64         `db:"user_id" json:"user_id"`
	Name        string         `db:"name" json:"name"`
	Description string         `db:"description" json:"description"`
	Start       timeutils.Time `db:"start" json:"start"`
	Finish      timeutils.Time `db:"finish" json:"finish"`
}

type EventFromDB struct {
	Id          uint64    `db:"id" json:"id"`
	UserId      uint64    `db:"user_id" json:"user_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Start       time.Time `db:"start" json:"start"`
	Finish      time.Time `db:"finish" json:"finish"`
}

type EventFromDBForUser struct {
	Id          uint64        `json:"id"`
	UserId      uint64        `json:"user_id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Start       time.Time     `json:"start"`
	Finish      time.Time     `json:"finish"`
	Points      []LabelPoints `json:"points"`
}

type PointFromDbForUser struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	Question    string `json:"question"`
	Container   string `json:"container"`
	Naviaddress string `json:"naviaddress"`
	IsSolved    bool   `json:"is_solved"`
	IsFound     bool   `json:"is_found"`
}

type PointId struct {
	Id uint64 `db:"id" json:"id"`
}

type LabelPoints struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	Question    string `json:"question"`
	Container   string `json:"question"`
	Naviaddress string `json:"question"`
	IsSolved    bool   `json:"is_solved"`
	IsFound     bool   `json:"is_found"`
}

type PointOfEvent struct {
	Id          uint64        `db:"id" json:"id"`
	EventId     uint64        `db:"event_id" json:"event_id"`
	Container   string        `db:"container" json:"container"`
	Naviaddress string        `db:"naviaddress" json:"naviaddress"`
	Name        string        `db:"name" json:"name"`
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
	IsFound  bool   `db:"is_found" json:"is_found"`
}

type NaviaddressPost struct {
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	DefaultLang string  `json:"default_lang"`
	AddressType string  `json:"address_type"`
}

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
	DefaultLang string `json:"default_lang"`
	Lang        string `json:"lang"`
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
	Name     string  `json:"name"`
	Question string  `json:"question"`
	Answer   string  `json:"answer"`
	Token    string  `json:"token"`
}

type EnterToken struct {
	PointId int64  `json:"id"`
	Token   string `json:"token"`
}

type AnswerQuestion struct {
	PointId int64  `json:"point_id"`
	Answer  string `json:"answer"`
}
