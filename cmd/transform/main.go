package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/MontgomeryWatts/SpotifyDBPlaylistTransformLambda/internal/datalake/datalakelayer"
	"github.com/MontgomeryWatts/SpotifyDBPlaylistTransformLambda/internal/db/mongodatabase"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(evt events.SQSEvent) {
	datalake := datalakelayer.NewDatalake(datalakelayer.S3)
	database := mongodatabase.NewMongoDatabase()

	for _, record := range evt.Records {
		s3Event := events.S3Event{}
		err := json.Unmarshal([]byte(record.Body), &s3Event)
		if err != nil {
			log.Fatal("Malformed SQSEvent received")
		}
		obj := &s3Event.Records[0].S3.Object
		key := obj.Key
		keyBytes := []byte(key)
		if strings.HasPrefix(key, "artists") {
			artist := datalake.GetArtist(keyBytes)
			err := database.InsertArtist(artist)
			if err != nil {
				log.Fatalf("Error inserting artist into database: %v", err)
			}
		} else if strings.HasPrefix(key, "albums") {
			album := datalake.GetAlbum(keyBytes)
			err := database.InsertAlbum(album)
			if err != nil {
				log.Fatalf("Error inserting album into database: %v", err)
			}
		} else {
			log.Printf("Unrecognized entity type encountered: %s", key)
		}
	}
}

func main() {
	lambda.Start(Handler)
}
