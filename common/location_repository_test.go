package common

import (
	"testing"
)

type MockLocationRepository struct {
	Location *Location
}

func (m *MockLocationRepository) Find(locationID string) *Location {
	return m.Location
}

func (m *MockLocationRepository) Store(location Location) {
	m.Location = &location
}
