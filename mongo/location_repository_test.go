package mongo

import (
	"testing"

	"fmt"

	"github.com/timhugh/ticket_service/root"
)

func TestFindsLocations(t *testing.T) {
	s := &MockSession{}
	r := LocationRepository{s}

	s.OneFn = func(out interface{}) error {
		result := out.(*root.Location)
		result.ID = "location"
		result.SignatureKey = "key"
		return nil
	}

	location, err := r.Find("id")
	if err != nil {
		t.Errorf("expected nil error but got %s", err)
	}
	if location.ID != "location" || location.SignatureKey != "key" {
		t.Errorf(
			"expected location with id 'location' and signature key 'key' but got id '%s' and signature key '%s'",
			location.ID, location.SignatureKey,
		)
	}
}

func TestErrorsForLocationNotFound(t *testing.T) {
	s := &MockSession{}
	r := LocationRepository{s}

	s.OneFn = func(out interface{}) error {
		return fmt.Errorf("not found")
	}

	_, err := r.Find("id")
	if err == nil {
		t.Errorf("expected error but received none")
	}
}

func TestCreatesNewLocations(t *testing.T) {
	s := &MockSession{}
	r := LocationRepository{s}

	s.InsertFn = func(docs ...interface{}) error {
		return nil
	}

	location := root.Location{
		ID:           "id",
		SignatureKey: "key",
	}

	if err := r.Create(location); err != nil {
		t.Errorf("expected nil error but got %s", err)
	}
	if s.InsertInvoked != true {
		t.Error("expected insert method to be invoked")
	}
}
