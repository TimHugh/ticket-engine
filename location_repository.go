package ticket_service

type LocationRepository interface {
	Create(Location) error
	Find(string) (*Location, error)
}
