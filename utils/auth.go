package utils

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

type FBError struct {
	Message string `json:"message"`
	Type string `json:"type"`
	OAuthException string `json:"OAuthException"`
	Code int `json:"code"`
	Fbtrace_id string `json:"fbtrace_id"`
}

type AccessResponse struct {
	Name string `json:"name"`
	Id string  `json:"id"`
	Email string `json:"email"`
	Error FBError `json:"error,omitempty"`
}
