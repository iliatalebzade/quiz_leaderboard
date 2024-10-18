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
	ctx := context.Background() // Create a background context for the application

	// Load application configuration from environment variables or .env file
	cfg := config.LoadConfig()

	// Connect to MongoDB using the URI from the configuration
	mongoClient := repositories.NewMongoDBClient(ctx, cfg.MongoDBURI)
	mongoClient.Connect()     // Establish the MongoDB connection
	defer mongoClient.Close() // Ensure the connection is closed on exit

	// Connect to Redis using the address, password, and database index from the configuration
	redisClient := repositories.NewRedisClient(ctx, cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDBIndex)
	redisClient.Connect()     // Establish the Redis connection
	defer redisClient.Close() // Ensure the connection is closed on exit

	// Create a production logger using Uber's Zap library
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err) // Exit if logger creation fails
	}
	defer logger.Sync() // Ensure all log entries are flushed before exiting

	// Log an informational message indicating the application is starting
	logger.Info("Application starting",
		zap.String("mongo_uri", cfg.MongoDBURI), // Log MongoDB URI
		zap.String("redis_addr", cfg.RedisAddr), // Log Redis address
	)

	// Setup the Player Score service with dependencies
	playerScoresService := service.NewPlayerScoreService(
		mongoClient, // MongoDB client
		redisClient, // Redis client
		ctx,         // Context for cancellation and deadlines
		logger,      // Logger for the service
	)

	// Setup the HTTP handlers for player scores
	playerScoresHandler := http.NewPlayerScoreHandler(playerScoresService)

	// Initialize the Gin router and setup routes grouped under the /points subroute
	router := gin.Default()
	v1 := router.Group("/points")
	{
		// Route to add or update player scores
		v1.POST("/add_or_update", playerScoresHandler.AddOrUpdateHandler)

		// Route to get the top players' scores
		v1.GET("/top_players", playerScoresHandler.TopPlayersHandler)

		// Route to get points for a specific player by ID
		v1.GET("/get_points/:id", playerScoresHandler.GetPointsHandler)
	}

	// Start the HTTP server on port 8000
	router.Run(":8000")
}
