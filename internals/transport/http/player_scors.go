package http

import (
	"quiz/internals/domain/player_score"
	"quiz/internals/service"

	"github.com/gin-gonic/gin"
)

type PlayerScoresHandler struct {
	Service *service.PlayerScoreService // Service to handle player score operations
}

// NewPlayerScoreHandler initializes a new PlayerScoresHandler with the provided service.
func NewPlayerScoreHandler(service *service.PlayerScoreService) *PlayerScoresHandler {
	return &PlayerScoresHandler{Service: service}
}

// AddOrUpdateHandler handles requests to add or update a player's score.
func (psh *PlayerScoresHandler) AddOrUpdateHandler(c *gin.Context) {
	var req player_score.PlayerScore
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// Update or insert the player score via the service
	if err := psh.Service.AddOrUpdatePlayerScore(req); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update player score"})
		return
	}

	c.JSON(200, gin.H{"message": "Player score added or updated"})
}

// TopPlayersHandler retrieves and returns the top players.
func (psh *PlayerScoresHandler) TopPlayersHandler(c *gin.Context) {
	// Fetch the top players via the service
	topPlayers, err := psh.Service.GetTopPlayers(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve top players"})
		return
	}

	c.JSON(200, gin.H{"top_players": topPlayers})
}

// GetPointsHandler fetches and returns the score for a specific player by their ID.
func (psh *PlayerScoresHandler) GetPointsHandler(c *gin.Context) {
	playerID := c.Param("id")

	// Get the player score via the service
	playerScore, err := psh.Service.GetPlayerScore(playerID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Player not found"})
		return
	}

	c.JSON(200, gin.H{"player_id": playerID, "score": playerScore})
}
