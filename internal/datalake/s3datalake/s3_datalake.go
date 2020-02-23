package s3datalake

import (
	"encoding/json"
	"log"
	"os"

	"github.com/MontgomeryWatts/SpotifyDBPlaylistTransformLambda/internal/datalake"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/zmb3/spotify"
)

type S3Datalake struct {
	Downloader *s3manager.Downloader
	Bucket     string
}

func NewS3Datalake() datalake.Datalake {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		log.Fatalf("Error occurred while initializing Session: %v", err)
	}

	bucket, ok := os.LookupEnv("BUCKET_NAME")
	if !ok {
		log.Fatalf("BUCKET_NAME not set in environment variables")
	}

	return &S3Datalake{
		Downloader: s3manager.NewDownloader(sess),
		Bucket:     bucket,
	}
}

func (datalake *S3Datalake) download(key string) []byte {
	var entity []byte
	buffer := aws.NewWriteAtBuffer(entity)
	input := &s3.GetObjectInput{
		Bucket: &datalake.Bucket,
		Key:    &key,
	}
	_, err := datalake.Downloader.Download(buffer, input)

	if err != nil {
		log.Fatalf("Error occurred while downloading entity from S3: %v", err)
	}
	return buffer.Bytes()
}

func (datalake *S3Datalake) GetArtist(artistKey []byte) spotify.FullArtist {
	keyString := string(artistKey)
	artist := spotify.FullArtist{}
	artistBytes := datalake.download(keyString)
	err := json.Unmarshal(artistBytes, &artist)
	if err != nil {
		log.Fatalf("Error unmarshalling artist %v", err)
	}
	return artist
}

func (datalake *S3Datalake) GetAlbum(albumKey []byte) spotify.FullAlbum {
	keyString := string(albumKey)
	album := spotify.FullAlbum{}
	albumBytes := datalake.download(keyString)
	err := json.Unmarshal(albumBytes, &album)
	if err != nil {
		log.Fatalf("Error unmarshalling album %v", err)
	}
	return album
}
