package mongo

import (
	"testing"

	"github.com/timhugh/ticket_service/root"
)

func TestFind(t *testing.T) {
	var s MockSession

	var o OrderRepository
	o.Session = s
}

func TestCreateNew(t *testing.T) {
	var s MockSession
	s.InsertFn = func(docs ...interface{}) error {
		return nil
	}

	var o OrderRepository
	o.Session = s

	order := root.Order{
		ID:         "id",
		LocationID: "location_id",
	}

	if err := o.Create(order); err != nil {
		t.Errorf("expected nil error but got %s", err)
	}

	// TODO: finish this
}
