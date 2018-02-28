package common

import (
	"errors"
)

type InMemoryLocationRepository struct {
	memoryStore map[string]Location
}

func NewInMemoryLocationRepository() LocationRepository {
	return &InMemoryLocationRepository{
		make(map[string]Location),
	}
}

func (repo *InMemoryLocationRepository) Find(id string) (*Location, error) {
	location := repo.memoryStore[id]
	if location.ID == "" {
		return nil, errors.New("Location doesn't exist")
	}
	return &location, nil
}

func (repo *InMemoryLocationRepository) Store(location Location) error {
	repo.memoryStore[location.ID] = location
	return nil
}
