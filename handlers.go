package main

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"

    "github.com/gorilla/mux"
)

const timeLayout = "2006-01-02T15:04:05Z07:00"

func createEvent(w http.ResponseWriter, r *http.Request) {
    var event Event
    if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var err error
    event.StartTime, err = time.Parse(timeLayout, event.StartTime.Format(timeLayout))
    if err != nil {
        http.Error(w, "Invalid start time format", http.StatusBadRequest)
        return
    }

    event.EndTime, err = time.Parse(timeLayout, event.EndTime.Format(timeLayout))
    if err != nil {
        http.Error(w, "Invalid end time format", http.StatusBadRequest)
        return
    }

    // Set the creator
    event.Creator = "unset" // Replace with actual logic to determine the creator

    // Check for conflict
    conflict, err := hasConflict(event)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if conflict {
        http.Error(w, "Event conflicts with an existing event", http.StatusConflict)
        return
    }

    id, err := insertEvent(event)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    event.ID = id

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(event)
}

func removeEvent(w http.ResponseWriter, r *http.Request) {
    var event Event
    if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var err error
    event.StartTime, err = time.Parse(timeLayout, event.StartTime.Format(timeLayout))
    if err != nil {
        http.Error(w, "Invalid start time format", http.StatusBadRequest)
        return
    }

    event.EndTime, err = time.Parse(timeLayout, event.EndTime.Format(timeLayout))
    if err != nil {
        http.Error(w, "Invalid end time format", http.StatusBadRequest)
        return
    }

    // Check permissions
	creator := "unset"

    err = deleteEvent(event, creator)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(event)
}

func listEvents(w http.ResponseWriter, r *http.Request) {
    events, err := getEvents()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(events)
}

func viewEvent(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid event ID", http.StatusBadRequest)
        return
    }

    event, err := getEventByID(id)
    if err != nil {
        http.Error(w, "Event not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(event)
}

