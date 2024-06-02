package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
    var err error
    db, err = sql.Open("sqlite3", "./events.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if err := createTable(); err != nil {
        log.Fatal(err)
    }

    r := mux.NewRouter()
    r.HandleFunc("/events", createEvent).Methods("POST")
    r.HandleFunc("/events", listEvents).Methods("GET")
    r.HandleFunc("/events/{id:[0-9]+}", viewEvent).Methods("GET")
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

    fmt.Println("Server started at :8080")
    http.ListenAndServe(":8080", r)
}

