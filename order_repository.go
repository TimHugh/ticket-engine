package ticket_service

type OrderRepository interface {
	Create(Order) error
	Find(string, string) (*Order, error)
}
