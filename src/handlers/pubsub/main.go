package pubsub

import (
	"math"
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

	persistKeys, persistOk := request.URL.Query()["persist"]
	persist := persistOk && len(persistKeys) == 1 && persistKeys[0] == "true"

	if _, dataChannelOk := requestHandler.dataChannelMap.Load(request.URL.Path); !dataChannelOk {
		requestHandler.dataChannelMap.Store(request.URL.Path, make(chan *[]byte))
	}

	dataChannelI, _ := requestHandler.dataChannelMap.Load(request.URL.Path)
	dataChannel := dataChannelI.(chan *[]byte)

	if _, comChannelOk := requestHandler.comChannelMap.Load(request.URL.Path); !comChannelOk {
		requestHandler.comChannelMap.Store(request.URL.Path, make(chan struct{}, math.MaxInt64))
	}

	comChannelI, _ := requestHandler.comChannelMap.Load(request.URL.Path)
	comChannel := comChannelI.(chan struct{})

	if request.Method == http.MethodPost {
		requestHandler.produce(request, responseWriter, dataChannel, comChannel)
	} else if request.Method == http.MethodGet {
		requestHandler.consume(request, responseWriter, dataChannel, comChannel, persist)
	}
}
