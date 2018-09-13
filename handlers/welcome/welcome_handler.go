package welcome

import (
	"net/http"
	"io/ioutil"
	"fmt"
)

func getDefaultHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

// welcomeHandler shows a welcome message and login button.
func (wh *WelcomeHandler) WelcomePage(w http.ResponseWriter, r *http.Request) {
	//getDefaultHeader(w)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if isAuthenticated(r) {
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}
	page, _ := ioutil.ReadFile("home.html")
	fmt.Fprintf(w, string(page))
}

// isAuthenticated returns true if the user has a signed session cookie.
func isAuthenticated(req *http.Request) bool {
	if _, err := sessionStore.Get(req, sessionName); err == nil {
		return true
	}
	return false
}
