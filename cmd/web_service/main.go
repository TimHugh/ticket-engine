package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/timhugh/ticket_service/mongo"
	"github.com/timhugh/ticket_service/rollbar"
)

type ErrorReporter interface {
	Error(error) error
}

type Logger interface{}

type Config map[string]string

type App struct {
	Config
	ErrorReporter
	Logger
	OrderRepository
	LocationRepository
}

func main() {
	config := Config{
		"environment":    os.Getenv("ENVIRONMENT"),
		"port":           os.Getenv("PORT"),
		"newrelic_token": os.Getenv("NEW_RELIC_TOKEN"),
		"newrelic_app":   os.Getenv("NEW_RELIC_APP_NAME"),
		"rollbar_token":  os.Getenv("ROLLBAR_TOKEN"),
		"mongodb_uri":    os.Getenv("MONGODB_URI"),
	}

	logger := log.New(os.Stdout, "web", log.LstdFlags)
	rollbarReporter := rollbar.New(config["rollbar_token"], config["environment"])

	mongoSession, mongoErr := mongo.NewMongoSession(config["mongodb_uri"])
	fail(mongoErr)

	app := App{
		Config:             config,
		ErrorReporter:      rollbarReporter,
		Logger:             logger,
		OrderRepository:    OrderRepository{mongoSession},
		LocationRepository: LocationRepository{mongoSession},
	}

	mux := http.NewServeMux()
	mux.Handle("/event", NewEventHandler(app))

	log.Printf("listening on %s\n", config["port"])
	serveErr := http.ListenAndServe(fmt.Sprintf(":%s", config["port"]), mux)
	fail(serveErr)
}

func fail(err error) {
	if err != nil {
		panic(err)
	}
}
