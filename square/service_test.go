package square

import (
	"testing"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	root "github.com/timhugh/ticket_service"
	"github.com/timhugh/ticket_service/mock"
)

func TestPaymentFetch(t *testing.T) {
	c := &mock.HTTPClient{}
	s := Service{
		client: c,
		logger: &mock.Logger{},
	}

	location := root.Location{
		ID:       "location",
		APIToken: "token",
	}
	paymentID := "payment"

	body := fmt.Sprintf(`{"id": "%s"}`, paymentID)
	c.DoFunc = func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("Authentication") != "Bearer token" {
			return nil, fmt.Errorf("bad auth header")
		}

		buf := bytes.NewBuffer([]byte(body))
		return &http.Response{
			Body: ioutil.NopCloser(buf),
		}, nil
	}

	payment, err := s.GetPayment(location, paymentID)
	if err != nil {
		t.Errorf("expected nil error but got '%s'", err)
	}

	if payment.ID != paymentID {
		t.Errorf("expected payment ID '%s' but got '%s'", paymentID, payment.ID)
	}
}
