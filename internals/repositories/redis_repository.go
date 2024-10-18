package repositories

import (
	"context"
	"log"
	"quiz/internals/domain/player_score"
	"strconv"

	"github.com/go-redis/redis"
)

type RedisClient struct {
	Ctx      context.Context
	Addr     string
	Password string
	DB       int

	Client *redis.Client
}

func NewRedisClient(ctx context.Context, addr, password string, db int) *RedisClient {
	return &RedisClient{Ctx: ctx, Addr: addr, Password: password, DB: db}
}

// UpdatePlayerScore updates the player's score in the ZSET (leaderboard).
func (rr *RedisClient) UpdatePlayerScore(key, playerID string, score float64) error {
	// Add or update the player's score in the ZSET
	err := rr.Client.ZAdd(key, redis.Z{
		Score:  score,    // The player's new score
		Member: playerID, // The player ID is the member in the ZSET
	}).Err()
	if err != nil {
		log.Println("Failed to update player score in Redis ZSET:", err)
		return err
	}

	return nil
}

func (rr *RedisClient) UpdatePlayerCache(key string, playerScore player_score.PlayerScore) error {
	// Update the ZSET leaderboard (find and replace player's score)
	err := rr.UpdatePlayerScore(key, playerScore.PlayerID, float64(playerScore.Score))
	if err != nil {
		log.Println("Failed to update player score in Redis ZSET:", err)
		return err
	}

	// Update the HASH for the player with new details
	playerHash := map[string]interface{}{
		"PlayerID":   playerScore.PlayerID,
		"PlayerName": playerScore.PlayerName,
		"Score":      playerScore.Score,
	}

	err = rr.Client.HMSet("player:"+playerScore.PlayerID, playerHash).Err()
	if err != nil {
		log.Println("Failed to update Redis HASH for player:", playerScore.PlayerID, "err:", err)
		return err
	}

	return nil
}

// GetSetByKey fetches a sorted set from Redis and retrieves additional player details from the HASH.
func (rr *RedisClient) GetSetByKey(key string) ([]player_score.PlayerScore, error) {
	zSet, err := rr.Client.ZRevRangeWithScores(key, 0, -1).Result()
	if err != nil {
		log.Println("Failed to retrieve sorted set from Redis:", err)
		return nil, err
	}

	playerScores := make([]player_score.PlayerScore, len(zSet))
	for i, z := range zSet {
		playerID := z.Member.(string)

		// Fetch the playername from the HASH
		playername, err := rr.Client.HGet("player:"+playerID, "PlayerName").Result()
		if err != nil {
			log.Println("Failed to retrieve playername from Redis:", err)
			return nil, err
		}

		playerScores[i] = player_score.PlayerScore{
			PlayerID:   playerID,
			Score:      int(z.Score),
			PlayerName: playername,
		}
	}

	return playerScores, nil
}

func (rr *RedisClient) GetRecordByKey(playerID string) (player_score.PlayerScore, error) {
	// Retrieve the score from the sorted set (assuming the key is for a sorted set of scores)
	scoreResult, err := rr.Client.Get(playerID).Result()
	if err != nil {
		log.Println("Failed to get score from Redis:", err)
		return player_score.PlayerScore{}, err
	}

	// Convert the score result to an integer
	score, _ := strconv.Atoi(scoreResult)

	// Retrieve the player's name from the hash stored under "player:playerID"
	playerName, err := rr.Client.HGet("player:"+playerID, "name").Result()
	if err != nil {
		log.Println("Failed to get player name from Redis:", err)
		return player_score.PlayerScore{}, err
	}

	// Return the PlayerScore struct with PlayerID, PlayerName, and Score
	return player_score.PlayerScore{
		PlayerID:   playerID,
		PlayerName: playerName,
		Score:      score,
	}, nil
}

// InsertRecord inserts a player score into the Redis leaderboard and stores the playername in a HASH.
func (rr *RedisClient) InsertRecord(key, playerID, playername string, score float64) error {
	// Add the player score to the ZSET
	err := rr.Client.ZAdd(key, redis.Z{
		Score:  score,
		Member: playerID,
	}).Err()
	if err != nil {
		log.Println("Failed to insert score into Redis:", err)
		return err
	}

	// Store the player's playername in a HASH
	err = rr.Client.HSet("player:"+playerID, "playername", playername).Err()
	if err != nil {
		log.Println("Failed to store player details in Redis HASH:", err)
		return err
	}

	return nil
}

// Connect to Redis
func (rc *RedisClient) Connect() {
	client := redis.NewClient(&redis.Options{
		Addr:     rc.Addr,
		Password: rc.Password,
		DB:       rc.DB,
	})

	rc.Client = client

	if err := rc.Client.Ping().Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis!")
}

// Close the Redis connection
func (rc *RedisClient) Close() {
	rc.Client.Close()
}
