package welcome

import (
	"github.com/dghubble/sessions"
	"net/http"
)

const (
	sessionName    = "example-facebook-app"
	sessionSecret  = "example cookie signing secret"
	sessionUserKey = "facebookID"
)

// sessionStore encodes and decodes session data stored in signed cookies
var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

// Config configures the main ServeMux.
type Config struct {
	FacebookClientID     string
	FacebookClientSecret string
}


// Decorator RequireLogin
func RequireLogin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if !isAuthenticated(req) {
			http.Redirect(w, req, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}
