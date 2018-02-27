package common

import (
	"testing"
)

type MockSession struct{}

func (s MockSession) DB(name string) DataLayer {
	return &MockDataLayer{}
}

func (s *MockSession) Close() {}

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
	result = Result{Data: "result"}
	return nil
}

type Result struct {
	Data string
}

func TestAdapter(t *testing.T) {
	session := &MockSession{}
	adapter := NewMongoAdapter(session, "test_database")

	doc, err := adapter.Find("collection", "query")
	if err != nil {
		t.Error("Expected to find without error")
	}
	doc = Result(*doc)
	if doc.Data != "result" {
		t.Error("Expected to retrieve mock record 'result'")
	}
}
