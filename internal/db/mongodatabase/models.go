package mongodatabase

import "github.com/zmb3/spotify"

type MongoArtist struct {
	URI    string   `bson:"_id"`
	ID     string   `bson:"id"`
	Name   string   `bson:"name"`
	Genres []string `bson:"genres,omitempty"`
}

type MongoTrack struct {
	ID string `bson:"_id"`
}

func NewMongoArtist(artist spotify.FullArtist) MongoArtist {
	return MongoArtist{
		URI:    string(artist.URI),
		ID:     string(artist.ID),
		Name:   artist.Name,
		Genres: artist.Genres,
	}
}
