package config

import (
	"github.com/namsral/flag"
)

type Config struct {
	store map[string]*ConfigEntry
}

type ConfigEntry struct {
	Name         string
	DefaultValue string
	Description  string
	Value        *string
}

func New() *Config {
	return &Config{
		make(map[string]*ConfigEntry),
	}
}

func (c *Config) Load() {
	for _, entry := range c.store {
		entry.Value = flag.String(entry.Name, entry.DefaultValue, entry.Description)
	}
	flag.Parse()
}

func (c *Config) Get(name string) string {
	entry, ok := c.store[name]
	if !ok {
		return ""
	}
	return *entry.Value
}

func (c *Config) Define(name string, defaultValue string, description string) {
	c.store[name] = &ConfigEntry{
		Name:         name,
		DefaultValue: defaultValue,
		Description:  description,
	}
}
