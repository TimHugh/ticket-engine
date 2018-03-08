package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
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
	locationRepository locationFinder
}

type locationFinder interface {
	Find(string) (*common.Location, error)
}

func (s SquareRequestValidator) Validate(req *SquareRequest) error {
	location, err := s.locationRepository.Find(req.Event.LocationID)
	if err != nil {
		return fmt.Errorf("Error finding location '%s': %s", req.Event.LocationID, err)
	}

	if !s.verifySignature(req.URL, location.SignatureKey, req.Body, req.Signature) {
		return fmt.Errorf("Request failed signature validation")
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
