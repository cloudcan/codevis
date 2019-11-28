package graphdb

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
)

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Uri      string `json:"uri"`
}

var (
	driver  neo4j.Driver
	session neo4j.Session
)

// init db
func Init(config Config) {
	var err error
	driver, err = neo4j.NewDriver(config.Uri, neo4j.BasicAuth(config.Username, config.Password, ""))
	if err != nil {
		log.Fatal("cannot load neo4j driver,cause:", err)
	}
	session, err = driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		log.Fatal("cannot open db session,cause:", err)
	}
}

// close db
func Close() {
	if session != nil {
		_ = session.Close()
	}
	if driver != nil {
		_ = driver.Close()
	}
}

// exec cypher QL
func Exec(cql string, params map[string]interface{}) (neo4j.Result, error) {
	return session.Run(cql, params)
}
