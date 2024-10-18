package repositories

import (
	"context"
	"log"
	"quiz/internals/domain/player_score"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBClient handles the connection to MongoDB and operations related to player scores.
type MongoDBClient struct {
	Ctx    context.Context // Context for managing database operations
	URI    string          // MongoDB connection URI
	Client *mongo.Client   // MongoDB client instance
}

// NewMongoDBClient creates a new instance of MongoDBClient with the given context and URI.
func NewMongoDBClient(ctx context.Context, uri string) *MongoDBClient {
	return &MongoDBClient{Ctx: ctx, URI: uri}
}

// UpdateOrInsertPlayerScore updates a player's score if it exists, or inserts it if it doesn't.
func (mdb *MongoDBClient) UpdateOrInsertPlayerScore(player player_score.PlayerScore) error {
	collection := mdb.Client.Database("game").Collection("players")
	_, err := collection.UpdateOne(
		mdb.Ctx,
		bson.M{"player_id": player.PlayerID}, // Filter by player ID
		bson.M{"$set": bson.M{"score": player.Score, "player_name": player.PlayerName}}, // Update player score and name
		options.Update().SetUpsert(true),                                                // Insert new document if none exists
	)
	if err != nil {
		log.Println("Failed to update player score in MongoDB:", err)
		return err
	}
	return nil
}

// GetTopPlayers retrieves the top players sorted by score in descending order.
func (mdb *MongoDBClient) GetTopPlayers() ([]player_score.PlayerScore, error) {
	collection := mdb.Client.Database("game").Collection("players")
	cursor, err := collection.Find(mdb.Ctx, bson.D{}, options.Find().SetSort(bson.M{"score": -1})) // Sort by score (highest first)
	if err != nil {
		log.Println("Failed to get top players from MongoDB:", err)
		return nil, err
	}
	defer cursor.Close(mdb.Ctx)

	var topPlayers []player_score.PlayerScore
	for cursor.Next(mdb.Ctx) {
		var player player_score.PlayerScore
		if err := cursor.Decode(&player); err != nil {
			log.Println("Failed to decode player data:", err)
			return nil, err
		}
		topPlayers = append(topPlayers, player)
	}

	return topPlayers, nil
}

// GetPlayerScore retrieves the score of a specific player by their ID.
func (mdb *MongoDBClient) GetPlayerScore(playerID string) (int, error) {
	collection := mdb.Client.Database("game").Collection("players")
	var result player_score.PlayerScore
	err := collection.FindOne(mdb.Ctx, bson.M{"player_id": playerID}).Decode(&result)
	if err != nil {
		return 0, err
	}
	return result.Score, nil
}

// Connect establishes a connection to MongoDB using the provided URI.
func (mc *MongoDBClient) Connect() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mc.URI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(mc.Ctx, opts)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	mc.Client = client
	log.Println("Connected to MongoDB!")
}

// Close gracefully closes the connection to MongoDB.
func (mc *MongoDBClient) Close() {
	mc.Client.Disconnect(mc.Ctx)
}
