package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // required
	// _ "github.com/mattes/migrate/source/file" // required

	"github.com/gkiryaziev/go-gorilla-sqlx/conf"
	"github.com/gkiryaziev/go-gorilla-sqlx/handlers/users"
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

	db, err := sqlx.Connect("postgres", sqlBind)
	checkError(err)
	db.SetMaxIdleConns(100)
	defer db.Close()

	// handlers
	userHandler := users.NewUserHandler(db)

	mx := mux.NewRouter()

	// user handler
	mx.HandleFunc("/api/v1/users", userHandler.GetUsers).Methods("GET")
	mx.HandleFunc("/api/v1/users/{id}", userHandler.GetUser).Methods("GET")
	mx.HandleFunc("/api/v1/users", userHandler.UpdateUser).Methods("PUT")
	mx.HandleFunc("/api/v1/users/{id}", userHandler.DeleteUser).Methods("DELETE")
	mx.HandleFunc("/api/v1/users", userHandler.InsertUser).Methods("POST")

	// static
	mx.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	// negroni
	ng := negroni.New()
	ng.UseHandler(mx)

	// start server
	log.Println("Listening on", hostBind)
	err = http.ListenAndServe(hostBind, ng)
	checkError(err)
}
