package main

import (
	"./handlers"
	"flag"
	"log"
	"net/http"
)

func requestHandler(maxReqSizeInMb int64) func(responseWriter http.ResponseWriter, request *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if (request.Method != http.MethodGet) && (request.Method != http.MethodPost) {
			http.Error(responseWriter, "wrong http method", http.StatusBadRequest)
			return
		}

		var handler handlers.Handler

		pubSubKeys, ok := request.URL.Query()["pubsub"]
		if ok && len(pubSubKeys) == 1 && pubSubKeys[0] == "true" {
			handler = handlers.NewHandlerPubSub(request.URL.Path, maxReqSizeInMb)
		} else {
			handler = handlers.NewHandlerStandard(request.URL.Path, maxReqSizeInMb)
		}

		if request.Method == http.MethodPost {
			handler.HandleProducer(request, responseWriter)
		} else if request.Method == http.MethodGet {
			handler.HandleConsumer(request, responseWriter)
		}
	}
}

func main() {
	host := flag.String("host", "0.0.0.0:9001", "host and port where this application should run")
	maxReqSizeInMb := flag.Int64("max_req_size", 10, "maximum request size in MB")

	flag.Parse()

	log.Println("running on", *host)

	if err := http.ListenAndServe(*host, http.HandlerFunc(requestHandler(*maxReqSizeInMb))); err != nil {
		log.Fatal("FATAL ERROR:", err)
	}
}
