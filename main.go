package main

import (
	"github.com/a-ryoo/overseer/cmd"
	log "github.com/sirupsen/logrus"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "not-set"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.Println("Running Overseer")
	log.Printf("Build Version: %s", Version)
	log.Printf("Build Commit: %s", Commit)
	log.Printf("Build Time: %s", Date)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}
