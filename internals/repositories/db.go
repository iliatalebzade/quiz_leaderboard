package repositories

import (
	"quiz/internals/domain/player_score"
)

// IDBRepository defines the operations for interacting with the database,
// specifically for managing player scores, including retrieval, insertion, and updates.
type IDBRepository interface {
	UpdateOrInsertPlayerScore(player player_score.PlayerScore) error // Insert a new player score or update an existing one in the database
	GetTopPlayers() ([]player_score.PlayerScore, error)              // Retrieve the top players' scores from the database (in case of cache miss)
	GetPlayerScore(playerID string) (int, error)                     // Retrieve a single player's score from the database by their ID
	Connect()                                                        // Establish a connection to the database
	Close()                                                          // Close the database connection
}
