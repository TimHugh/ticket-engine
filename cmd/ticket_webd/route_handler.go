package main

import (
	"io"
	"net/http"
)

type RouteHandler struct {
	App App

	Processor RequestProcessor
}

type RequestProcessor interface {
	Process(*http.Request) error
}

func NewRouteHandler(app App) RouteHandler {
	processor := NewSquareRequestProcessor(app)
	return RouteHandler{
		App:       app,
		Processor: processor,
	}
}

func (h RouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ok := func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"status": "OK"}`)
	}

	unprocessable := func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"error": "unable to process"}`)
	}

	if r.Method != "POST" {
		ok(w)
	} else if err := h.Processor.Process(r); err != nil {
		h.App.Logger.Printf(`event=error message="%s"`, err)
		h.App.ErrorReporter.Error(err)
		unprocessable(w)
	} else {
		ok(w)
	}
}
