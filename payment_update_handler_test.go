package main

import (
	"testing"

	"fmt"

	"github.com/timhugh/ticket_service/common"
)

type mockOrderRepository struct {
	Order *common.Order
}

func (r mockOrderRepository) Find(id string) (*common.Order, error) {
	return r.Order, nil
}

func (r *mockOrderRepository) Store(order common.Order) error {
	if r.Order == nil {
		r.Order = &order
		return nil
	} else {
		return fmt.Errorf("duplicate")
	}
}

func TestCreatesNewOrders(t *testing.T) {
	orderRepo := &mockOrderRepository{}
	event := Event{
		OrderID:    "order_id",
		LocationID: "location_id",
	}

	handler := NewPaymentUpdateHandler(orderRepo)

	if err := handler.Handle(event); err != nil {
		t.Errorf("Expected to create order without error but got %s", err)
	}

	// duplicate error check
	if err := handler.Handle(event); err == nil {
		t.Error("Expected to receive error for creating a duplicate order.")
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
