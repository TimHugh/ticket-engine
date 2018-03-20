package main

import (
	"testing"

	"fmt"

	"github.com/timhugh/ticket_service/mock"
	"github.com/timhugh/ticket_service/root"
)

const goodSignature = "vsTe0jrY7ypjTdir98ES097hqN0="

func mockRequest(signature string, location_id string) *SquareRequest {
	return &SquareRequest{
		Body:      `{"event": "test"}`,
		Signature: signature,
		URL:       "http://www.example.com/event",
		Event: Event{
			LocationID: location_id,
		},
	}
}

func TestGoodSignature(t *testing.T) {
	r := &mock.LocationRepository{}
	v := SquareRequestValidator{r}

	r.FindFn = func(id string) (*root.Location, error) {
		return &root.Location{
			ID:           "location_id",
			SignatureKey: "test_key",
		}, nil
	}

	request := mockRequest(goodSignature, "location_id")
	if err := v.Validate(request); err != nil {
		t.Errorf("got unexpected error: %s", err.Error())
	}
}

func TestBadSignature(t *testing.T) {
	r := &mock.LocationRepository{}
	v := SquareRequestValidator{r}

	r.FindFn = func(id string) (*root.Location, error) {
		return &root.Location{
			ID:           "location_id",
			SignatureKey: "test_key",
		}, nil
	}

	request := mockRequest("bad_signature", "location_id")
	if err := v.Validate(request); err == nil {
		t.Error("expected error but received none")
	}
}

func TestBadLocation(t *testing.T) {
	r := &mock.LocationRepository{}
	v := SquareRequestValidator{r}

	r.FindFn = func(id string) (*root.Location, error) {
		return nil, fmt.Errorf("not found")
	}

	request := mockRequest("signature", "location_id")
	if err := v.Validate(request); err == nil {
		t.Error("expected error but received none")
	}
}
