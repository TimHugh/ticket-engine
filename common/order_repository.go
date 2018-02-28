package common

import ()

type Order struct {
	ID         string
	LocationID string
}

type OrderRepository interface {
	Store(Order) error
	Find(string) (*Order, error)
}
