package config

import (
	"encoding/json"
	"errors"
	"github.com/cloudcan/codevis/analysis"
	"github.com/cloudcan/codevis/graphdb"
	"io/ioutil"
	"log"
)

type Config struct {
	GraphDB  graphdb.Config  `json:"graph_db"`
	Analysis analysis.Config `json:"analysis"`
}

func (config *Config) Check() error {
	if config.GraphDB.Uri == "" || config.GraphDB.Username == "" || config.GraphDB.Password == "" {
		return errors.New("graph db config is not valid")
	}
	return nil
}

func Load(path string) (config *Config) {
	config = new(Config)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("cannot load config file, cause:", err)
	}
	err = json.Unmarshal(bytes, config)
	if err != nil {
		log.Fatal("invalid config file, cause:", err)
	}

	if err = config.Check(); err != nil {
		log.Fatal(err)
	}
	log.Print("load config from :", path)
	return
}
