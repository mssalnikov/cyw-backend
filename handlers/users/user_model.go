package users

import (
	"time"
)

// User struct
type User struct {
	ID        int64     `db:"id" json:"id"`
	FBUid     string    `db:"fb_uid" json:"fb_uid"`
	UserName  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	DateAdded time.Time `db:"date_added" json:"date_added"`
}
