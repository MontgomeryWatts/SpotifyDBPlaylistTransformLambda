package mongo

import (
	"github.com/MontgomeryWatts/SpotifyDBPlaylistTransformLambda/internal/db"
	"github.com/zmb3/spotify"
)

type MongoDatabase struct {
}

func NewMongoDatabase(connection string) (db.Database, error) {
	return nil, nil
}

func (mongo *MongoDatabase) InsertArtist(artist spotify.FullArtist) error {
	return nil
}

func (mongo *MongoDatabase) InsertTracks(album spotify.FullAlbum) error {
	return nil
}
