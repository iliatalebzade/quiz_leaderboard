package player_score

// PlayerScore represents the structure for storing player information and their score.
// It includes the player's ID, name, and current score, with corresponding JSON and BSON annotations
// for serialization and storage in MongoDB.
type PlayerScore struct {
	PlayerID   string `json:"player_id" bson:"player_id"`     // Unique identifier for the player
	PlayerName string `json:"player_name" bson:"player_name"` // Player's name
	Score      int    `json:"score" bson:"score"`             // Player's current score
}
