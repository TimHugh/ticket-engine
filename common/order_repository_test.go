package common

import (
	"testing"
)

func TestOrderRepo(t *testing.T) {
	adapter := NewMockAdapter()
	repo := OrderRepository{adapter}

	newOrder := Order{
		ID:         "id",
		LocationID: "location_id",
	}
	storeErr := repo.Store(newOrder)
	if storeErr != nil {
		t.Errorf("Expected to store without error but received %s", storeErr)
	}

	findOrder, findErr := repo.Find("id")
	if findErr != nil {
		t.Errorf("Expected to find without error but received %s", findErr)
	}
	if findOrder.ID != newOrder.ID || findOrder.LocationID != newOrder.LocationID {
		t.Errorf("Expected to find %s but found %s", newOrder, findOrder)
	}
}
