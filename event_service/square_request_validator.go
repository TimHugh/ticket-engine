package event_service

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/timhugh/ticket-engine/common"
)

type SquareRequestValidator struct {
	LocationRepository common.LocationRepository
}

func (s SquareRequestValidator) Validate(request Request) error {
	location := s.LocationRepository.Find(request.LocationID)
	if location == nil {
		return errors.New(fmt.Sprintf("Received request for unknown location: %s.", request.LocationID))
	}

	if !verifySignature(request.URL, location.SignatureKey, request.Body, request.Signature) {
		return errors.New("Request failed signature validation.")
	}

	return nil
}

func verifySignature(webhookURL string, webhookSignatureKey string, body string, signature string) bool {
	mac := hmac.New(sha1.New, []byte(webhookSignatureKey))
	mac.Write([]byte(webhookURL + body))
	expectedMAC := mac.Sum(nil)
	expectedSignature := base64.StdEncoding.EncodeToString([]byte(expectedMAC))
	return expectedSignature == signature
}
