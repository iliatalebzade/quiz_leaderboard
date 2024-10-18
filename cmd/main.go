package main

import (
	"context"
	"log"
	"quiz/internals/repositories"
	"quiz/internals/service"
	"quiz/internals/transport/http"

	"quiz/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	mongoClient := repositories.NewMongoDBClient(ctx, cfg.MongoDBURI)
	mongoClient.Connect()
	defer mongoClient.Close()

	// Connect to Redis
	redisClient := repositories.NewRedisClient(ctx, cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDBIndex)
	redisClient.Connect()
	defer redisClient.Close()

	// Create a logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Log an informational message
	logger.Info("Application starting",
		zap.String("mongo_uri", cfg.MongoDBURI),
		zap.String("redis_addr", cfg.RedisAddr),
	)

	// Setup Service
	playerScoresService := service.NewPlayerScoreService(
		mongoClient,
		redisClient,
		ctx,
		logger, // Pass the zap logger to your service
	)

	// Setup Handlers
	playerScoresHandler := http.NewPlayerScoreHandler(playerScoresService)

	router := gin.Default()
	v1 := router.Group("/points")
	{
		v1.POST("/add_or_update", playerScoresHandler.AddOrUpdateHandler)
		v1.GET("/top_players", playerScoresHandler.TopPlayersHandler)
		v1.GET("/get_points/:id", playerScoresHandler.GetPointsHandler)
	}

	router.Run(":8000")
}
