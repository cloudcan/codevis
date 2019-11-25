package main

import (
	"github.com/cloudcan/codevis/config"
	"github.com/cloudcan/codevis/graphdb"
	"log"
)

func main() {
	// init config
	config := config.Init("config.json")
	// init db
	graphdb.Init(config.GraphDB)
	defer graphdb.Close()
	log.Print(config)
}
