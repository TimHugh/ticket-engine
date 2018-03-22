package main

import (
	"log"
	"net/http"
	"os"

	root "github.com/timhugh/ticket_service"
	"github.com/timhugh/ticket_service/mongo"
	"github.com/timhugh/ticket_service/rollbar"
)

type ErrorReporter interface {
	Error(error) error
}

type Logger interface {
	Printf(string, ...interface{})
}

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
	config := Config{}
	config.Load()

	log.Printf("running with %s environment config\n", config.Get("environment"))

	mongoSession, mongoErr := mongo.NewMongoSession(config.Get("mongodb_uri"))
	fail(mongoErr)

	app := App{
		Config:             config,
		ErrorReporter:      rollbar.New(config.Get("rollbar_token"), config.Get("environment")),
		Logger:             log.New(os.Stdout, "web", log.LstdFlags),
		OrderRepository:    mongo.OrderRepository{Session: mongoSession},
		LocationRepository: mongo.LocationRepository{Session: mongoSession},
	}

	mux := http.NewServeMux()
	mux.Handle("/event", NewRouteHandler(app))

	log.Printf("listening on %s\n", config.Get("http"))
	serveErr := http.ListenAndServe(config.Get("http"), mux)
	fail(serveErr)
}

func fail(err error) {
	if err != nil {
		panic(err)
	}
}
