package utils

import (
	"log"
	"net/http"
	"runtime/debug"
)

func LogError(err error, request *http.Request) {
	if err != nil {
		log.Println("ERROR", request.Method, request.URL.Path, err)
		debug.PrintStack()
	}
}
