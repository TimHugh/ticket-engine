package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	_ "github.com/joho/godotenv/autoload"

	"github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"
	"github.com/stvp/rollbar"
)

func initRollbar(env string, token string) {
	rollbar.Environment = env
	rollbar.Token = token
}

func main() {
	settings := loadSettings()

	initRollbar(
		settings["environment"],
		settings["rollbar.token"],
	)

	nrAgent := initNewRelic(
		settings["environment"],
		settings["newrelic.appName"],
		settings["newrelic.token"],
	)

	db := initDB()
	defer db.Close()

	controller := Controller{
		db: db,
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(nrgin.Middleware(nrAgent))
	r.Use(RequestLogger())
	r.POST("/event", controller.HandleEvent)
	r.Run()
}

type App struct {
	settings      SettingsGroup
	db            *gorm.DB
	router        gin.Engine
	errorReporter ErrorReporter
}

type ErrorReporter interface {
	Error(string, error, ...interface{})
}

type SettingsGroup map[string]string

func loadSettings() SettingsGroup {
	settings := SettingsGroup{
		"environment": os.Getenv("ENVIRONMENT"),

		"rollbar.token": os.Getenv("ROLLBAR_TOKEN"),

		"newrelic.token":   os.Getenv("NEW_RELIC_TOKEN"),
		"newrelic.appName": os.Getenv("NEW_RELIC_APP_NAME"),
	}
	return settings
}

func initNewRelic(env string, appName string, token string) newrelic.Application {
	if env != "production" {
		appName = fmt.Sprintf("%s (%s)", appName, env)
	}

	cfg := newrelic.NewConfig(appName, token)
	app, err := newrelic.NewApplication(cfg)
	if err != nil {
		log.Fatal("Error initialize New Relic Agent:", err)
	}
	return app
}

func RequestLogger() gin.HandlerFunc {
	readBody := func(reader io.Reader) string {
		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)
		return buf.String()
	}

	return func(c *gin.Context) {
		t := time.Now()

		buf, _ := ioutil.ReadAll(c.Request.Body)
		logReader := ioutil.NopCloser(bytes.NewBuffer(buf))
		newReader := ioutil.NopCloser(bytes.NewBuffer(buf))
		c.Request.Body = newReader
		c.Next()

		event_string := `event="request"`
		path_string := fmt.Sprintf(`path="%s"`, c.Request.URL.Path)
		latency := time.Since(t)
		latency_string := fmt.Sprintf(`time_elapsed=%v`, latency)
		params := readBody(logReader)
		params_string := fmt.Sprintf(`params="%s"`, params)

		log.Println(event_string, path_string, latency_string, params_string)
	}
}

func initDB() *gorm.DB {
	db, openErr := gorm.Open("sqlite3", "./gorm.db")
	if openErr != nil {
		rollbar.Error(rollbar.ERR, openErr)
		log.Fatal("Error opening database:", openErr)
	}
	if migrateErr := db.AutoMigrate(&PaymentEvent{}).Error; migrateErr != nil {
		rollbar.Error(rollbar.ERR, migrateErr)
		log.Fatal("Error migrating database:", migrateErr)
	}
	db.LogMode(false)
	return db
}

type Event struct {
	EventType  string `json:"event_type"`
	EntityID   string `json:"entity_id"`
	LocationID string `json:"location_id"`
	MerchantID string `json:"merchant_id"`
}

type PaymentEvent struct {
	ID         string    `json:"id" binding:"required" gorm:"primary_key;not null;"`
	LocationID string    `json:"location_id" binding:"required" gorm:"primary_key;not null"`
	CreatedAt  time.Time `json:"created_at"`
}

type Controller struct {
	db *gorm.DB
}

func (ctrl Controller) HandleEvent(c *gin.Context) {
	var event Event
	c.ShouldBindJSON(&event)

	switch event.EventType {
	case "PAYMENT_UPDATED":
		if wrapError(c, 422, ctrl.savePaymentIfNew(event.EntityID, event.LocationID)) {
			return
		}
	}

	c.JSON(200, event)
}

func wrapError(c *gin.Context, status int, err error) bool {
	if err == nil {
		return false
	} else {
		rollbar.Error(rollbar.ERR, err)
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return true
	}
}

func (ctrl Controller) savePaymentIfNew(paymentID string, locationID string) error {
	paymentEvent := PaymentEvent{
		ID:         paymentID,
		LocationID: locationID,
		CreatedAt:  time.Now(),
	}
	return ctrl.db.Create(&paymentEvent).Error
}
