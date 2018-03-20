package main

import (
	"testing"

	"net/httptest"
	"strings"
)

var requestTests = []struct {
	method string
	body   string
}{
	method: "POST",
}

func Test(t *testing.T) {
	handler := EventHandler{}

	for _, tt := range requestTests {
		body := strings.NewReader(tt.body)
		req := httptest.NewRequest(tt.method, "/event", body)
		w := httptest.NewRecorder()
	}
}
