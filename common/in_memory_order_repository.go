package common

import ()

type InMemoryOrderRepository struct {
	memoryStore map[string]Order
}

func NewInMemoryOrderRepository() OrderRepository {
	return &InMemoryOrderRepository{
		make(map[string]Order),
	}
}

func (repo *InMemoryOrderRepository) Find(id string) *Order {
	order := repo.memoryStore[id]
	if order.ID == "" {
		return nil
	}
	return &order
}

func (repo *InMemoryOrderRepository) Store(order Order) {
	repo.memoryStore[order.ID] = order
}
