package datastore

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Save stores a player in the database
func (d *Datastore) Save(player *Player) error {
	_, err := d.db.Exec(
		"INSERT OR REPLACE INTO players (name) VALUES (?)",
		player.Name,
	)
	if err != nil {
		return fmt.Errorf("error saving player: %v", err)
	}
	return nil
}

// Create new player
func (d *Datastore) CreateNewPlayer(userName string) (*Player, error) {
	_, err := d.db.Exec(
		"INSERT INTO players (name) VALUES (?)",
		userName,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating player: %v", err)
	}
	return &Player{
		Name: userName,
	}, nil
}

// GetByName retrieves a player by their name
func (d *Datastore) GetByName(name string) (*Player, error) {
	player := &Player{}
	err := d.db.QueryRow(
		"SELECT name FROM players WHERE name = ?",
		name,
	).Scan(&player.Name)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error getting player: %v", err)
	}

	return player, nil
}

// GetByID retrieves a player by their ID
func (d *Datastore) GetByID(id string) (*Player, error) {
	player := &Player{}
	err := d.db.QueryRow(
		"SELECT name FROM players WHERE name = ?",
		id,
	).Scan(&player.Name)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting player: %v", err)
	}
	return player, nil
}

// GetAll retrieves all players from the database
func (d *Datastore) GetAll() ([]*Player, error) {
	rows, err := d.db.Query("SELECT name FROM players")
	if err != nil {
		return nil, fmt.Errorf("error querying players: %v", err)
	}
	defer rows.Close()

	var players []*Player
	for rows.Next() {
		player := &Player{}
		err := rows.Scan(&player.Name)
		if err != nil {
			return nil, fmt.Errorf("error scanning player: %v", err)
		}
		players = append(players, player)
	}
	return players, nil
}

// Close closes the database connection
func (d *Datastore) Close() error {
	return d.db.Close()
}

func Example() {
	// Create a new PlayerDB instance
	playerDB, err := NewDatastore("./players.db")
	if err != nil {
		panic(err)
	}
	defer playerDB.Close()

	// Create a new player
	player := &Player{
		Name: "John Doe",
	}

	// Save the player
	if err := playerDB.Save(player); err != nil {
		panic(err)
	}

	// Retrieve the player
	retrieved, err := playerDB.GetByID("player1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Retrieved player: %+v\n", retrieved)

	// Get all players
	allPlayers, err := playerDB.GetAll()
	if err != nil {
		panic(err)
	}
	fmt.Printf("All players: %+v\n", allPlayers)
}
