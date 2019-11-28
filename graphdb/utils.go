package graphdb

import (
	"fmt"
	"log"
	"strings"
)

func Escape(raw string) string {
	escapeValue := strings.ReplaceAll(raw, "\\", "\\\\")
	escapeValue = strings.ReplaceAll(escapeValue, "'", "\\'")
	escapeValue = strings.ReplaceAll(escapeValue, "\"", "\\\"")
	return escapeValue
}

// create index
func CreateIndex(label string, field ...string) {
	for _, f := range field {
		_, err := Exec(fmt.Sprintf("create index on :%s(%s)", label, f), nil)
		if err != nil {
			log.Printf("create index on :%s(%s) err,cause :%s", label, f, err)
		}
	}
}
