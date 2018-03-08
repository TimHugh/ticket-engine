package common

import (
	"fmt"
)

func NewMockAdapter() Adapter {
	return &MockAdapter{}
}

type MockAdapter struct {
	Doc interface{}
}

func (a *MockAdapter) Find(collection string, id string, result interface{}) error {
	switch value := result.(type) {
	case *Location:
		*value = a.Doc.(Location)
	case *Order:
		*value = a.Doc.(Order)
	default:
		return fmt.Errorf("Received unknown type %T", result)
	}
	return nil
}

func (a *MockAdapter) Create(collection string, doc interface{}) error {
	a.Doc = doc
	return nil
}

func (a *MockAdapter) Close() {}

type MockDocument struct {
	ID    string
	Value string
}
