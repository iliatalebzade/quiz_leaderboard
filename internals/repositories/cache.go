package repositories

import "quiz/internals/domain/player_score"

type ICacheRepository interface {
	UpdatePlayerCache(string, player_score.PlayerScore) error           // Invalidate cache for leaderboard or user score
	GetSetByKey(key string) ([]player_score.PlayerScore, error)         // Get the leaderboard by key
	GetRecordByKey(key string) (player_score.PlayerScore, error)        // Get a player's score by key
	InsertRecord(key, playerID, playername string, score float64) error // Insert one player's score
	Connect()                                                           // Open a connection
	Close()                                                             // Close the connection
}
