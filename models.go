package main

import (
	"fmt"
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
    err := db.Exec(query)
    return err
}

func insertEvent(event Event) (int, error) {
    stmt, _, err := db.Prepare(`INSERT INTO events (title, description, location, start_time, end_time, creator) VALUES (?, ?, ?, ?, ?, ?)`)
    if err != nil {
        return 0, err
    }
	defer stmt.Close()
	stmt.BindText(1, event.Title)
	stmt.BindText(2, event.Description)
	stmt.BindText(3, event.Location)
	stmt.BindTime(4, event.StartTime, timeLayout)
	stmt.BindTime(5, event.EndTime, timeLayout)
	stmt.BindText(6, event.Creator)
	err = stmt.Exec()
    if err != nil {
        return 0, err
    }
    //id, err := result.LastInsertId()
	id := 0
    return int(id), err
}

func deleteEvent(event Event, creator string) error {
    stmt, _, err := db.Prepare(`DELETE FROM events WHERE id = ? AND creator = ?`)
    if err != nil {
        return err
    }
	defer stmt.Close()
	stmt.BindInt(1, event.ID)
	stmt.BindText(2, creator)
	err = stmt.Exec()
    if err != nil {
        return err
    }
	affected := db.Changes()
	if affected == 0 {
		return fmt.Errorf("Invalid event ID or insufficient permissions")
	}
	return nil
}

func getEvents() ([]Event, error) {
    stmt, _, err := db.Prepare("SELECT id, title, description, location, start_time, end_time, creator FROM events")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()


    events := []Event{}
	for stmt.Step() {
		event := Event{
			ID:          stmt.ColumnInt(0),
			Title:       stmt.ColumnText(1),
			Description: stmt.ColumnText(2),
			Location:    stmt.ColumnText(3),
			StartTime:   stmt.ColumnTime(4, timeLayout),
			EndTime:     stmt.ColumnTime(5, timeLayout),
			Creator:     stmt.ColumnText(6),
		}
        events = append(events, event)
    }
	if err := stmt.Err(); err != nil {
		return nil, err
	}
    return events, nil
}

func getEventByID(id int) (Event, error) {
    stmt, _, err := db.Prepare("SELECT id, title, description, location, start_time, end_time, creator FROM events WHERE id = ?")
	if err != nil {
		return Event{}, err
	}

	defer stmt.Close()
	stmt.Step()

	event := Event{
		ID:          stmt.ColumnInt(0),
		Title:       stmt.ColumnText(1),
		Description: stmt.ColumnText(2),
		Location:    stmt.ColumnText(3),
		StartTime:   stmt.ColumnTime(4, timeLayout),
		EndTime:     stmt.ColumnTime(5, timeLayout),
		Creator:     stmt.ColumnText(6),
	}
    return event, stmt.Err()
}

func hasConflict(event Event) (bool, error) {
    stmt, _, err := db.Prepare(`SELECT COUNT(*) FROM events WHERE location = ? AND (
                (start_time < ? AND end_time > ?) OR
                (start_time < ? AND end_time > ?) OR
                (start_time >= ? AND start_time < ?) OR
                (end_time > ? AND end_time <= ?)
            )`)
	if err != nil {
		return false, err
	}

	stmt.BindText(1, event.Location)
	stmt.BindTime(2, event.EndTime, timeLayout)
	stmt.BindTime(3, event.StartTime, timeLayout)
	stmt.BindTime(4, event.EndTime, timeLayout)
	stmt.BindTime(5, event.StartTime, timeLayout)
	stmt.BindTime(6, event.StartTime, timeLayout)
	stmt.BindTime(7, event.EndTime, timeLayout)
	stmt.BindTime(8, event.StartTime, timeLayout)
	stmt.BindTime(9, event.EndTime, timeLayout)

	defer stmt.Close()

	stmt.Step()
	count := stmt.ColumnInt(0)
	return count > 0, stmt.Err()
}

