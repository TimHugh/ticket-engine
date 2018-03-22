package main

import (
	"log"
	"net/http"
	"os"

	root "github.com/timhugh/ticket_service"
	"github.com/timhugh/ticket_service/config"
	"github.com/timhugh/ticket_service/mongo"
	"github.com/timhugh/ticket_service/rollbar"
)

type ErrorReporter interface {
	Error(error) error
}

type Logger interface {
	Printf(string, ...interface{})
}

type AppContext struct {
	Config             Config
	Report             ErrorReporter
	Logger             Logger
	OrderRepository    OrderRepository
	LocationRepository LocationRepository
}

type Config interface {
	Get(string) string
	Define(string, string, string)
	Load()
}

type OrderRepository interface {
	Find(string, string) (*root.Order, error)
	Create(root.Order) error
}

type LocationRepository interface {
	Find(string) (*root.Location, error)
}

func main() {
	app := AppContext{}

	app.Logger = log.New(os.Stdout, "", 0)

	app.Config = config.New()
	app.Config.Define("environment", "development", "application environment")
	app.Config.Define("http", ":8080", "HTTP service address")
	app.Config.Define("newrelic_token", "", "New Relic API token")
	app.Config.Define("newrelic_app", "TicketService", "New Relic application name")
	app.Config.Define("rollbar_token", "", "Rollbar API token")
	app.Config.Define("mongodb_uri", "", "MongoDB host URI")
	app.Config.Load()

	app.Report = rollbar.New(app.Config.Get("rollbar_token"), app.Config.Get("environment"))

	mongoSession, mongoErr := mongo.NewMongoSession(app.Config.Get("mongodb_uri"))
	fail(mongoErr)
	app.OrderRepository = mongo.OrderRepository{Session: mongoSession}
	app.LocationRepository = mongo.LocationRepository{Session: mongoSession}

	mux := http.NewServeMux()
	mux.Handle("/event", NewRouteHandler(app))

	app.Logger.Printf("listening on %s\n", app.Config.Get("http"))
	serveErr := http.ListenAndServe(app.Config.Get("http"), mux)
	fail(serveErr)
}

func fail(err error) {
	if err != nil {
		panic(err)
	}
}
