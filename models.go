package main

import (
    _ "database/sql"
    "time"
)

type Event struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Location    string    `json:"location"`
    StartTime   time.Time `json:"start_time"`
    EndTime     time.Time `json:"end_time"`
    Creator     string    `json:"creator"`
}

func createTable() error {
    query := `
    CREATE TABLE IF NOT EXISTS events (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        description TEXT,
        location TEXT,
        start_time DATETIME,
        end_time DATETIME,
        creator TEXT
    );`
    _, err := db.Exec(query)
    return err
}

func insertEvent(event Event) (int, error) {
    query := `INSERT INTO events (title, description, location, start_time, end_time, creator) VALUES (?, ?, ?, ?, ?, ?)`
    result, err := db.Exec(query, event.Title, event.Description, event.Location, event.StartTime, event.EndTime, event.Creator)
    if err != nil {
        return 0, err
    }
    id, err := result.LastInsertId()
    return int(id), err
}

func getEvents() ([]Event, error) {
    rows, err := db.Query("SELECT id, title, description, location, start_time, end_time, creator FROM events")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    events := []Event{}
    for rows.Next() {
        var event Event
        err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.Location, &event.StartTime, &event.EndTime, &event.Creator)
        if err != nil {
            return nil, err
        }
        events = append(events, event)
    }
    return events, nil
}

func getEventByID(id int) (Event, error) {
    var event Event
    query := "SELECT id, title, description, location, start_time, end_time, creator FROM events WHERE id = ?"
    row := db.QueryRow(query, id)
    err := row.Scan(&event.ID, &event.Title, &event.Description, &event.Location, &event.StartTime, &event.EndTime, &event.Creator)
    return event, err
}

func hasConflict(event Event) (bool, error) {
    query := `SELECT COUNT(*) FROM events WHERE location = ? AND (
                (start_time < ? AND end_time > ?) OR
                (start_time < ? AND end_time > ?) OR
                (start_time >= ? AND start_time < ?) OR
                (end_time > ? AND end_time <= ?)
            )`
    var count int
    err := db.QueryRow(query, event.Location, event.EndTime, event.StartTime, event.EndTime, event.StartTime, event.StartTime, event.EndTime, event.StartTime, event.EndTime).Scan(&count)
    return count > 0, err
}

