package datastore

import (
	"database/sql"
	"fmt"
)

var playersSchema string = `
CREATE TABLE IF NOT EXISTS players (
	name TEXT NOT NULL PRIMARY KEY
)`

var gamesSchema string = `
CREATE TABLE IF NOT EXISTS games (
	id TEXT PRIMARY KEY,
	players TEXT NOT NULL, -- JSON array of player IDs
	tries INTEGER NOT NULL,
	duration INTEGER NOT NULL,
	date TEXT NOT NULL,
	won BOOLEAN NOT NULL
)`

// Datastore handles database operations for players
type Datastore struct {
	db *sql.DB
}

// NewDatastore creates a new Datastore instance
func NewDatastore(dbPath string) (*Datastore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Create table if it doesn't exist
	_, err = db.Exec(playersSchema)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	_, err = db.Exec(gamesSchema)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	return &Datastore{db: db}, nil
}
