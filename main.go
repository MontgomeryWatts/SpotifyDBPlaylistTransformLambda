package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/zmb3/spotify"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func getEnv(key string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("%s is not set in environment variables", key)
	}
	return env
}

func Handler(evt events.S3Event) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		log.Fatalf("Error occurred while initializing Session: %v", err)
	}

	downloader := s3manager.NewDownloader(sess)
	c := make(chan bool)
	receives := len(evt.Records)

	for _, record := range evt.Records {
		s3Entity := record.S3
		bucket := &s3Entity.Bucket
		obj := &s3Entity.Object
		if strings.HasPrefix(obj.Key, "artists") {
			go processArtist(c, downloader, bucket, obj)
		} else if strings.HasPrefix(obj.Key, "albums") {
			go processAlbum(c, downloader, bucket, obj)
		} else {
			log.Fatalf("Unrecognized entity type encountered: %s", obj.Key)
		}
	}

	for i := 0; i < receives; i++ {
		<-c
	}
}

func downloadEntity(downloader *s3manager.Downloader, bucket, key *string) []byte {
	if downloader == nil {
		log.Fatal("Nil pointer passed to downloadEntity for s3manager.Downloader")
	}

	if bucket == nil {
		log.Fatal("Nil pointer passed to downloadEntity for bucket string")
	}

	if key == nil {
		log.Fatal("Nil pointer passed to downloadEntity for key string")
	}

	var entity []byte
	buffer := aws.NewWriteAtBuffer(entity)
	input := &s3.GetObjectInput{
		Bucket: bucket,
		Key:    key,
	}
	_, err := downloader.Download(buffer, input)

	if err != nil {
		log.Fatalf("Error occurred while downloading entity from S3: %v", err)
	}
	return buffer.Bytes()
}

func processArtist(c chan bool, downloader *s3manager.Downloader, bucket *events.S3Bucket, object *events.S3Object) {
	artistBytes := downloadEntity(downloader, &bucket.Name, &object.Key)
	var artist spotify.FullArtist
	err := json.Unmarshal(artistBytes, &artist)
	if err != nil {
		log.Fatalf("Error unmarshalling artist: %v", err)
	}
	log.Print(artist)
	c <- true
}

func processAlbum(c chan bool, downloader *s3manager.Downloader, bucket *events.S3Bucket, object *events.S3Object) {
	albumBytes := downloadEntity(downloader, &bucket.Name, &object.Key)
	var album spotify.FullAlbum
	err := json.Unmarshal(albumBytes, &album)
	if err != nil {
		log.Fatalf("Error unmarshalling album: %v", err)
	}
	log.Print(album)
	c <- true
}

func main() {
	lambda.Start(Handler)
}
