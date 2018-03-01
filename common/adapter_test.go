package common

import ()

func NewMockAdapter() Adapter {
	return &MockAdapter{}
}

type MockAdapter struct {
	doc Document
}

func (a *MockAdapter) Find(collection string, id string, result *interface{}) error {
	switch result.(type) {
	case *Location:
		*result = Location{ID: id}
	case *Order:
		*result = Order{ID: id}
	}
	return nil
}

func (a *MockAdapter) Create(collection string, doc interface{}) error {
	a.doc = doc.(Document)
	return nil
}

func (a *MockAdapter) Close() {}
