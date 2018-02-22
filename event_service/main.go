package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"

	"github.com/timhugh/ticket-engine/common"
)

func main() {
	orderRepository := common.NewInMemoryOrderRepository()

	locationRepository := common.NewInMemoryLocationRepository()

	requestValidator := SquareRequestValidator{
		locationRepository: locationRepository,
	}
	eventRouter := NewEventRouter()
	eventRouter.Register("PAYMENT_UPDATED", &PaymentUpdateHandler{
		OrderCreator{orderRepository},
	})

	router := mux.NewRouter()
	router.HandleFunc("/event", eventHandler(requestValidator, eventRouter))
	n := negroni.Classic()
	n.UseHandler(router)
	http.ListenAndServe(":8080", n)
}

type JSON map[string]string

func eventHandler(requestValidator RequestValidator, eventRouter *EventRouter) http.HandlerFunc {
	r := render.New()

	return func(w http.ResponseWriter, req *http.Request) {
		event := parseEvent(req)
		request := serializeRequest(req, event)
		if err := requestValidator.Validate(request); err != nil {
			log.Println("Failed signature validation.", err)
			r.JSON(w, http.StatusNotFound, JSON{
				"error":         "Failed signature validation",
				"bad_signature": request.Signature,
				"entity_id":     event.OrderID,
			})
		}

		if err := eventRouter.Dispatch(event); err != nil {
			log.Println(err)
		}
	}
}

func parseEvent(req *http.Request) Event {
	var event Event
	body := readBody(req)
	json.Unmarshal(body, &event)
	return event
}

func cloneBody(req *http.Request) io.ReadCloser {
	buf, _ := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	return ioutil.NopCloser(bytes.NewBuffer(buf))
}

func readBody(req *http.Request) []byte {
	body := cloneBody(req)
	buf, _ := ioutil.ReadAll(body)
	return buf
}

func serializeRequest(req *http.Request, event Event) Request {
	return Request{
		URL:        "https://" + req.Host + req.URL.Path,
		Body:       string(readBody(req)),
		Signature:  req.Header.Get("X-Square-Signature"),
		LocationID: event.LocationID,
	}
}
