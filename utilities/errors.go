package utilities

import (
	"fmt"
	"log"
	"time"
)

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func LogStrErr(prefix string, text string) {
	log.Panicf("%s [%s] %s\n", prefix, returnPrettyDate(), text)
}

func returnPrettyDate() string {
	return fmt.Sprintf("%d/%d/%d %d:%d:%d", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
}

func RequestErrorInputData() map[string]string {
	return map[string]string{
		"status":  "err",
		"details": "Database query execution error",
	}
}

func RequestErrorQueryPg() map[string]string {
	return map[string]string{
		"status":  "err",
		"details": "Database query execution error",
	}
}
