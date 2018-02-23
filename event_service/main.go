package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"

	_ "github.com/joho/godotenv/autoload"
	"github.com/newrelic/go-agent"
	"github.com/stvp/rollbar"

	"github.com/timhugh/ticket-engine/common"
)

var config = map[string]string{
	"environment":    os.Getenv("ENVIRONMENT"),
	"port":           os.Getenv("PORT"),
	"newrelic_token": os.Getenv("NEW_RELIC_TOKEN"),
	"newrelic_app":   os.Getenv("NEW_RELIC_APP_NAME"),
	"rollbar_token":  os.Getenv("ROLLBAR_TOKEN"),
}

type RequestProcessor interface {
	AddValidator(RequestValidator)
	Process(*http.Request) error
}

func main() {
	orderRepository := common.NewInMemoryOrderRepository()
	locationRepository := common.NewInMemoryLocationRepository()

	eventRouter := NewEventRouter()
	eventRouter.Register("PAYMENT_UPDATED", NewPaymentUpdateHandler(orderRepository))

	requestProcessor := NewSquareRequestProcessor(eventRouter)
	requestProcessor.AddValidator(SquareRequestValidator{locationRepository})

	initRollbar(config["rollbar_token"], config["environment"])

	nrApp := initNewRelic(config["newrelic_token"], config["newrelic_app"])

	router := mux.NewRouter()
	router.HandleFunc(newrelic.WrapHandleFunc(nrApp, "/event", eventHandler(requestProcessor)))

	n := negroni.Classic()
	n.UseHandler(router)

	err := http.ListenAndServe(fmt.Sprintf(":%s", config["port"]), n)
	if err != nil {
		log.Fatal(err)
	}
}

func initNewRelic(token string, appName string) newrelic.Application {
	config := newrelic.NewConfig(appName, token)
	app, err := newrelic.NewApplication(config)
	if err != nil {
		log.Fatal("Unable to initialize NewRelic reporting.")
	}
	return app
}

func initRollbar(token string, environment string) {
	rollbar.Environment = environment
	rollbar.Token = token
}

type JSON map[string]string

func eventHandler(processor RequestProcessor) http.HandlerFunc {
	r := render.New()

	ok := func(w http.ResponseWriter) {
		r.JSON(w, http.StatusOK, JSON{
			"status": "OK",
		})
	}

	unprocessable := func(w http.ResponseWriter) {
		r.JSON(w, http.StatusUnprocessableEntity, JSON{
			"error": "unable to process",
		})
	}

	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			unprocessable(w)
		} else if err := processor.Process(req); err != nil {
			rollbar.Error(rollbar.ERR, err)
			unprocessable(w)
		} else {
			ok(w)
		}
	}
}
