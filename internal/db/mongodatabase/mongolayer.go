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
	trackCollection  = "tracks"
)

type MongoDatabase struct {
	client *mongo.Client
}

type TrackDocument struct {
	URI string `bson:"_id"`
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
		bson.E{Key: "_id", Value: string(artist.URI)}}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "name", Value: artist.Name},
			bson.E{Key: "genres", Value: artist.Genres},
		}}}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (mg *MongoDatabase) InsertTracks(album spotify.FullAlbum) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mg.client.Database(database).Collection(trackCollection)
	for _, t := range album.Tracks.Tracks {
		artists := make([]string, len(t.Artists))
		for i, a := range t.Artists {
			artists[i] = string(a.URI)
		}

		filter := bson.D{
			bson.E{Key: "_id", Value: string(t.URI)}}
		update := bson.D{
			bson.E{Key: "$set", Value: bson.D{
				bson.E{Key: "artists", Value: artists},
				bson.E{Key: "duration_ms", Value: t.Duration},
				bson.E{Key: "explicit", Value: t.Explicit},
			}}}
		opts := options.Update().SetUpsert(true)
		_, err := collection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mg *MongoDatabase) GetPlaylistTracks() []string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sampleStage := bson.D{
		bson.E{Key: "$sample", Value: bson.D{
			bson.E{Key: "size", Value: 100},
		}},
	}
	projectStage := bson.D{
		bson.E{Key: "$project", Value: bson.D{
			bson.E{Key: "_id", Value: 1},
		}},
	}

	collection := mg.client.Database(database).Collection(trackCollection)
	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{sampleStage, projectStage})
	if err != nil {
		return make([]string, 0)
	}

	var trackDocs []TrackDocument
	if err = cursor.All(ctx, &trackDocs); err != nil {
		return make([]string, 0)
	}

	ids := make([]string, len(trackDocs))
	for i, doc := range trackDocs {
		ids[i] = doc.URI[14:]
	}

	return ids
}
