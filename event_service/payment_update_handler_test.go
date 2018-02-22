package main

import (
	"testing"

	"github.com/timhugh/ticket-engine/common"
)

func TestCreatesNewOrders(t *testing.T) {
	orderRepo := common.NewInMemoryOrderRepository()
	event := Event{
		OrderID:    "order_id",
		LocationID: "location_id",
	}

	handler := NewPaymentUpdateHandler(orderRepo)

	err := handler.Handle(event)
	if err != nil {
		t.Error("Expected to successfully create a new order.")
	}

	order := orderRepo.Find("order_id")
	if order.ID != event.OrderID {
		t.Errorf("Expected new order with ID %s but got %s", event.OrderID, order.ID)
	}
	if order.LocationID != event.LocationID {
		t.Errorf("Expected new order with locationID %s but got %s", event.LocationID, order.LocationID)
	}
}
