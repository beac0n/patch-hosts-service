package pubsub

import "net/http"

type RequestHandler struct {
	maxReqSize int64
}

var multiplesChannelWrap = newChannelMapWrap()
var singlesChannelWrap = newChannelMapWrap()

func NewRequestHandler(maxReqSizeIn int64) *RequestHandler {
	return &RequestHandler{maxReqSize: maxReqSizeIn}
}

func (requestHandler *RequestHandler) ServeHttp(responseWriter http.ResponseWriter, request *http.Request) {
	if (request.Method != http.MethodGet) && (request.Method != http.MethodPost) {
		http.Error(responseWriter, "wrong http method", http.StatusBadRequest)
		return
	}

	var wrap channelWrap

	pubSubKeys, ok := request.URL.Query()["pubsub"]
	if ok && len(pubSubKeys) == 1 && pubSubKeys[0] == "true" {
		wrap = channelWrap{
			data:       multiplesChannelWrap.getDataChannel(request.URL.Path, 0),
			com:        multiplesChannelWrap.getComChannel(request.URL.Path, 100),
			maxReqSize: requestHandler.maxReqSize,
		}
	} else {
		wrap = channelWrap{
			data:       singlesChannelWrap.getDataChannel(request.URL.Path, 0),
			maxReqSize: requestHandler.maxReqSize,
		}
	}

	if request.Method == http.MethodPost {
		wrap.produce(request, responseWriter)
	} else if request.Method == http.MethodGet {
		wrap.consume(request, responseWriter)
	}
}
