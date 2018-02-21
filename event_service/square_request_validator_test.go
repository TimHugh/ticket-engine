package event_service

import (
	"testing"

	"fmt"
)

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
		LocationRepository: &MockLocationRepository{
			Location: &Location{
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
		LocationRepository: &MockLocationRepository{
			Location: &Location{
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
		LocationRepository: &MockLocationRepository{},
	}

	request = mockRequest("/8sKiBs8HSiamqGkG2vSHOopv+w=")
	err = validator.Validate(request)
	if err == nil {
		t.Error("Expected error but received none.")
	}
}
