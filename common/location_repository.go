package common

import ()

type Location struct {
	ID           string
	SignatureKey string
}

type LocationRepository struct {
	adapter Adapter
}

func (r LocationRepository) Find(id string) (Location, error) {
	var loc Location
	err := r.adapter.Find("locations", id, &loc)
	return loc, err
}

func (r LocationRepository) Store(loc Location) error {
	return r.adapter.Create("locations", loc)
}
