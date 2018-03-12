package main

import (
	"testing"

	"fmt"

	"github.com/timhugh/ticket_service/common"
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

type mockLocationFinder struct {
	Location *common.Location
}

func (m mockLocationFinder) Find(id string) (*common.Location, error) {
	if m.Location.ID == id {
		return m.Location, nil
	} else {
		return nil, fmt.Errorf("Unable to find location")
	}
}

func mockRepo() locationFinder {
	locationRepository := mockLocationFinder{
		&common.Location{
			ID:           "location_id",
			SignatureKey: "test_key",
		},
	}
	return locationRepository
}

func TestSignatureValidation(t *testing.T) {
	validator := SquareRequestValidator{mockRepo()}

	// good location, good signature
	request := mockRequest(goodSignature, "location_id")
	err := validator.Validate(request)
	if err != nil {
		t.Errorf("Got unexpected error: %s", err.Error())
	}

	// good location, bad signature
	request = mockRequest("/8sKiBs8HSiamqGkG2vSHOopv+w=", "location_id")
	err = validator.Validate(request)
	if err == nil {
		t.Error("Expected error but received none.")
	}

	// bad location
	request = mockRequest(goodSignature, "not_location_id")
	err = validator.Validate(request)
	if err == nil {
		t.Error("Expected error but received none.")
	}
}
