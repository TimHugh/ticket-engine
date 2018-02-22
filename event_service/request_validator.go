package main

import ()

type RequestValidator interface {
	Validate(Request) error
}

type Request struct {
	Body       string
	Signature  string
	LocationID string
	URL        string
}
