package common

import ()

type Location struct {
	ID           string
	SignatureKey string
}

type LocationRepository struct {
	Adapter Adapter
}

func (r LocationRepository) Find(id string) (*Location, error) {
	var loc Location
	if err := r.Adapter.Find("locations", id, &loc); err == nil {
		return &loc, nil
	} else {
		return nil, err
	}
}

func (r LocationRepository) Store(loc Location) error {
	return r.Adapter.Create("locations", loc)
}
