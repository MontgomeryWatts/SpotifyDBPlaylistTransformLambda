package datalakelayer

import (
	"github.com/MontgomeryWatts/SpotifyDBPlaylistTransformLambda/internal/datalake"
	"github.com/MontgomeryWatts/SpotifyDBPlaylistTransformLambda/internal/datalake/s3datalake"
)

type DatalakeType string

const (
	S3 DatalakeType = "s3"
)

func NewDatalake(datalakeType DatalakeType) datalake.Datalake {
	switch datalakeType {
	case S3:
		return s3datalake.NewS3Datalake()
	}
	return nil
}
