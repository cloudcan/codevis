package main

import (
	"github.com/cloudcan/codevis/analysis"
	"github.com/cloudcan/codevis/config"
	"github.com/cloudcan/codevis/graphdb"
	"log"
)

func main() {
	// load config
	config := config.Load("config.json")
	// init db
	graphdb.Init(config.GraphDB)
	defer graphdb.Close()
	// code analysis
	result, err := analysis.Analysis(config.Analysis)
	if err != nil {
		log.Fatal(err)
	}
	// result refine
	g, err := result.Refine()
	if err != nil {
		log.Fatal(err)
	}
	// save graph
	g.Persistence()
}
