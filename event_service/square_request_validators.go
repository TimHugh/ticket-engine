package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/timhugh/ticket-engine/common"
)

type SquareRequest struct {
	Body      string
	Signature string
	URL       string
	Event     Event
}

type SquareRequestValidator struct {
	locationRepository common.LocationRepository
}

func (s SquareRequestValidator) Validate(req SquareRequest) error {
	location := s.locationRepository.Find(req.Event.LocationID)
	if location == nil {
		return errors.New(fmt.Sprintf("Received request for unknown location: %s.", req.Event.LocationID))
	}
	if !s.verifySignature(req.URL, location.SignatureKey, req.Body, req.Signature) {
		return errors.New("Request failed signature validation")
	}
	return nil
}

func (s SquareRequestValidator) verifySignature(webhookURL string, webhookSignatureKey string, body string, signature string) bool {
	mac := hmac.New(sha1.New, []byte(webhookSignatureKey))
	mac.Write([]byte(webhookURL + body))
	expectedMAC := mac.Sum(nil)
	expectedSignature := base64.StdEncoding.EncodeToString([]byte(expectedMAC))
	return expectedSignature == signature
}
