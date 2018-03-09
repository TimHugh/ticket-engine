package common

import (
	"testing"

	"fmt"
)

type mockOrderAdapter struct {
	Order *Order
}

func (a mockOrderAdapter) Find(collection string, id string, order interface{}) error {
	if id == "good_id" && collection == "orders" {
		pOrder := order.(*Order)
		pOrder.ID = "good_id"
		pOrder.LocationID = "test_location"
		return nil
	} else {
		return fmt.Errorf("not found")
	}
}

func (a *mockOrderAdapter) Create(collection string, order interface{}) error {
	value := order.(Order)
	a.Order = &value
	return nil
}

func (a mockOrderAdapter) Close() {}

func TestOrderRepoFind(t *testing.T) {
	adapter := &mockOrderAdapter{}
	repo := OrderRepository{adapter}

	goodOrder, goodErr := repo.Find("good_id")
	if goodErr != nil {
		t.Errorf("Expected to find without error but received %s", goodErr)
	}
	if goodOrder.ID != "good_id" {
		t.Errorf("Expected to find order with id 'good_id' but got %s", goodOrder.ID)
	}
	if goodOrder.LocationID != "test_location" {
		t.Errorf("Expected to find order with location_id 'test_location' but got %s", goodOrder.LocationID)
	}

	badOrder, badErr := repo.Find("bad_id")
	if badErr == nil {
		t.Errorf("Expected 'not found' error but received none.")
	}
	if badOrder != nil {
		t.Errorf("Expected missing order to be nil but got %s", badOrder)
	}
}

func TestOrderRepoStore(t *testing.T) {
	adapter := &mockOrderAdapter{}
	repo := OrderRepository{adapter}

	newOrder := Order{
		ID:         "good_id",
		LocationID: "location_id",
	}
	storeErr := repo.Store(newOrder)
	if storeErr != nil {
		t.Errorf("Expected to store without error but recieved: %s", storeErr)
	}
	if adapter.Order.ID != "good_id" {
		t.Errorf("Adapter did not receive order to store")
	}
}
