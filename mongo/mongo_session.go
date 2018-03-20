package mongo

import (
	"gopkg.in/mgo.v2"
)

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
	return c.Collection.Insert(docs...)
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
