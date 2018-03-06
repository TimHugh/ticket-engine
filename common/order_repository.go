package common

import ()

type Order struct {
	ID         string
	LocationID string
}

type OrderRepository struct {
	adapter Adapter
}

func (r OrderRepository) Find(id string) (Order, error) {
	var order Order
	err := r.adapter.Find("orders", id, &order)
	return order, err
}

func (r OrderRepository) Store(order Order) error {
	return r.adapter.Create("orders", order)
}
