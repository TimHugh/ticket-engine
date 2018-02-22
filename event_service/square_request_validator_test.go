package main

import (
	"testing"

	"fmt"

	"github.com/timhugh/ticket-engine/common"
)

type MockLocationRepository struct {
	Location *common.Location
}

func (m *MockLocationRepository) Find(locationID string) *common.Location {
	return m.Location
}

func (m *MockLocationRepository) Store(location common.Location) {
	m.Location = &location
}

const goodSignature = "vsTe0jrY7ypjTdir98ES097hqN0="

func mockRequest(signature string) Request {
	return Request{
		Body:       `{"event": "test"}`,
		Signature:  signature,
		LocationID: "location_id",
		URL:        "http://www.example.com/event",
	}
}

func TestValidatesSignatureHeader(t *testing.T) {
	validator := SquareRequestValidator{
		locationRepository: &MockLocationRepository{
			Location: &common.Location{
				ID:           "location_id",
				SignatureKey: "test_key",
			},
		},
	}

	// good signature
	request := mockRequest(goodSignature)
	err := validator.Validate(request)
	if err != nil {
		t.Errorf(fmt.Sprintf("Got unexpected error: %s", err.Error()))
	}

	// bad signature
	request = mockRequest("/8sKiBs8HSiamqGkG2vSHOopv+w=")
	err = validator.Validate(request)
	if err == nil {
		t.Error("Expected error but received none.")
	}
}

func TestValidatesLocation(t *testing.T) {
	validator := SquareRequestValidator{
		locationRepository: &MockLocationRepository{
			Location: &common.Location{
				SignatureKey: "test_key",
			},
		},
	}

	request := mockRequest(goodSignature)
	err := validator.Validate(request)
	if err != nil {
		t.Errorf(fmt.Sprintf("Got unexpected error: %s", err.Error()))
	}

	validator = SquareRequestValidator{
		locationRepository: &MockLocationRepository{},
	}

	request = mockRequest("/8sKiBs8HSiamqGkG2vSHOopv+w=")
	err = validator.Validate(request)
	if err == nil {
		t.Error("Expected error but received none.")
	}
}
