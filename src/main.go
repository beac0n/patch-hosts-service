package main

import (
	"./handlers"
	"log"
	"net/http"
)

var address = ":9001"

func requestHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if (request.Method != http.MethodGet) && (request.Method != http.MethodPost) {
		http.Error(responseWriter, "wrong http method", http.StatusBadRequest)
		return
	}

	var handler handlers.Handler

	pubSubKeys, ok := request.URL.Query()["pubsub"]
	if ok && len(pubSubKeys) == 1 && pubSubKeys[0] == "true" {
		handler = handlers.NewHandlerPubSub(request.URL.Path)
	} else {
		handler = handlers.NewHandlerStandard(request.URL.Path)
	}

	if request.Method == http.MethodPost {
		handler.HandleProducer(request, responseWriter)
	} else if request.Method == http.MethodGet {
		handler.HandleConsumer(request, responseWriter)
	}
}

func main() {
	log.Println("running on", address)

	if err := http.ListenAndServe(address, http.HandlerFunc(requestHandler)); err != nil {
		log.Fatal("FATAL ERROR:", err)
	}
}
