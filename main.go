package main

import (
	"embed"
	"fmt"

	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/embed"

	"github.com/joho/godotenv"
)


var db *sqlite3.Conn

// content holds our static web server content.
//go:embed static/*
var content embed.FS
var store *sessions.CookieStore

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	var port string
	var exists bool
	if port, exists = os.LookupEnv("PORT"); !exists {
		port = "8080"
	}

	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

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

    r := mux.NewRouter()

	// set up OIDC
	setupOIDC()

	// set up HTTP
	r.HandleFunc("/login", initOIDCAuth)
	r.HandleFunc("/oidc-response", handleOIDCResponse)
	r.HandleFunc("/profile", getProfile)
	r.HandleFunc("/logout", logout)

    r.HandleFunc("/events", createEvent).Methods("POST")
    r.HandleFunc("/events", removeEvent).Methods("DELETE")
    r.HandleFunc("/events", listEvents).Methods("GET")
    r.HandleFunc("/events/{id:[0-9]+}", viewEvent).Methods("GET")
	
	if _, exists := os.LookupEnv("DEBUG"); exists {
		r.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))
	} else {
		static, _ := fs.Sub(content, "static")
		r.PathPrefix("/").Handler(http.FileServer(http.FS(static)))
	}

	if _, exists := os.LookupEnv("PROD"); exists {
		port = fmt.Sprintf(":%s", port)
		fmt.Println("Server started at", port)
		http.ListenAndServe(port, r)
	} else {
		log.Println("Starting server on :443")
		log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", r))
	}
}
