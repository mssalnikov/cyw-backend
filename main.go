package main

import (
	"fmt"
	"log"
	"net/http"

	"./conf"
	_ "./models"
	_ "github.com/lib/pq"

	"./handlers/users"
	"./utils"
	"github.com/jmoiron/sqlx"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/facebook"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	facebookOAuth2 "golang.org/x/oauth2/facebook"
	"github.com/go-redis/redis"
)

// checkError check errors
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func YourHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

type authenticationMiddleware struct {
	tokenUsers map[string]string
}

func main() {

	// load config
	config, err := conf.NewConfig("config.yaml").Load()
	checkError(err)

	// redis connection client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := redisClient.Ping().Result()
	fmt.Println(pong, err)

	// mysql connection string
	sqlBind := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		config.DB.UserName,
		config.DB.UserPassword,
		config.DB.Database,
	)

	// http server address and port
	hostBind := fmt.Sprintf("%s:%s",
		config.Host.IP,
		config.Host.Port,
	)
	hostFacebookBind := fmt.Sprintf("%s:%s",
		config.Auth.IP,
		config.Auth.Port,
	)

	db, err := sqlx.Connect("postgres", sqlBind)
	checkError(err)
	db.SetMaxIdleConns(100)
	defer db.Close()

	// 	// handlers
	//userHandler := users.NewUserHandler(db)
	//welcomeHandler := welcome.NewWelcomeHandler()

	//mx := mux.NewRouter()
	//mux := http.NewServeMux()
	//mux.HandleFunc("/", welcomeHandler.WelcomePage)
	//mux.Handle("/profile", utils.RequireLogin(http.HandlerFunc(utils.ProfileHandler)))
	//

	uh := users.NewUserHandler(db)

	oauth2Config := &oauth2.Config{
		ClientID:     config.Auth.FBClient,
		ClientSecret: config.Auth.FBSecret,
		RedirectURL:  "http://" + config.AuthHost.IP + ":" + config.AuthHost.Port + "/facebook/callback",
		Endpoint:     facebookOAuth2.Endpoint,
		Scopes:       []string{"email"},
	}
	// state param cookies require HTTPS by default; disable for localhost development

	r := mux.NewRouter()
	r.HandleFunc("/", YourHandler).Methods("GET")
	r.HandleFunc("/profile", YourHandler).Methods("GET")

	//r.Handle("/facebook/login", facebook.StateHandler(stateConfig, facebook.LoginHandler(oauth2Config, nil)))
	//r.Handle("/facebook/callback", facebook.StateHandler(stateConfig, facebook.CallbackHandler(oauth2Config, uh.IssueSession(), nil)))

	r.Use(utils.AuthenticationMiddleware)

	mx := http.NewServeMux()
	stateConfig := gologin.DebugOnlyCookieConfig
	mx.Handle("/facebook/login", facebook.StateHandler(stateConfig, facebook.LoginHandler(oauth2Config, nil)))
	mx.Handle("/facebook/callback", facebook.StateHandler(stateConfig, facebook.CallbackHandler(oauth2Config, uh.IssueSession(hostFacebookBind), nil)))

	// start server
	log.Println("Listening on", hostBind)
	go http.ListenAndServe(hostBind, r)
	//checkError(err)
	http.ListenAndServe(hostFacebookBind, mx)
	//checkError(err)
}
