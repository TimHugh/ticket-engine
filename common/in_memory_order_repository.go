package common

import (
	"errors"
)

type InMemoryOrderRepository struct {
	memoryStore map[string]Order
}

func NewInMemoryOrderRepository() OrderRepository {
	return &InMemoryOrderRepository{
		make(map[string]Order),
	}
}

func (repo *InMemoryOrderRepository) Find(id string) (*Order, error) {
	order := repo.memoryStore[id]
	if order.ID == "" {
		return nil, errors.New("Order does not exist")
	}
	return &order, nil
}

func (repo *InMemoryOrderRepository) Store(order Order) error {
	repo.memoryStore[order.ID] = order
	return nil
}
