package repos

import ()

type Order struct {
	ID         string
	LocationID string
}

type OrderRepository interface {
	Store(Order)
	Find(string) *Order
}
