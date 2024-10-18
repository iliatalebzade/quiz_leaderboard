package player_score

type PlayerScore struct {
	PlayerID   string `json:"player_id" bson:"player_id"`
	PlayerName string `json:"player_name" bson:"player_name"`
	Score      int    `json:"score" bson:"score"`
}
