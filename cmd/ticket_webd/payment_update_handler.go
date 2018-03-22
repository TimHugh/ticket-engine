package main

import (
	"fmt"

	root "github.com/timhugh/ticket_service"
)

type PaymentUpdateHandler struct {
	OrderRepository OrderRepository
}

func (p PaymentUpdateHandler) Handle(event Event) error {
	_, err := p.OrderRepository.Find(event.LocationID, event.OrderID)
	if err.Error() != "not found" {
		return err
	} else if err != nil {
		return fmt.Errorf("Couldn't create duplicate order %s.", event.OrderID)
	}

	return p.OrderRepository.Create(root.Order{
		ID:         event.OrderID,
		LocationID: event.LocationID,
	})
}
