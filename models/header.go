package models

// Header struct
type Header struct {
	Status    string      `json:"status"`
	AuthToken string      `json:"auth-token"`
	Count     int         `json:"count"`
	Data      interface{} `json:"data"`
}
