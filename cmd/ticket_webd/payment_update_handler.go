package main

import ()

type PaymentUpdateHandler struct {
	orderCreator OrderCreator
}

func NewPaymentUpdateHandler(orderRepository OrderRepository) PaymentUpdateHandler {
	return PaymentUpdateHandler{
		OrderCreator{orderRepository},
	}
}

func (p PaymentUpdateHandler) Handle(event Event) error {
	return p.orderCreator.Create(event.OrderID, event.LocationID)
}
