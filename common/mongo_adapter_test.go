package common

import (
	"testing"
)

type MockSession struct{}

func (s MockSession) DB(name string) DataLayer {
	return &MockDataLayer{}
}

func (s MockSession) Clone() Session {
	return s
}

func (s MockSession) Close() {}

type MockDataLayer struct{}

func (d MockDataLayer) C(name string) Collection {
	return &MockCollection{}
}

type MockCollection struct{}

func (c MockCollection) Find(query interface{}) Query {
	return &MockQuery{}
}

func (c MockCollection) Insert(docs ...interface{}) error {
	return nil
}

type MockQuery struct {
	Queried bool
}

func (q MockQuery) One(result interface{}) error {
	return nil
}

func TestAdapter(t *testing.T) {
	session := MockSession{}
	adapter := MongoAdapter{session}

	doc := MockDocument{
		ID:    "id",
		Value: "value",
	}
	err := adapter.Create("collection", doc)
	if err != nil {
		t.Errorf("Expected to create document without error but received %s", err)
	}
}
