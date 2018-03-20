package main

import (
	"net/http"

	"github.com/timhugh/ticket_service/root"
)

type EventHandler struct {
	Processor RequestProcessor
}

func (h EventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	} else if err := h.processor.Process(r); err != nil {
		h.logger.Printf(`event=error message="%s"`, err)
		h.reporter.Error(err)
		unprocessable(w)
	} else {
		ok(w)
	}
}
