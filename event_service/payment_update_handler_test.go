package event_service

import (
	"testing"
)

type MockLocationRepository struct {
	Location *Location
}

func (m *MockLocationRepository) Find(locationID string) *Location {
	return m.Location
}

func (m *MockLocationRepository) Store(location Location) {
	m.Location = &location
}

func TestRejectsUnknownLocations(t *testing.T) {
	locationRepo := &MockLocationRepository{}
	event := Event{
		LocationID: "location_id",
	}

	handler := PaymentUpdateHandler{
		LocationRepository: locationRepo,
	}

	err := handler.Handle(event)
	if err == nil {
		t.Error("Expected exception for unknown location")
	}
}

func TestCreatesNewOrders(t *testing.T) {
	orderRepo := &MockOrderRepository{}
	event := Event{
		OrderID:    "order_id",
		LocationID: "location_id",
	}
	locationRepo := &MockLocationRepository{
		Location: &Location{
			ID: event.LocationID,
		},
	}

	handler := PaymentUpdateHandler{
		OrderCreator: OrderCreator{
			OrderRepository: orderRepo,
		},
		LocationRepository: locationRepo,
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
