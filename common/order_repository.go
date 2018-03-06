package common

import ()

type Order struct {
	ID         string
	LocationID string
}

type OrderRepository struct {
	Adapter Adapter
}

func (r OrderRepository) Find(id string) (*Order, error) {
	var order Order
	err := r.Adapter.Find("orders", id, &order)
	return &order, err
}

func (r OrderRepository) Store(order Order) error {
	return r.Adapter.Create("orders", order)
}
