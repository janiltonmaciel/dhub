package main

import (
	"log"
	"os"

	"github.com/janiltonmaciel/dhub/cmd"
)

var (
	version string
	commit  string
	date    string
	token   string
)

func main() {
	if err := cmd.Execute(version, commit, date, getToken()); err != nil {
		log.Fatal(err)
	}
}

func getToken() string {
	if token != "" {
		return token
	}
	return os.Getenv("GITHUB_TOKEN")
}
