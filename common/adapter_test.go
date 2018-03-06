package common

import (
	"fmt"
)

func NewMockAdapter() Adapter {
	return &MockAdapter{}
}

type MockAdapter struct {
	doc Document
}

func (a *MockAdapter) Find(collection string, id string, result interface{}) error {
	switch result.(type) {
	case *Location:
		value := result.(*Location)
		value.ID = id
		value.SignatureKey = "signature"
	case *Order:
		value := result.(*Order)
		value.ID = id
	default:
		return fmt.Errorf("Received unknown type %T", result)
	}
	return nil
}

func (a *MockAdapter) Create(collection string, doc interface{}) error {
	a.doc = doc.(Document)
	return nil
}

func (a *MockAdapter) Close() {}
