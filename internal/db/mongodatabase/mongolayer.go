package mongodatabase

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/MontgomeryWatts/SpotifyDBPlaylistTransformLambda/internal/db"
	"github.com/zmb3/spotify"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	database         = "music"
	artistCollection = "artists"
	albumCollection  = "albums"
)

type MongoDatabase struct {
	client *mongo.Client
}

func NewMongoDatabase() db.Database {
	connectionString, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		log.Fatalf("MONGODB_URI not set in environment variables")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error while initializing MongoDB Client: %v", err)
	}

	return &MongoDatabase{
		client: client,
	}
}

func (mg *MongoDatabase) InsertArtist(artist spotify.FullArtist) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := mg.client.Database(database).Collection(artistCollection)

	filter := bson.D{
		bson.E{Key: "_id", Value: string(artist.ID)}}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "name", Value: artist.Name},
			bson.E{Key: "image", Value: artist.Images[0].URL},
			bson.E{Key: "uri", Value: string(artist.URI)},
			bson.E{Key: "genres", Value: artist.Genres},
		}}}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (mg *MongoDatabase) InsertAlbum(album spotify.FullAlbum) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mg.client.Database(database).Collection(albumCollection)

	artists := make([]bson.D, len(album.Artists))
	for index, artist := range album.Artists {
		artists[index] = bson.D{
			bson.E{Key: "name", Value: artist.Name},
			bson.E{Key: "uri", Value: artist.URI},
		}
	}

	filter := bson.D{
		bson.E{Key: "_id", Value: string(album.ID)}}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "name", Value: album.Name},
			bson.E{Key: "image", Value: album.Images[0].URL},
			bson.E{Key: "uri", Value: album.URI},
			bson.E{Key: "artists", Value: artists},
		}}}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}
