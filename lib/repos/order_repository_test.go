package repos

import (
	"testing"
)

type MockOrderRepository struct {
	Order Order
}

func (m *MockOrderRepository) Store(order Order) {
	m.Order = order
}

func (m *MockOrderRepository) Find(id string) *Order {
	return m.Order
}
