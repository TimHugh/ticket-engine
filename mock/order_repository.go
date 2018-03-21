package mock

import (
	root "github.com/timhugh/ticket_service"
)

var _ root.OrderRepository = &OrderRepository{}

type OrderRepository struct {
	FindFn      func(string, string) (*root.Order, error)
	FindInvoked bool

	CreateFn      func(root.Order) error
	CreateInvoked bool
}

func (r *OrderRepository) Find(locationID, id string) (*root.Order, error) {
	r.FindInvoked = true
	return r.FindFn(locationID, id)
}

func (r *OrderRepository) Create(order root.Order) error {
	r.CreateInvoked = true
	return r.CreateFn(order)
}
