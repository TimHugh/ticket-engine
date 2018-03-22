package main

import (
	"fmt"

	"github.com/namsral/flag"
)

type Config struct {
	store map[string]*string
}

func (c *Config) Load() {
	c.store = map[string]*string{
		"environment":    flag.String("env", "development", "application environment"),
		"http":           flag.String("http", ":8080", "HTTP service address"),
		"newrelic_token": flag.String("new_relic_token", "", "New Relic API token"),
		"newrelic_app":   flag.String("new_relic_app", "", "New Relic application name"),
		"rollbar_token":  flag.String("rollbar_token", "", "Rollbar API token"),
		"mongodb_uri":    flag.String("mongodb_uri", "", "MongoDB host URI"),
	}
	flag.Parse()

	for name, val := range c.store {
		fmt.Printf("%s: %s\n", name, *val)
	}
}

func (c *Config) Get(name string) string {
	v, ok := c.store[name]
	if !ok {
		return ""
	}
	return *v
}
