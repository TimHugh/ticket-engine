package common

import (
	"testing"
)

func TestLocationRepo(t *testing.T) {
	adapter := NewMockAdapter()
	repo := LocationRepository{adapter}

	newLoc := Location{
		ID:           "id",
		SignatureKey: "signature",
	}
	storeErr := repo.Store(newLoc)
	if storeErr != nil {
		t.Errorf("Expected to store without error but received %s", storeErr)
	}

	findLoc, findErr := repo.Find("id")
	if findErr != nil {
		t.Errorf("Expected to find without error but received %s", findErr)
	}
	if findLoc.ID != newLoc.ID || findLoc.SignatureKey != newLoc.SignatureKey {
		t.Errorf("Expected to find %s but found %s", newLoc, findLoc)
	}
}
