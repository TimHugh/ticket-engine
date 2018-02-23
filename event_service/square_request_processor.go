package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type RequestValidator interface {
	Validate(SquareRequest) error
}

type SquareRequestProcessor struct {
	validators  []RequestValidator
	eventRouter EventRouter
}

func NewSquareRequestProcessor(eventRouter EventRouter) RequestProcessor {
	return &SquareRequestProcessor{
		validators:  []RequestValidator{},
		eventRouter: eventRouter,
	}
}

func (proc *SquareRequestProcessor) AddValidator(validator RequestValidator) {
	proc.validators = append(proc.validators, validator)
}

func (proc *SquareRequestProcessor) Process(r *http.Request) error {
	squareRequest := proc.serializeRequest(r)
	if err := proc.validate(squareRequest); err != nil {
		return err
	}
	return proc.eventRouter.Dispatch(squareRequest.Event)
}

func (proc *SquareRequestProcessor) validate(req SquareRequest) error {
	for _, validator := range proc.validators {
		if err := validator.Validate(req); err != nil {
			return err
		}
	}
	return nil
}

func (proc *SquareRequestProcessor) parseEvent(body []byte) Event {
	var event Event
	json.Unmarshal(body, &event)
	return event
}

func (proc *SquareRequestProcessor) serializeRequest(req *http.Request) SquareRequest {
	buf := proc.cloneBody(req)
	body, _ := ioutil.ReadAll(buf)

	return SquareRequest{
		URL:       "https://" + req.Host + req.URL.Path,
		Body:      string(body),
		Signature: req.Header.Get("X-Square-Signature"),
		Event:     proc.parseEvent(body),
	}
}

func (proc *SquareRequestProcessor) cloneBody(req *http.Request) io.ReadCloser {
	buf, _ := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	return ioutil.NopCloser(bytes.NewBuffer(buf))
}
