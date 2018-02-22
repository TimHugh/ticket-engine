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

func (proc SquareRequestProcessor) Process(req *http.Request) error {
	event, err := proc.serializeEvent(req)
	if err != nil {
		return err
	}

	squareRequest := proc.serializeRequest(req, event)
	if err := proc.validate(squareRequest); err != nil {
		return err
	}

	return proc.eventRouter.Dispatch(event)
}

func (proc SquareRequestProcessor) validate(req SquareRequest) error {
	for _, validator := range proc.validators {
		if err := validator.Validate(req); err != nil {
			return err
		}
	}
	return nil
}

func (proc SquareRequestProcessor) serializeEvent(req *http.Request) (Event, error) {
	var event Event
	buf := proc.cloneBody(req)
	body, _ := ioutil.ReadAll(buf)
	err := json.Unmarshal(body, &event)
	return event, err
}

func (proc SquareRequestProcessor) serializeRequest(req *http.Request, event Event) SquareRequest {
	buf := proc.cloneBody(req)
	body, _ := ioutil.ReadAll(buf)
	return SquareRequest{
		URL:        "https://" + req.Host + req.URL.Path,
		Body:       string(body),
		Signature:  req.Header.Get("X-Square-Signature"),
		LocationID: event.LocationID,
	}
}

func (proc SquareRequestProcessor) cloneBody(req *http.Request) io.ReadCloser {
	buf, _ := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	return ioutil.NopCloser(bytes.NewBuffer(buf))
}
