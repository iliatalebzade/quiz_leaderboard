package service

import (
	"context"
	"fmt"
	"quiz/internals/domain/player_score"
	"quiz/internals/repositories"

	"go.uber.org/zap"
)

type PlayerScoreService struct {
	DBClient    repositories.IDBRepository    // Interface for database operations
	CacheClient repositories.ICacheRepository // Interface for cache operations
	CTX         context.Context               // Context for managing request-scoped values
	Logger      *zap.Logger                   // Logger for structured logging
}

// NewPlayerScoreService initializes a new PlayerScoreService with the provided database, cache clients, context, and logger.
func NewPlayerScoreService(db_client repositories.IDBRepository, cache_client repositories.ICacheRepository, ctx context.Context, custom_logger *zap.Logger) *PlayerScoreService {
	return &PlayerScoreService{
		DBClient:    db_client,
		CacheClient: cache_client,
		CTX:         ctx,
		Logger:      custom_logger,
	}
}

// AddOrUpdatePlayerScore adds or updates the player's score in the database and cache.
func (pss *PlayerScoreService) AddOrUpdatePlayerScore(playerScore player_score.PlayerScore) error {
	pss.Logger.Info("AddOrUpdatePlayerScore method called", zap.String("player_id", playerScore.PlayerID))

	// Update or insert player score in the database
	if err := pss.DBClient.UpdateOrInsertPlayerScore(playerScore); err != nil {
		pss.Logger.Error("Error updating or inserting player score in DB", zap.String("player_id", playerScore.PlayerID), zap.Error(err))
		return err
	}

	pss.Logger.Info("Player score updated/inserted in DB", zap.String("player_id", playerScore.PlayerID))

	// Update the player's cache asynchronously (ZSET and HASH)
	go func() {
		if err := pss.CacheClient.UpdatePlayerCache("leaderboard", playerScore); err != nil {
			pss.Logger.Error("Error updating the cache for player", zap.String("player_id", playerScore.PlayerID), zap.Error(err))
		} else {
			pss.Logger.Info("Player cache updated successfully", zap.String("player_id", playerScore.PlayerID))
		}
	}()

	pss.Logger.Info(fmt.Sprintf("Create or update operations were successful for player: %v", playerScore))
	return nil
}

// GetTopPlayers retrieves the top players from cache or database.
func (pss *PlayerScoreService) GetTopPlayers(ctx context.Context) ([]player_score.PlayerScore, error) {
	pss.Logger.Info("GetTopPlayers method called")

	// Attempt to retrieve leaderboard from cache
	leaderboard, err := pss.CacheClient.GetSetByKey("leaderboard")
	if err != nil {
		pss.Logger.Error("Error retrieving records from Cache", zap.Error(err))
	}

	// Cache miss, retrieve from database if cache is empty
	if len(leaderboard) == 0 {
		pss.Logger.Info("Cache miss, retrieving from DB")

		// Fetch top players from the database
		topPlayers, err := pss.DBClient.GetTopPlayers()
		if err != nil {
			pss.Logger.Error("Error retrieving records from DB", zap.Error(err))
			return nil, err
		}

		pss.Logger.Info("Top players retrieved from DB", zap.Int("count", len(topPlayers)))

		// Cache the leaderboard asynchronously
		go func() {
			for _, player := range topPlayers {
				if err := pss.CacheClient.InsertRecord("leaderboard", player.PlayerID, player.PlayerName, float64(player.Score)); err != nil {
					pss.Logger.Error("Error inserting new records into cache", zap.String("player_id", player.PlayerID), zap.Error(err))
				} else {
					pss.Logger.Info("Player score cached successfully", zap.String("player_id", player.PlayerID))
				}
			}
		}()

		return topPlayers, nil
	}

	// Return leaderboard from cache if available
	pss.Logger.Info("Cached response provided", zap.Int("count", len(leaderboard)))
	return leaderboard, nil
}

// GetPlayerScore fetches a player's score from the database.
func (pss *PlayerScoreService) GetPlayerScore(playerID string) (int, error) {
	pss.Logger.Info("GetPlayerScore method called", zap.String("player_id", playerID))

	// Retrieve player score from the database
	score, err := pss.DBClient.GetPlayerScore(playerID)
	if err != nil {
		pss.Logger.Error("Error fetching player score from DB", zap.String("player_id", playerID), zap.Error(err))
		return 0, err
	}

	pss.Logger.Info("Player score retrieved successfully", zap.String("player_id", playerID), zap.Int("score", score))
	return score, nil
}
