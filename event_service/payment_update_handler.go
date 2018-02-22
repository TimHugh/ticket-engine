package main

import (
	"github.com/timhugh/ticket-engine/common"
)

type PaymentUpdateHandler struct {
	orderCreator OrderCreator
}

func NewPaymentUpdateHandler(orderRepository common.OrderRepository) PaymentUpdateHandler {
	return PaymentUpdateHandler{
		OrderCreator{orderRepository},
	}
}

func (p PaymentUpdateHandler) Handle(event Event) error {
	return p.orderCreator.Create(event.OrderID, event.LocationID)
}
