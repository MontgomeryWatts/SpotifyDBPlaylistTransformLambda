package db

import (
	"github.com/zmb3/spotify"
)

type Database interface {
	InsertArtist(spotify.FullArtist) error
	InsertAlbum(spotify.FullAlbum) error
}
