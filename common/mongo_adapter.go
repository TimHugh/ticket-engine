package common

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoAdapter satisfies Adapter interface
func NewMongoAdapter(host_uri string) (Adapter, error) {
	session, err := NewMongoSession(host_uri)
	return MongoAdapter{session}, err
}

type MongoAdapter struct {
	session Session
}

func (m MongoAdapter) collection(name string) (Session, Collection) {
	session := m.session.Clone()
	collection := session.DB("").C(name)
	return session, collection
}

func (m MongoAdapter) Find(collection_name string, id string, result interface{}) error {
	session, collection := m.collection(collection_name)
	defer session.Close()

	return collection.Find(bson.M{"id": id}).One(result)
}

func (m MongoAdapter) Create(collection_name string, doc interface{}) error {
	session, collection := m.collection(collection_name)
	defer session.Close()

	return collection.Insert(doc)
}

func (m MongoAdapter) Close() {
	m.session.Close()
}

// MongoSession wraps mgo.Session
type Session interface {
	DB(name string) DataLayer
	Clone() Session
	Close()
}

func NewMongoSession(host string) (Session, error) {
	session, err := mgo.Dial(host)
	return MongoSession{session}, err
}

type MongoSession struct {
	*mgo.Session
}

func (s MongoSession) DB(name string) DataLayer {
	return MongoDatabase{s.Session.DB(name)}
}

func (s MongoSession) Clone() Session {
	return MongoSession{s.Session.Clone()}
}

func (s MongoSession) Close() {
	s.Session.Close()
}

// MongoDatabase wraps mgo.Database
type DataLayer interface {
	C(name string) Collection
}

type MongoDatabase struct {
	*mgo.Database
}

func (d MongoDatabase) C(name string) Collection {
	return MongoCollection{d.Database.C(name)}
}

// MongoCollection wraps mgo.Collection
type Collection interface {
	Find(query interface{}) Query
	Insert(docs ...interface{}) error
}

type MongoCollection struct {
	*mgo.Collection
}

func (c MongoCollection) Find(query interface{}) Query {
	return MongoQuery{c.Collection.Find(query)}
}

func (c MongoCollection) Insert(docs ...interface{}) error {
	return c.Collection.Insert(docs)
}

// MongoQuery wraps mgo.Query
type Query interface {
	One(result interface{}) error
}

type MongoQuery struct {
	*mgo.Query
}

func (q MongoQuery) One(result interface{}) error {
	return q.Query.One(result)
}
