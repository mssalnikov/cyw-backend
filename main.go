package main

import (
	"fmt"
	"log"
	"net/http"

	"./conf"
	_ "./models"
	_ "github.com/lib/pq"

	"./utils"
	"github.com/gorilla/mux"
	"github.com/go-redis/redis"
	"database/sql"
	"./handlers/users"
	"./handlers/events"
)

// checkError check errors
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// load config
	config, err := conf.NewConfig("config.yaml").Load()
	checkError(err)
	sqlBind := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		config.DB.UserName,
		config.DB.UserPassword,
		config.DB.Database,
	)

	// http server address and port
	hostBind := fmt.Sprintf("%s:%s",
		"localhost",
		"3000",
	)

	utils.DBCon, err = sql.Open("postgres", sqlBind)
	utils.RedisCon = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	uh := users.NewUserHandler()
	eh := events.NewEventHandler()

	r := mux.NewRouter()

	// users
	r.HandleFunc("/auth", uh.Auth).Methods("POST")
	r.HandleFunc("/profile", uh.GetUser).Methods("GET")

	// events
	r.HandleFunc("/new_event", eh.NewEvent).Methods("POST")
	r.HandleFunc("/my_events", eh.MyEvents).Methods("GET")
	r.HandleFunc("/check_event", eh.CheckEvent).Methods("POST")
	r.HandleFunc("/events", eh.AllEvents).Methods("GET")
	r.HandleFunc("/event", eh.GetEvent).Methods("GET")
	r.HandleFunc("/point", eh.GetPoint).Methods("GET")
	r.HandleFunc("/join_event", eh.JoinEvent).Methods("POST")
	r.HandleFunc("/enter_token", eh.EnterToken).Methods("POST")
	r.HandleFunc("/answer_question", eh.AnswerQuestion).Methods("POST")

	r.Use(utils.AuthenticationMiddleware)

	// start server
	log.Println("Listening on", hostBind)
	http.ListenAndServe(hostBind, r)
}
