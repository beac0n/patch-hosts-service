package mpmc

import (
	"net/http"
	"sync"
)

type RequestHandler struct {
	maxReqSize     int64
	dataChannelMap *sync.Map
}

func NewRequestHandler(maxReqSize int64) *RequestHandler {
	return &RequestHandler{maxReqSize: maxReqSize, dataChannelMap: &sync.Map{}}
}

func (requestHandler *RequestHandler) ServeHttp(responseWriter http.ResponseWriter, request *http.Request) {
	if (request.Method != http.MethodGet) && (request.Method != http.MethodPost) {
		http.Error(responseWriter, "wrong http method", http.StatusBadRequest)
		return
	}

	dataChannelI, _ := requestHandler.dataChannelMap.LoadOrStore(request.URL.Path, make(chan *[]byte))
	dataChannel := dataChannelI.(chan *[]byte)

	if request.Method == http.MethodPost {
		requestHandler.produce(request, responseWriter, dataChannel)
	} else if request.Method == http.MethodGet {
		requestHandler.consume(dataChannel, responseWriter, request)
	}
}
