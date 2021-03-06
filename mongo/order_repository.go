package mongo

import (
	"gopkg.in/mgo.v2/bson"

	root "github.com/timhugh/ticket_service"
)

var _ root.OrderRepository = OrderRepository{}

type OrderRepository struct {
	Session Session
}

func (s OrderRepository) collection() (Session, Collection) {
	session := s.Session.Clone()
	collection := session.DB("").C("orders")
	return session, collection
}

func (s OrderRepository) Create(order root.Order) error {
	session, collection := s.collection()
	defer session.Close()

	return collection.Insert(order)
}

func (s OrderRepository) Find(id string, locationID string) (*root.Order, error) {
	session, collection := s.collection()
	defer session.Close()

	var order root.Order
	err := collection.Find(bson.M{"id": id, "location_id": locationID}).One(&order)
	return &order, err
}
