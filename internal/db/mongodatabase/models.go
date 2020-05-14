package mongodatabase

import (
	"github.com/zmb3/spotify"
)

type MongoTrack struct {
	URI      string   `bson:"_id"`
	Artists  []string `bson:"artists"`
	Duration int      `bson:"duration_ms"`
}

func NewMongoTracks(album spotify.FullAlbum) []interface{} {
	tracks := make([]interface{}, len(album.Tracks.Tracks))
	for i, track := range album.Tracks.Tracks {
		t := MongoTrack{
			URI:      string(track.URI),
			Duration: track.Duration,
		}
		artists := make([]string, len(track.Artists))
		for j, artist := range track.Artists {
			artists[j] = string(artist.URI)
		}
		t.Artists = artists
		tracks[i] = t
	}
	return tracks
}
