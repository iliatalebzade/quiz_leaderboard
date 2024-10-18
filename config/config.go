package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all the necessary configuration settings for the application,
// including database URIs and Redis connection details.
type Config struct {
	MongoDBURI    string // MongoDB connection URI
	RedisAddr     string // Redis server address
	RedisPassword string // Redis password (if required)
	RedisDBIndex  int    // Redis database index to use
}

// LoadConfig reads the configuration from the .env file or environment variables.
// It sets default values for each configuration setting if not found in .env.
// Logs a message if the .env file is not found, indicating fallback to environment variables.
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables directly")
	}

	return &Config{
		MongoDBURI:    getEnv("MONGODB_URI", "mongodb://localhost:27017"), // Default MongoDB URI
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),             // Default Redis address
		RedisPassword: getEnv("REDIS_PASSWORD", ""),                       // Default Redis password (empty)
		RedisDBIndex:  getEnvAsInt("REDIS_DB_INDEX", 0),                   // Default Redis DB index
	}
}

// getEnv retrieves the value of the environment variable identified by key.
// If the variable is not set, it returns the provided fallback value.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getEnvAsInt retrieves the value of the environment variable identified by key
// and converts it to an integer. If the variable is not set or conversion fails,
// it returns the provided fallback integer value.
func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}
