package main

import (
	"encoding/json"

	"github.com/MontgomeryWatts/SpotifyDBPlaylistTransformLambda/internal/db"
	"github.com/MontgomeryWatts/SpotifyDBPlaylistTransformLambda/internal/db/mongodatabase"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type TestEvent struct {
	Field string `json:"yeet"`
}

type Placeholder struct {
	IDs []string `json:"ids"`
}

func Handler(evt TestEvent) (events.APIGatewayProxyResponse, error) {
	var database db.Database = mongodatabase.NewMongoDatabase()
	body := Placeholder{IDs: database.GetPlaylistTracks()}
	b, err := json.Marshal(body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			IsBase64Encoded: false,
			StatusCode:      500,
		}, err
	}
	return events.APIGatewayProxyResponse{
		IsBase64Encoded: false,
		StatusCode:      200,
		Body:            string(b),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
