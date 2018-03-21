package mongo

import (
	"testing"

	"github.com/timhugh/ticket_service/root"
)

func TestFind(t *testing.T) {
}

func TestCreateNew(t *testing.T) {
	s := &MockSession{}
	o := &OrderRepository{s}

	s.InsertFn = func(docs ...interface{}) error {
		return nil
	}

	order := root.Order{
		ID:         "id",
		LocationID: "location_id",
	}

	if err := o.Create(order); err != nil {
		t.Errorf("expected nil error but got %s", err)
	}
	// if s.InsertInvoked != true {
	// 	t.Errorf("expected insert method to be invoked")
	// }
}
