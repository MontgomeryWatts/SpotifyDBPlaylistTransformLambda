package main

import (
	"log"
	"strings"

	"github.com/MontgomeryWatts/SpotifyDBPlaylistTransformLambda/internal/datalake/datalakelayer"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(evt events.S3Event) {

	datalake := datalakelayer.NewDatalake(datalakelayer.S3)

	for _, record := range evt.Records {
		s3Entity := record.S3
		obj := &s3Entity.Object
		key := obj.Key
		keyBytes := []byte(key)
		if strings.HasPrefix(key, "artists") {
			datalake.GetArtist(keyBytes)
		} else if strings.HasPrefix(key, "albums") {
			datalake.GetAlbum(keyBytes)
		} else {
			log.Fatalf("Unrecognized entity type encountered: %s", key)
		}
	}
}

func main() {
	lambda.Start(Handler)
}
