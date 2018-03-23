package main

import (
	"testing"

	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/timhugh/ticket_service/mock"
)

type MockProcessor struct {
	ProcessFn      func(*http.Request) error
	ProcessInvoked bool
}

func (p *MockProcessor) Process(r *http.Request) error {
	p.ProcessInvoked = true
	return p.ProcessFn(r)
}

func TestSuccessfulProcessor(t *testing.T) {
	p := &MockProcessor{}
	h := RouteHandler{Processor: p}

	p.ProcessFn = func(req *http.Request) error {
		return nil
	}
	req := httptest.NewRequest("POST", "/event", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("expected 200 response but got %d", status)
	}

	expected := `{"status": "OK"}`
	if body := w.Body.String(); body != expected {
		t.Errorf("expected body:\n\t%s\ngot:\n\t%s", expected, body)
	}
}

func TestErrorProcessor(t *testing.T) {
	p := &MockProcessor{}

	l := &mock.Logger{}
	r := &mock.ErrorReporter{}

	h := RouteHandler{
		App: AppContext{
			Logger: l,
			Report: r,
		},
		Processor: p,
	}

	errorText := "some error"
	p.ProcessFn = func(req *http.Request) error {
		return errors.New(errorText)
	}
	req := httptest.NewRequest("POST", "/event", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 status but got %d", status)
	}

	expectedBody := `{"error": "unable to process"}`
	if body := w.Body.String(); body != expectedBody {
		t.Errorf("expected body:\n\t%s\ngot:\n\t%s", expectedBody, body)
	}

	if err := r.Store.Error(); err != errorText {
		t.Errorf(`expected reporter to recieve error "%s" but got "%s"`, errorText, err)
	}
	if !l.Contains(errorText) {
		t.Errorf("expected log to contain '%s':\n%s", errorText, l.Out())
	}
}

func TestNonPost(t *testing.T) {
	r := &mock.ErrorReporter{}
	l := &mock.Logger{}
	h := RouteHandler{
		App: AppContext{
			Report: r,
			Logger: l,
		},
	}

	req := httptest.NewRequest("GET", "/event", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusNotFound {
		t.Errorf("expected 404 status but got %d", status)
	}

	expectedBody := `{"error": "not found"}`
	if body := w.Body.String(); body != expectedBody {
		t.Errorf("expected body:\n\t%s\ngot:\n\t%s", expectedBody, body)
	}

	errorText := "route does not exist: GET /event"
	if err := r.Store.Error(); err != errorText {
		t.Errorf(`expected reporter to recieve error "%s" but got "%s"`, errorText, err)
	}
	if !l.Contains(errorText) {
		t.Errorf("expected log to contain '%s':\n%s", errorText, l.Out())
	}
}
