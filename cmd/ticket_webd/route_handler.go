package main

import (
	"fmt"
	"io"
	"net/http"
)

type RouteHandler struct {
	App AppContext

	Processor RequestProcessor
}

type RequestProcessor interface {
	Process(*http.Request) error
}

func NewRouteHandler(app AppContext) RouteHandler {
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

	notfound := func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"error": "not found"}`)
	}

	if r.Method != "POST" {
		err := fmt.Errorf("route does not exist: %s %s", r.Method, r.URL.Path)
		h.App.Logger.Printf(`event=error message="%s"`, err)
		h.App.Report.Error(err)
		notfound(w)
	} else if err := h.Processor.Process(r); err != nil {
		h.App.Logger.Printf(`event=error message="%s"`, err)
		h.App.Report.Error(err)
		unprocessable(w)
	} else {
		ok(w)
	}
}
