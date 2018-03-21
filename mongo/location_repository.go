package mongo

import (
	"gopkg.in/mgo.v2/bson"

	root "github.com/timhugh/ticket_service"
)

type LocationRepository struct {
	Session Session
}

func (s LocationRepository) collection() (Session, Collection) {
	session := s.Session.Clone()
	collection := session.DB("").C("locations")
	return session, collection
}

func (s LocationRepository) Create(location root.Location) error {
	session, collection := s.collection()
	defer session.Close()

	return collection.Insert(location)
}

func (s LocationRepository) Find(id string) (*root.Location, error) {
	session, collection := s.collection()
	defer session.Close()

	location := &root.Location{}
	err := collection.Find(bson.M{"id": id}).One(location)
	return location, err
}
