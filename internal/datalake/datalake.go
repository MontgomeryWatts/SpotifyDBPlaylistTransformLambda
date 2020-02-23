package datalake

import (
	"github.com/zmb3/spotify"
)

type Datalake interface {
	GetArtist([]byte) spotify.FullArtist
	GetAlbum([]byte) spotify.FullAlbum
}
