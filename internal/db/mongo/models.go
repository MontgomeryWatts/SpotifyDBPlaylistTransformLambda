package mongo

type MongoArtist struct {
	ID     string   `bson:"_id"`
	URI    string   `bson:"uri"`
	Name   string   `bson:"name"`
	Genres []string `bson:"genres"`
}

type MongoTrack struct {
	ID string `bson:"_id"`
}
