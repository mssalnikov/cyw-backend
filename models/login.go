package models

// Header struct
type LoginHeader struct {
}

type LoginRequest struct {
	FbUid   string `json:"fb_uid"`
	FbToken string `json:"fb_token"`
}

type LoginResponse struct {
	AuthToken string `json:"auth-token"`
}
