package mongo

import ()

type MockSession struct {
	InsertFn      func(...interface{}) error
	InsertInvoked bool

	OneFn      func(interface{}) error
	OneInvoked bool
}

func (s *MockSession) DB(name string) DataLayer {
	return &MockDataLayer{s}
}

func (s *MockSession) Clone() Session {
	return s
}

func (s *MockSession) Close() {}

type MockDataLayer struct {
	Session *MockSession
}

func (d *MockDataLayer) C(name string) Collection {
	return &MockCollection{d.Session}
}

type MockCollection struct {
	Session *MockSession
}

func (c *MockCollection) Find(query interface{}) Query {
	return &MockQuery{c.Session}
}

func (c *MockCollection) Insert(docs ...interface{}) error {
	c.Session.InsertInvoked = true
	return c.Session.InsertFn(docs...)
}

type MockQuery struct {
	Session *MockSession
}

func (q *MockQuery) One(result interface{}) error {
	q.Session.OneInvoked = true
	return q.Session.OneFn(result)
}
