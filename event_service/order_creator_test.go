package event_service

import (
	"testing"

	"strings"
)

type MockOrderRepository struct {
	IncludesDuplicate bool
	Order             Order
}

func (m *MockOrderRepository) Store(order Order) {
	m.Order = order
}

func (m *MockOrderRepository) Find(id string) *Order {
	if m.IncludesDuplicate {
		return &Order{}
	} else {
		return nil
	}
}

func TestRejectsDuplicates(t *testing.T) {
	orderCreator := OrderCreator{
		&MockOrderRepository{
			IncludesDuplicate: true,
		},
	}

	order := Order{
		ID: "order_id",
	}

	err := orderCreator.Create(order)
	if err == nil {
		t.Error("Expected duplicate order to fail creating.")
	}
	if !strings.Contains(err.Error(), "duplicate order") {
		t.Errorf("Expected error to include 'duplicate order' but received %s", err.Error())
	}
}

func TestStoresValidOrder(t *testing.T) {
	mockRepo := &MockOrderRepository{
		IncludesDuplicate: false,
	}

	orderCreator := OrderCreator{
		OrderRepository: mockRepo,
	}

	order := Order{
		ID: "order_id",
	}

	if err := orderCreator.Create(order); err != nil {
		t.Error("Expected new order to create successfully.")
	}

	if mockRepo.Order != order {
		t.Error("Expected create to store new order in repo.")
	}
}
