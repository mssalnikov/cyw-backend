package utils

import (
	"net/http"
	"log"
	"github.com/go-redis/redis"
	"fmt"
)


//const (
//	sessionName    = "example-facebook-app"
//	sessionSecret  = "example cookie signing secret"
//	sessionUserKey = "facebookID"
//)
//
//// sessionStore encodes and decodes session data stored in signed cookies
//var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)
//
//
//func TryToGetUserIdByToken(token string) (int64, error){
//	if val, ok := TokenUser[token]; ok {
//		return val, nil
//	}
//	return 0, errors.New("can't find user by token")
//}
// Middleware function, which will be called for each request
func AuthenticationMiddleware(next http.Handler) http.Handler {
	log.Println(next)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")
		log.Println(r.URL.String())
		switch url := r.URL.String(); url {
		case "/":
			next.ServeHTTP(w, r)
		case "/auth":
			next.ServeHTTP(w, r)
		default:
			userId, err := RedisCon.Get(fmt.Sprintf("TOKEN:%s", token)).Result()
			if err == redis.Nil {
				log.Println("Unauthorized request at: ", r.URL)
				http.Error(w, "Forbidden", http.StatusForbidden)
			} else if err != nil {
				log.Println(err)
				http.Error(w, "Bad token", http.StatusBadRequest)
			} else {
				log.Printf("Authenticated user %d\n", userId)
				// Pass down the request to the next middleware (or final handler)
				next.ServeHTTP(w, r)
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

