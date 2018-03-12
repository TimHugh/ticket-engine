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

	"github.com/timhugh/ticket_service/common"
)

var config = map[string]string{
	"environment":    os.Getenv("ENVIRONMENT"),
	"port":           os.Getenv("PORT"),
	"newrelic_token": os.Getenv("NEW_RELIC_TOKEN"),
	"newrelic_app":   os.Getenv("NEW_RELIC_APP_NAME"),
	"rollbar_token":  os.Getenv("ROLLBAR_TOKEN"),
	"mongodb_uri":    os.Getenv("MONGODB_URI"),
}

type RequestProcessor interface {
	AddValidator(RequestValidator)
	Process(*http.Request) error
}

func main() {
	adapter, err := common.NewMongoAdapter(config["mongodb_uri"])
	if err != nil {
		rollbar.Error(rollbar.ERR, err)
		log.Fatal(err)
	}
	orderRepository := common.OrderRepository{Adapter: adapter}
	locationRepository := common.LocationRepository{Adapter: adapter}

	eventRouter := NewEventRouter()
	eventRouter.Register("PAYMENT_UPDATED", NewPaymentUpdateHandler(orderRepository))
	eventRouter.Register("TEST_NOTIFICATION", NoopHandler{})

	requestProcessor := NewSquareRequestProcessor(eventRouter)
	requestProcessor.AddValidator(SquareRequestValidator{locationRepository})

	initRollbar(config["rollbar_token"], config["environment"])

	nrApp := initNewRelic(config["newrelic_token"], config["newrelic_app"])

	router := mux.NewRouter()
	router.HandleFunc(newrelic.WrapHandleFunc(nrApp, "/event", requestHandler(requestProcessor)))

	n := negroni.Classic()
	n.UseHandler(router)

	log.Printf("Starting server on port %s", config["port"])
	err = http.ListenAndServe(fmt.Sprintf(":%s", config["port"]), n)
	if err != nil {
		rollbar.Error(rollbar.ERR, err)
		log.Fatal(err)
	}
}

func initNewRelic(token string, appName string) newrelic.Application {
	config := newrelic.NewConfig(appName, token)
	app, err := newrelic.NewApplication(config)
	if err != nil {
		rollbar.Error(rollbar.ERR, err)
		log.Fatal("Unable to initialize NewRelic reporting.")
	}
	log.Printf("Reporting to New Relic as '%s'", appName)
	return app
}

func initRollbar(token string, environment string) {
	rollbar.Environment = environment
	rollbar.Token = token
}

type JSON map[string]string

func requestHandler(processor RequestProcessor) http.HandlerFunc {
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
			ok(w)
		} else if err := processor.Process(req); err != nil {
			log.Printf(`event=error message="%s"`, err)
			rollbar.RequestError(rollbar.ERR, req, err)
			unprocessable(w)
		} else {
			ok(w)
		}
	}
}
