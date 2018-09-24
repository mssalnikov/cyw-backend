package events

import (
	"time"
	"github.com/kellydunn/golang-geo"
)

// Event struct
type Event struct {
	Id        int64     `db:"id" json:"id"`
	UserId    int64     `db:"user_id" json:"user_id"`
	Name      string     `db:"name" json:"name"`
	Description    string     `db:"description" json:"description"`
	Naviaddress    string     `db:"naviaddress" json:"naviaddress"`
	StartAt     time.Time    `db:"start" json:"start_at"`
	FinishAt time.Time `db:"finish" json:"finish_at"`
	IsPrivate bool `db:"is_private" json:"is_private"`
}

type PointOfEvent struct {
	Id        int64     `db:"id" json:"id"`
	EventId        int64     `db:"event_id" json:"event_id"`
	Point     geo.Point    `db:"point" json:"point"`
	HaveQuestion bool `db:"have_question" json:"have_question"`
	Question string `db:"question" json:"question"`
	Answer string `db:"answer" json:"answer"`
	IsChained bool `db:"is_chained" json:"is_chained"`
	PrevPointId int64 `db:"prev_point_id" json:"prev_point_id"`
}

type UserPoint struct {
	Id        int64     `db:"id" json:"id"`
	UserId     int64    `db:"user_id" json:"user_id"`
	PointId        int64     `db:"point_id" json:"point_id"`
	IsSolved bool `db:"is_solved" json:"is_solved"`
}
