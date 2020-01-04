package pubsub

import (
	"net/http"
	"sync"
)

type RequestHandler struct {
	maxReqSize     int64
	dataChannelMap *sync.Map
	comChannelMap  *sync.Map
}

func NewRequestHandler(maxReqSize int64) *RequestHandler {
	return &RequestHandler{maxReqSize: maxReqSize, dataChannelMap: &sync.Map{}, comChannelMap: &sync.Map{}}
}

func (requestHandler *RequestHandler) ServeHttp(responseWriter http.ResponseWriter, request *http.Request) {
	if (request.Method != http.MethodGet) && (request.Method != http.MethodPost) {
		http.Error(responseWriter, "wrong http method", http.StatusBadRequest)
		return
	}

	dataChannelI, _ := requestHandler.dataChannelMap.LoadOrStore(request.URL.Path, make(chan *[]byte))
	dataChannel := dataChannelI.(chan *[]byte)

	comChannelI, _ := requestHandler.comChannelMap.LoadOrStore(request.URL.Path, make(chan struct{}, 100))
	comChannel := comChannelI.(chan struct{})

	if request.Method == http.MethodPost {
		requestHandler.produce(request, responseWriter, comChannel, dataChannel)
	} else if request.Method == http.MethodGet {
		requestHandler.consume(comChannel, dataChannel, responseWriter, request)
	}
}
