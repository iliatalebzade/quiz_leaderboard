package http

import (
	"quiz/internals/domain/player_score"
	"quiz/internals/service"

	"github.com/gin-gonic/gin"
)

type PlayerScoresHandler struct {
	Service *service.PlayerScoreService
}

func NewPlayerScoreHandler(service *service.PlayerScoreService) *PlayerScoresHandler {
	return &PlayerScoresHandler{Service: service}
}

// AddOrUpdateHandler updates or inserts the player's score.
func (psh *PlayerScoresHandler) AddOrUpdateHandler(c *gin.Context) {
	var req player_score.PlayerScore
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// Call the service method to update the player score
	if err := psh.Service.AddOrUpdatePlayerScore(req); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update player score"})
		return
	}

	c.JSON(200, gin.H{"message": "Player score added or updated"})
}

// TopPlayersHandler returns the top 10 players.
func (psh *PlayerScoresHandler) TopPlayersHandler(c *gin.Context) {
	// Retrieve the top players from the service
	topPlayers, err := psh.Service.GetTopPlayers(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve top players"})
		return
	}

	c.JSON(200, gin.H{"top_players": topPlayers})
}

// GetPointsHandler returns the score for a specific player.
func (psh *PlayerScoresHandler) GetPointsHandler(c *gin.Context) {
	playerID := c.Param("id")

	// Fetch the player score from the service
	playerScore, err := psh.Service.GetPlayerScore(playerID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Player not found"})
		return
	}

	c.JSON(200, gin.H{"player_id": playerID, "score": playerScore})
}
