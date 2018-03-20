package main

import (
	"testing"

	"fmt"

	"github.com/timhugh/ticket_service/mock"
	"github.com/timhugh/ticket_service/root"
)

func TestCreatesNewOrders(t *testing.T) {
	r := &mock.OrderRepository{}
	h := NewPaymentUpdateHandler(r)

	r.FindFn = func(locationID, id string) (*root.Order, error) {
		return nil, fmt.Errorf("not found")
	}
	r.CreateFn = func(root.Order) error {
		return nil
	}

	event := Event{
		OrderID:    "id",
		LocationID: "id",
	}

	if err := h.Handle(event); err != nil {
		t.Errorf("Expected to create order without error but got %s", err)
	}
}

func TestDoesNotCreateDuplicates(t *testing.T) {
	r := &mock.OrderRepository{}
	h := NewPaymentUpdateHandler(r)

	r.FindFn = func(locationID, id string) (*root.Order, error) {
		return &root.Order{}, nil
	}

	event := Event{
		OrderID:    "id",
		LocationID: "id",
	}

	if err := h.Handle(event); err == nil {
		t.Error("Expected to receive error for creating a duplicate order.")
	}
}
