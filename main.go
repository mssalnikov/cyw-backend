package main

import (
	"fmt"
	"log"
	"net/http"

	"./conf"
	_ "github.com/lib/pq"
	_ "./models"

	"./handlers/users"
	"./utils"
	"github.com/jmoiron/sqlx"

	facebookOAuth2 "golang.org/x/oauth2/facebook"
	"github.com/gorilla/mux"
	"github.com/dghubble/gologin/facebook"
	"golang.org/x/oauth2"
	"github.com/dghubble/gologin"
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
		RedirectURL:  "http://localhost:3001/facebook/callback",
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


	//mux.HandleFunc("/user/login", users.Login).Methods("GET")
	//mux.HandleFunc("/user/logout", userHandler.Logout).Methods("GET")

	//mux.HandleFunc("/logout", logoutHandler)
	////welcome handler
	//mx.HandleFunc("/", welcomeHandler.WelcomePage).Methods("GET")
	//// user handler
	//mx.HandleFunc("/profile", utils.RequireLogin(userHandler.GetUser)).Methods("GET")
	//mx.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	//mx.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	//mx.HandleFunc("/users", userHandler.UpdateUser).Methods("PUT")
	//mx.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")
	//// ToDo: fix method to post
	//mx.HandleFunc("/user/login", userHandler.Login).Methods("GET")
	//mx.HandleFunc("/user/logout", userHandler.Logout).Methods("GET")

	//// // read credentials from environment variables if available
	//fbConfig := &utils.Config{
	//	FacebookClientID:     config.Auth.FBClient,
	//	FacebookClientSecret: config.Auth.FBSecret,
	//}

	//log.Printf("Starting Server listening on %s\n", hostBind)
	//err = http.ListenAndServe(hostBind, utils.New(fbConfig))
	//checkError(err)
	//
	// static
	//mx.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	// negroni
	//ng := negroni.New()
	//ng.UseHandler(mx)

	// start server
	log.Println("Listening on", hostBind)
	go http.ListenAndServe(hostBind, r)
	//checkError(err)
	http.ListenAndServe(hostFacebookBind, mx)
	//checkError(err)
}
