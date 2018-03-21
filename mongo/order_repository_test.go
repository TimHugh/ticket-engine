package mongo

import (
	"testing"

	"fmt"

	"github.com/timhugh/ticket_service/root"
)

func TestFind(t *testing.T) {
	s := &MockSession{}
	o := OrderRepository{s}

	s.OneFn = func(out interface{}) error {
		result := out.(*root.Order)
		result.ID = "order"
		result.LocationID = "location"
		return nil
	}

	order, err := o.Find("id", "id")
	if err != nil {
		t.Errorf("expected nil error but got %s", err)
	}
	if order.ID != "order" || order.LocationID != "location" {
		t.Errorf(
			"expected order with id 'order' and location id 'location' but got id '%s' and location id '%s'",
			order.ID, order.LocationID,
		)
	}
}

func TestNotFound(t *testing.T) {
	s := &MockSession{}
	o := OrderRepository{s}

	s.OneFn = func(out interface{}) error {
		return fmt.Errorf("not found")
	}

	_, err := o.Find("id", "id")
	if err == nil {
		t.Errorf("expected error but received none")
	}
}

func TestCreateNew(t *testing.T) {
	s := &MockSession{}
	o := OrderRepository{s}

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
	if s.InsertInvoked != true {
		t.Error("expected insert method to be invoked")
	}
}
