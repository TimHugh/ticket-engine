package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type RequestValidator interface {
	Validate(*SquareRequest) error
}

type SquareRequestProcessor struct {
	validators  []RequestValidator
	eventRouter EventRouter
}

func NewSquareRequestProcessor(app App) RequestProcessor {
	router := NewEventRouter()
	router.Register("PAYMENT_UPDATED", PaymentUpdateHandler{app.OrderRepository})
	router.Register("TEST_NOTIFICATION", NoopHandler{})

	return &SquareRequestProcessor{
		validators: []RequestValidator{
			SquareRequestValidator{app.LocationRepository},
		},
		eventRouter: NewEventRouter(),
	}
}

func (proc *SquareRequestProcessor) AddValidator(validator RequestValidator) {
	proc.validators = append(proc.validators, validator)
}

func (proc *SquareRequestProcessor) Process(r *http.Request) error {
	squareRequest, err := proc.serializeRequest(r)
	if err != nil {
		return fmt.Errorf("Error processing request data: %s", err)
	}
	if err = proc.validate(squareRequest); err != nil {
		return err
	}
	return proc.eventRouter.Dispatch(squareRequest.Event)
}

func (proc *SquareRequestProcessor) validate(req *SquareRequest) error {
	for _, validator := range proc.validators {
		if err := validator.Validate(req); err != nil {
			return err
		}
	}
	return nil
}

func (proc *SquareRequestProcessor) parseEvent(body []byte) (Event, error) {
	var event Event
	err := json.Unmarshal(body, &event)
	return event, err
}

func (proc *SquareRequestProcessor) serializeRequest(req *http.Request) (*SquareRequest, error) {
	buf := proc.cloneBody(req)
	body, _ := ioutil.ReadAll(buf)
	event, err := proc.parseEvent(body)
	if err != nil {
		return nil, err
	}

	log.Printf("event=square_event_received location_id=%s event_type=%s entity_id=%s", event.LocationID, event.Type, event.OrderID)

	return &SquareRequest{
		URL:       "https://" + req.Host + req.URL.Path,
		Body:      string(body),
		Signature: req.Header.Get("X-Square-Signature"),
		Event:     event,
	}, nil
}

func (proc *SquareRequestProcessor) cloneBody(req *http.Request) io.ReadCloser {
	buf, _ := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	return ioutil.NopCloser(bytes.NewBuffer(buf))
}
