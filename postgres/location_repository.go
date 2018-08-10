package postgres

import (
	root "github.com/timhugh/ticket_service"
)

type LocationRepository struct {
}

func (s LocationRepository) Create(location root.Location) error {
}

func (s LocationRepository) Find(id string) (*root.Location, error) {
}
