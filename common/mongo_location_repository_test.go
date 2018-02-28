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

func (d *MockDataLayer) C(name string) Collection {
	return &MockCollection{}
}

type MockCollection struct {
	store []interface{}
}

func (c *MockCollection) Find(query interface{}) Query {
	return &MockQuery{}
}

func (c *MockCollection) Insert(docs ...interface{}) error {
	c.store = docs
	return nil
}

type MockQuery struct {
	Queried bool
}

func (q *MockQuery) One(result interface{}) error {
	result.(*Location).SignatureKey = "test key"
	return nil
}

func TestRepository(t *testing.T) {
	session := MockSession{}
	repo := NewMongoLocationRepository(session, "test_database")

	location, err := repo.Find("id")
	if err != nil {
		t.Error("Expected to find without error")
	}
	if location.SignatureKey != "test key" {
		t.Errorf("Expected to receive mock record")
	}
}
