package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/embed"
)


var db *sqlite3.Conn

// content holds our static web server content.
//go:embed static/*
var content embed.FS

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// set up DB
    var err error
    db, err = sqlite3.Open("./events.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if err := createTable(); err != nil {
        log.Fatal(err)
    }
	fmt.Println("Set up DB")

    r := mux.NewRouter()

	// set up OIDC
	setupOIDC()

	// set up HTTP
	r.HandleFunc("/auth", initOIDCAuth)
	r.HandleFunc("/oidc-response", handleOIDCResponse)

    r.HandleFunc("/events", createEvent).Methods("POST")
    r.HandleFunc("/events", removeEvent).Methods("DELETE")
    r.HandleFunc("/events", listEvents).Methods("GET")
    r.HandleFunc("/events/{id:[0-9]+}", viewEvent).Methods("GET")
	
	static, _ := fs.Sub(content, "static")
	r.PathPrefix("/").Handler(http.FileServer(http.FS(static)))
	//r.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	//http.Handle("/static/", )

	//port = fmt.Sprintf(":%s", port)
    //fmt.Println("Server started at", port)
    //http.ListenAndServe(port, r)

	log.Println("Starting server on :443")
	log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", r))
}

