package common

import (
	"gopkg.in/mgo.v2/bson"
)

type MongoLocationRepository struct {
	session         Session
	database_name   string
	collection_name string
}

func NewMongoLocationRepository(session Session, database string) LocationRepository {
	return MongoLocationRepository{session, database, "locations"}
}

func (r MongoLocationRepository) collection() (Session, Collection) {
	session := r.session.Clone()
	collection := session.DB(r.database_name).C(r.collection_name)
	return session, collection
}

func (r MongoLocationRepository) Find(id string) (*Location, error) {
	session, collection := r.collection()
	defer session.Close()

	var location Location
	if err := collection.Find(bson.M{"id": id}).One(&location); err != nil {
		return nil, err
	} else {
		return &location, err
	}
}

func (r MongoLocationRepository) Store(location Location) error {
	return nil
}
