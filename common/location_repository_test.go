package common

import (
	"testing"

	"fmt"
)

type mockLocationAdapter struct {
	Location *Location
}

func (a mockLocationAdapter) Find(collection string, id string, location interface{}) error {
	if id == "good_id" && collection == "locations" {
		pLocation := location.(*Location)
		pLocation.ID = "good_id"
		pLocation.SignatureKey = "test_signature"
		return nil
	} else {
		return fmt.Errorf("not found")
	}
}

func (a *mockLocationAdapter) Create(collection string, location interface{}) error {
	value := location.(Location)
	a.Location = &value
	return nil
}

func (a mockLocationAdapter) Close() {}

func TestLocationRepoFind(t *testing.T) {
	adapter := &mockLocationAdapter{}
	repo := LocationRepository{adapter}

	newLoc := Location{
		ID:           "id",
		SignatureKey: "test_signature",
	}
	err := repo.Store(newLoc)
	if err != nil {
		t.Errorf("Expected to store without error but received %s", err)
	}
}

func TestLocationRepoStore(t *testing.T) {
	adapter := &mockLocationAdapter{}
	repo := LocationRepository{adapter}

	goodLoc, goodErr := repo.Find("good_id")
	if goodErr != nil {
		t.Errorf("Expected to find without error but received %s", goodErr)
	}
	if goodLoc.ID != "good_id" {
		t.Errorf("Expected to find location with id 'good_id' but got %s", goodLoc.ID)
	}
	if goodLoc.SignatureKey != "test_signature" {
		t.Errorf("Expected to find location with signature key 'test_signature' but got %s", goodLoc.SignatureKey)
	}

	badLoc, badErr := repo.Find("bad_id")
	if badErr == nil {
		t.Errorf("Expected 'not found' error but received none.")
	}
	if badLoc != nil {
		t.Errorf("Expected missing location to be nil but got %s", badLoc)
	}
}
