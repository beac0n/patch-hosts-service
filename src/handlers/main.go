package handlers

import "net/http"
import "./pubsub"

func NewPubSubRequestHandler(maxReqSizeInMb int64) *PubSubRequestHandler {
	return &PubSubRequestHandler{maxReqSizeInMb: maxReqSizeInMb}
}

type PubSubRequestHandler struct {
	maxReqSizeInMb int64
}

func (pubSubRequestHandler *PubSubRequestHandler) ServeHttp(responseWriter http.ResponseWriter, request *http.Request) {
	if (request.Method != http.MethodGet) && (request.Method != http.MethodPost) {
		http.Error(responseWriter, "wrong http method", http.StatusBadRequest)
		return
	}

	var pubSubHandler pubsub.Handler

	pubSubKeys, ok := request.URL.Query()["pubsub"]
	if ok && len(pubSubKeys) == 1 && pubSubKeys[0] == "true" {
		pubSubHandler = pubsub.NewHandlerMulti(request.URL.Path, pubSubRequestHandler.maxReqSizeInMb)
	} else {
		pubSubHandler = pubsub.NewHandlerSingle(request.URL.Path, pubSubRequestHandler.maxReqSizeInMb)
	}

	if request.Method == http.MethodPost {
		pubSubHandler.HandleProducer(request, responseWriter)
	} else if request.Method == http.MethodGet {
		pubSubHandler.HandleConsumer(request, responseWriter)
	}
}
