package users

// User struct
type User struct {
	Id      int64     `db:"id" json:"id"`
	FBUid     string    `db:"fb_uid" json:"fb_uid"`
	UserName  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	//DateAdded time.Time `db:"date_added" json:"date_added"`
}

type GetUser struct {
	Id      int64     `json:"id"`
}
// Header struct
type FbAuth struct {
	FbAccessToken string `json:"fbAccessToken"`
}

type LoginRequest struct {
	FbUid   string `json:"fb_uid"`
	FbToken string `json:"fb_token"`
}

type LoginResponse struct {
	AuthToken string `json:"auth-token"`
}
