package utils

import (
	"net/http"
	"github.com/dghubble/gologin/facebook"
	"log"
	"io/ioutil"
	"fmt"
	"errors"
	"github.com/dghubble/sessions"
)


const (
	sessionName    = "example-facebook-app"
	sessionSecret  = "example cookie signing secret"
	sessionUserKey = "facebookID"
)

// sessionStore encodes and decodes session data stored in signed cookies
var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)
//
//// Config configures the main ServeMux.
//type Config struct {
//	FacebookClientID     string
//	FacebookClientSecret string
//}

// ToDo: need redis
var TokenUser = map[string]int64{
	"00000000": 1,
}

func TryToGetUserIdByToken(token string) (int64, error){
	if val, ok := TokenUser[token]; ok {
		return val, nil
	}
	return 0, errors.New("can't find user by token")
}
// Middleware function, which will be called for each request
func AuthenticationMiddleware(next http.Handler) http.Handler {
	log.Println(next)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")
		if r.URL.String() != "/" {
			if user, found := TokenUser[token]; found {
				// We found the token in our map
				log.Printf("Authenticated user %d\n", user)
				// Pass down the request to the next middleware (or final handler)
				next.ServeHTTP(w, r)
			} else {
				// Write an error and stop the handler chain
				log.Println("Unauthorized request at: ", r.URL)
				http.Error(w, "Forbidden", http.StatusForbidden)
				//http.Redirect(w, r, "/profile", http.StatusForbidden)
			}
		}

	})
	return nil
}

//
//// New returns a new ServeMux with app routes.
//func New(config *Config) *http.ServeMux {
//	mux := http.NewServeMux()
//	mux.HandleFunc("/", welcomeHandler)
//	mux.Handle("/profile", requireLogin(http.HandlerFunc(profileHandler)))
//	mux.HandleFunc("/logout", logoutHandler)
//	// 1. Register Login and Callback handlers
//	oauth2Config := &oauth2.Config{
//		ClientID:     config.FacebookClientID,
//		ClientSecret: config.FacebookClientSecret,
//		RedirectURL:  "http://localhost:8080/user/callback",
//		Endpoint:     facebookOAuth2.Endpoint,
//		Scopes:       []string{"email"},
//	}
//	// state param cookies require HTTPS by default; disable for localhost development
//	stateConfig := gologin.DebugOnlyCookieConfig
//	mux.Handle("/user/login", facebook.StateHandler(stateConfig, facebook.LoginHandler(oauth2Config, nil)))
//	mux.Handle("/user/callback", facebook.StateHandler(stateConfig, facebook.CallbackHandler(oauth2Config, issueSession(), nil)))
//	return mux
//}

// issueSession issues a cookie session after successful Facebook login
func issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		facebookUser, err := facebook.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 2. Implement a success handler to issue some form of session
		session := sessionStore.New(sessionName)
		session.Values[sessionUserKey] = facebookUser.ID
		log.Println(facebookUser.Name)
		log.Println(facebookUser.Email)
		log.Println(facebookUser.ID)
		session.Save(w)
		log.Println(session)
		http.Redirect(w, req, "/profile", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

// welcomeHandler shows a welcome message and login button.
func welcomeHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	if isAuthenticated(req) {
		http.Redirect(w, req, "/profile", http.StatusFound)
		return
	}
	page, _ := ioutil.ReadFile("home.html")
	fmt.Fprintf(w, string(page))
}

// profileHandler shows protected user content.
func ProfileHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, `<p>You are logged in!</p><form action="/logout" method="post"><input type="submit" value="Logout"></form>`)
}

// logoutHandler destroys the session on POSTs and redirects to home.
func logoutHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		sessionStore.Destroy(w, sessionName)
	}
	http.Redirect(w, req, "/", http.StatusFound)
}

// requireLogin redirects unauthenticated users to the login route.
//func RequireLogin(next http.Handler) (w http.ResponseWriter, req *http.Request) {
//
//	//_ := func(w http.ResponseWriter, req *http.Request) {
//		if !isAuthenticated(req) {
//			http.Redirect(w, req, "/", http.StatusFound)
//			return
//		}
//		next.ServeHTTP(w, req)
//	//}
//	return w, req
//}

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

// isAuthenticated returns true if the user has a signed session cookie.
func isAuthenticated(req *http.Request) bool {
	if _, err := sessionStore.Get(req, sessionName); err == nil {
		return true
	}
	return false
}
