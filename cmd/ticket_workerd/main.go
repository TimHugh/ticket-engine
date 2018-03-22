package main

import (
	"log"
	"os"

	"github.com/streadway/amqp"

	"github.com/timhugh/ticket_service/config"
	//"github.com/timhugh/ticket_service/mongo"
	"github.com/timhugh/ticket_service/rollbar"
)

type ErrorReporter interface {
	Error(error) error
}

type Logger interface {
	Printf(string, ...interface{})
}

type AppContext struct {
	Config *config.Config
	Report ErrorReporter
	Logger Logger
}

func main() {
	var app AppContext

	app.Logger = log.New(os.Stdout, "", 0)

	app.Config = config.New()
	app.Config.Define("environment", "development", "application environment")
	app.Config.Define("newrelic_token", "", "New Relic API token")
	app.Config.Define("newrelic_app", "TicketService", "New Relic application name")
	app.Config.Define("rollbar_token", "", "Rollbar API token")
	app.Config.Define("mongodb_uri", "", "MongoDB host URI")
	app.Config.Define("rabbitmq_uri", "", "RabbitMQ host URI")
	app.Config.Load()

	app.Report = rollbar.New(app.Config.Get("rollbar_token"), app.Config.Get("environment"))

	//	mongoSession, mongoErr := mongo.NewMongoSession(app.Config.Get("mongodb_uri"))
	//	fail(app.Logger, mongoErr)

	queueConn, queueErr := amqp.Dial(app.Config.Get("rabbitmq_uri"))
	fail(app.Logger, queueErr)
	defer queueConn.Close()
}

func fail(log Logger, err error) {
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		panic(err)
	}
}
