package repositories

import (
	"quiz/internals/domain/player_score"
)

type IDBRepository interface {
	UpdateOrInsertPlayerScore(player player_score.PlayerScore) error // Update or insert player score in DB
	GetTopPlayers() ([]player_score.PlayerScore, error)              // Cache miss, get from DB
	GetPlayerScore(playerID string) (int, error)                     // Get Single score from DB
	Connect()                                                        // Open a connection
	Close()                                                          // Close the connection
}
