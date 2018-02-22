package main

import (
	"testing"
)

func TestCreatesNewOrders(t *testing.T) {
	orderRepo := &MockOrderRepository{}
	event := Event{
		OrderID:    "order_id",
		LocationID: "location_id",
	}

	handler := PaymentUpdateHandler{
		OrderCreator{orderRepo},
	}

	err := handler.Handle(event)
	if err != nil {
		t.Error("Expected to successfully create a new order.")
	}

	order := orderRepo.Order
	if order.ID != event.OrderID {
		t.Errorf("Expected new order with ID %s but got %s", event.OrderID, order.ID)
	}
	if order.LocationID != event.LocationID {
		t.Errorf("Expected new order with locationID %s but got %s", event.LocationID, order.LocationID)
	}
}
