package datastore

import (
	"fmt"
)

func (d *Datastore) GetGamesByPlayerId(playerId string) ([]*Game, error) {
	rows, err := d.db.Query(`
SELECT
id, players, tries, duration, date, won
FROM games
WHERE players = ? AND won = FALSE
ORDER BY date ASC
`, playerId)
	if err != nil {
		return nil, fmt.Errorf("error querying games: %v", err)
	}

	defer rows.Close()

	var games []*Game

	for rows.Next() {
		game := &Game{}
		err := rows.Scan(&game.Id, &game.Players, &game.Tries, &game.Duration, &game.Date, &game.Won)
		if err != nil {
			return nil, fmt.Errorf("error scanning game: %v", err)
		}
		games = append(games, game)
	}

	return games, nil
}

func (d *Datastore) SaveGame(game *Game) error {
	_, err := d.db.Exec(`
INSERT INTO games
(id, players, tries, duration, date, won)
VALUES (?, ?, ?, ?, ?, ?)
`,
		game.Id, game.Players, game.Tries, game.Duration, game.Date, game.Won,
	)
	if err != nil {
		return fmt.Errorf("error saving game: %v", err)
	}

	return nil
}
