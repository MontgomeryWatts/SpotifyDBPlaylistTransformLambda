package mongodatabase

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/MontgomeryWatts/SpotifyDBPlaylistTransformLambda/internal/db"
	"github.com/zmb3/spotify"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	database   = "music"
	collection = "playlist"
)

type MongoDatabase struct {
	client *mongo.Client
}

func NewMongoDatabase() db.Database {
	connectionString, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		log.Fatalf("MONGODB_URI not set in environment variables")
	}

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Error while initializing MongoDB Client: %v", err)
	}

	return &MongoDatabase{
		client: client,
	}
}

func (mongo *MongoDatabase) InsertArtist(artist spotify.FullArtist) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	mongoArtist := NewMongoArtist(artist)
	collection := mongo.client.Database(database).Collection(collection)
	_, err := collection.InsertOne(ctx, mongoArtist)
	return err
}

func (mongo *MongoDatabase) InsertTracks(album spotify.FullAlbum) error {
	return nil
}
