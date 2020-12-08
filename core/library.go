package core

import "time"

type Library struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	StartCount  string    `yaml:"star_count"`
	PullCount   int       `yaml:"pull_count"`
	LastUpdated time.Time `yaml:"last_updated"`
}

