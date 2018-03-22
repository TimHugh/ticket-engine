package main

import (
	"fmt"

	root "github.com/timhugh/ticket_service"
)

type PaymentUpdateHandler struct {
	OrderRepository OrderRepository
}

func (p PaymentUpdateHandler) Handle(event Event) error {
	existingOrder, _ := p.OrderRepository.Find(event.LocationID, event.OrderID)
	if existingOrder != nil {
		return fmt.Errorf("Couldn't create duplicate order %s.", event.OrderID)
	}

	return p.OrderRepository.Create(root.Order{
		ID:         event.OrderID,
		LocationID: event.LocationID,
	})
}
