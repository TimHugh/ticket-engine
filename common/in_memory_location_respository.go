package common

import ()

type InMemoryLocationRepository struct {
	memoryStore map[string]Location
}

func NewInMemoryLocationRepository() LocationRepository {
	return &InMemoryLocationRepository{
		make(map[string]Location),
	}
}

func (repo *InMemoryLocationRepository) Find(id string) *Location {
	location := repo.memoryStore[id]
	if location.ID == "" {
		return nil
	}
	return &location
}

func (repo *InMemoryLocationRepository) Store(location Location) {
	repo.memoryStore[location.ID] = location
}
