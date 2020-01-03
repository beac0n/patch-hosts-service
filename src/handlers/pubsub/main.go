package pubsub

import "net/http"

type RequestHandler struct {
	maxReqSizeInMb int64
}

func NewRequestHandler(maxReqSizeInMb int64) *RequestHandler {
	return &RequestHandler{maxReqSizeInMb: maxReqSizeInMb}
}


func (requestHandler *RequestHandler) ServeHttp(responseWriter http.ResponseWriter, request *http.Request) {
	if (request.Method != http.MethodGet) && (request.Method != http.MethodPost) {
		http.Error(responseWriter, "wrong http method", http.StatusBadRequest)
		return
	}

	var handler handler

	pubSubKeys, ok := request.URL.Query()["pubsub"]
	if ok && len(pubSubKeys) == 1 && pubSubKeys[0] == "true" {
		handler = newHandlerMulti(request.URL.Path, requestHandler.maxReqSizeInMb)
	} else {
		handler = newHandlerSingle(request.URL.Path, requestHandler.maxReqSizeInMb)
	}

	if request.Method == http.MethodPost {
		handler.handleProducer(request, responseWriter)
	} else if request.Method == http.MethodGet {
		handler.handleConsumer(request, responseWriter)
	}
}
