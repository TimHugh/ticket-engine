package event_service

import ()

type PaymentUpdateHandler struct {
	OrderCreator OrderCreator
}

func (p PaymentUpdateHandler) Handle(event Event) error {
	return p.OrderCreator.Create(event.OrderID, event.LocationID)
}
