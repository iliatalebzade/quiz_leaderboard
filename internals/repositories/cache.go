package repositories

import "quiz/internals/domain/player_score"

// ICacheRepository defines the operations for interacting with a cache system,
// specifically for storing and retrieving player scores and leaderboard data.
type ICacheRepository interface {
	UpdatePlayerCache(string, player_score.PlayerScore) error           // Update or invalidate the cache for a player's score or leaderboard
	GetSetByKey(key string) ([]player_score.PlayerScore, error)         // Retrieve the leaderboard (set of player scores) by a cache key
	GetRecordByKey(key string) (player_score.PlayerScore, error)        // Retrieve a specific player's score from the cache by key
	InsertRecord(key, playerID, playername string, score float64) error // Insert or update a player's score in the cache
	Connect()                                                           // Establish a connection to the cache
	Close()                                                             // Close the cache connection
}
