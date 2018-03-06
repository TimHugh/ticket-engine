package common

import (
	"fmt"
)

type Location struct {
	ID           string
	SignatureKey string
}

type LocationRepository interface {
	Find(string) (Location, error)
	Store(Location) error
}

func NewLocationRepository(adapter Adapter) LocationRepository {
	return &locationRepository{adapter}
}

type locationRepository struct {
	adapter Adapter
}

func (r locationRepository) Find(id string) (Location, error) {
	var loc Location
	err := r.adapter.Find("locations", id, &loc)
	return loc, err
}

func (r locationRepository) Store(loc Location) error {
	return r.adapter.Create("locations", loc)
}
