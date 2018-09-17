package events

import (
	"time"
	"github.com/kellydunn/golang-geo"
)

// Event struct
type Event struct {
	Id        int64     `db:"id" json:"id"`
	UserId    int64     `db:"user_id" json:"user_id"`
	Description    string     `db:"description" json:"description"`
	StartAt     time.Time    `db:"start_at" json:"start_at"`
	FinishAt time.Time `db:"finish_at" json:"finish_at"`
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