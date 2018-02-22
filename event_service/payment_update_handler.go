package main

import ()

type PaymentUpdateHandler struct {
	orderCreator OrderCreator
}

func (p PaymentUpdateHandler) Handle(event Event) error {
	return p.orderCreator.Create(event.OrderID, event.LocationID)
}
