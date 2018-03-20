package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/timhugh/ticket_service/mongo"
	"github.com/timhugh/ticket_service/rollbar"
	"github.com/timhugh/ticket_service/root"
)

type ErrorReporter interface {
	Error(error) error
}

type Logger interface {
	Printf(string, ...interface{})
}

type Config map[string]string

type App struct {
	Config
	ErrorReporter
	Logger
	OrderRepository
	LocationRepository
}

type OrderRepository interface {
	Find(string, string) (*root.Order, error)
	Create(root.Order) error
}

type LocationRepository interface {
	Find(string) (*root.Location, error)
}

func main() {
	config := Config{
		"environment":    *flag.String("env", "development", "application environment"),
		"http":           *flag.String("http", ":8080", "HTTP service address"),
		"newrelic_token": *flag.String("new_relic_token", "", "New Relic API token"),
		"newrelic_app":   *flag.String("new_relic_app", "", "New Relic application name"),
		"rollbar_token":  *flag.String("rollbar_token", "", "Rollbar API token"),
		"mongodb_uri":    *flag.String("mongodb", "", "MongoDB host URI"),
	}

	log.Printf("running with %s environment config\n", config["environment"])

	mongoSession, mongoErr := mongo.NewMongoSession(config["mongodb_uri"])
	fail(mongoErr)

	app := App{
		Config:             config,
		ErrorReporter:      rollbar.New(config["rollbar_token"], config["environment"]),
		Logger:             log.New(os.Stdout, "web", log.LstdFlags),
		OrderRepository:    mongo.OrderRepository{mongoSession},
		LocationRepository: mongo.LocationRepository{mongoSession},
	}

	mux := http.NewServeMux()
	mux.Handle("/event", NewRouteHandler(app))

	log.Printf("listening on %s\n", config["port"])
	serveErr := http.ListenAndServe(fmt.Sprintf(":%s", config["port"]), mux)
	fail(serveErr)
}

func fail(err error) {
	if err != nil {
		panic(err)
	}
}
