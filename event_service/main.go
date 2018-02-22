package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"

	"github.com/timhugh/ticket-engine/common"
)

func main() {
	orderRepository := common.NewInMemoryOrderRepository()
	locationRepository := common.NewInMemoryLocationRepository()

	eventRouter := NewEventRouter()
	eventRouter.Register("PAYMENT_UPDATED", NewPaymentUpdateHandler(orderRepository))

	requestProcessor := NewSquareRequestProcessor(eventRouter)
	requestProcessor.AddValidator(SquareRequestValidator{locationRepository})

	router := mux.NewRouter()
	router.HandleFunc("/event", eventHandler(requestProcessor))

	n := negroni.Classic()
	n.UseHandler(router)

	http.ListenAndServe(":8080", n)
}

type RequestProcessor interface {
	AddValidator(RequestValidator)
	Process(*http.Request) error
}

type JSON map[string]string

func eventHandler(processor RequestProcessor) http.HandlerFunc {
	r := render.New()

	return func(w http.ResponseWriter, req *http.Request) {
		err := processor.Process(req)
		if err != nil {
			r.JSON(w, http.StatusUnprocessableEntity, JSON{
				"error": "unable to process",
			})
		} else {
			r.JSON(w, http.StatusOK, JSON{
				"status": "OK",
			})
		}
	}
}
