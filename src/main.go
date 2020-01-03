package main

import (
	"./handlers/pubsub"
	"flag"
	"log"
	"net/http"
)

type RequestHandler struct {
	maxReqSizeInMb int64
}

func (requestHandler *RequestHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if (request.Method != http.MethodGet) && (request.Method != http.MethodPost) {
		http.Error(responseWriter, "wrong http method", http.StatusBadRequest)
		return
	}

	var pubSubHandler pubsub.Handler

	pubSubKeys, ok := request.URL.Query()["pubsub"]
	if ok && len(pubSubKeys) == 1 && pubSubKeys[0] == "true" {
		pubSubHandler = pubsub.NewHandlerPubSub(request.URL.Path, requestHandler.maxReqSizeInMb)
	} else {
		pubSubHandler = pubsub.NewHandlerStandard(request.URL.Path, requestHandler.maxReqSizeInMb)
	}

	if request.Method == http.MethodPost {
		pubSubHandler.HandleProducer(request, responseWriter)
	} else if request.Method == http.MethodGet {
		pubSubHandler.HandleConsumer(request, responseWriter)
	}
}

func main() {
	host := flag.String("host", "0.0.0.0:9001", "host and port where this application should run")
	maxReqSizeInMb := flag.Int64("max_req_size", 10, "maximum request size in MB")

	flag.Parse()

	log.Println("running on", *host)

	requestHandler := &RequestHandler{maxReqSizeInMb: *maxReqSizeInMb}

	if err := http.ListenAndServe(*host, requestHandler); err != nil {
		log.Fatal("FATAL ERROR:", err)
	}
}
