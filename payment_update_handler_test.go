package main

import (
	"testing"

	"github.com/timhugh/ticket_service/common"
)

type mockOrderRepository struct {
	Order *common.Order
}

func (r mockOrderRepository) Find(id string) (*common.Order, error) {
	return r.Order, nil
}

func (r *mockOrderRepository) Store(order common.Order) error {
	r.Order = &order
	return nil
}

func TestCreatesNewOrders(t *testing.T) {
	orderRepo := &mockOrderRepository{}
	event := Event{
		OrderID:    "order_id",
		LocationID: "location_id",
	}

	handler := NewPaymentUpdateHandler(orderRepo)

	err := handler.Handle(event)
	if err != nil {
		t.Error("Expected to successfully create a new order.")
	}

	order, err := orderRepo.Find("order_id")
	if err != nil {
		t.Error("Expected to successfully find order")
	}
	if order.ID != event.OrderID {
		t.Errorf("Expected new order with ID %s but got %s", event.OrderID, order.ID)
	}
	if order.LocationID != event.LocationID {
		t.Errorf("Expected new order with locationID %s but got %s", event.LocationID, order.LocationID)
	}
}
