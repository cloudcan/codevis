package config

import (
	"encoding/json"
	"github.com/cloudcan/codevis/graphdb"
	"io/ioutil"
	"log"
)

type Config struct {
	GraphDB graphdb.Config `json:"graph_db"`
}

func Init(path string) (config *Config) {
	config = new(Config)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("cannot load config file, cause:%s", err)
	}
	err = json.Unmarshal(bytes, config)
	if err != nil {
		log.Fatalf("invalid config file, cause:%s", err)
	}
	log.Printf("load config from %s", path)
	return
}
