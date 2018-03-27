package square

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	root "github.com/timhugh/ticket_service"
)

type Service struct {
	client client
	logger logger
}

// nested struct: https://play.golang.org/p/a81HUuLgwl-
// extracted: https://play.golang.org/p/wCh3lQ2DuKE

type Payment struct {
	ID            string `json:"id"`
	LocationID    string `json:"merchant_id"`
	PaymentURL    string `json:"payment_url"`
	TransactionID string
	Items         []Item `json:"itemizations"`
}

type Item struct {
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
	Detail   struct {
		ItemID string `json:"item_id"`
	} `json:"item_detail"`
}

type client interface {
	Do(*http.Request) (*http.Response, error)
}

type logger interface {
	Printf(string, ...interface{})
}

func New(logger logger) Service {
	return Service{
		client: &http.Client{},
		logger: logger,
	}
}

const connectURL string = "https://connect.squareup.com/api"

func (s Service) GetPayment(location root.Location, id string) (*Payment, error) {
	paymentURL := fmt.Sprintf(connectURL+"/v1/%s/payments/%s", location.ID, id)

	body, err := s.getBody(paymentURL, location.APIToken)
	if err != nil {
		return nil, err
	}

	var payment Payment
	if parseErr := json.Unmarshal(body, &payment); parseErr != nil {
		return nil, parseErr
	}
	return &payment, nil
}

func (s Service) getBody(url, token string) ([]byte, error) {
	resp, getErr := s.get(url, token)
	if getErr != nil {
		return nil, getErr
	}
	defer resp.Body.Close()

	buf, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	return buf, nil
}

func (s Service) get(url, token string) (*http.Response, error) {
	req, newErr := http.NewRequest("GET", url, nil)
	if newErr != nil {
		return nil, newErr
	}

	req.Header.Add("Authentication", fmt.Sprintf("Bearer %s", token))

	// TODO log arounnd request
	resp, doErr := s.client.Do(req)
	if doErr != nil {
		return nil, doErr
	}

	return resp, nil
}
