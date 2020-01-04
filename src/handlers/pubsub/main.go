package pubsub

import "net/http"

type RequestHandler struct {
	maxReqSizeInMb int64
}

var multiplesChannelWrap = newChannelMapWrap()
var singlesChannelWrap = newChannelMapWrap()

func NewRequestHandler(maxReqSizeInMb int64) *RequestHandler {
	return &RequestHandler{maxReqSizeInMb: maxReqSizeInMb}
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
			data:           multiplesChannelWrap.getDataChannel(request.URL.Path, 0),
			com:            multiplesChannelWrap.getComChannel(request.URL.Path, 0),
			maxReqSizeInMb: requestHandler.maxReqSizeInMb,
		}
	} else {
		wrap = channelWrap{
			data:           singlesChannelWrap.getDataChannel(request.URL.Path, 0),
			maxReqSizeInMb: requestHandler.maxReqSizeInMb,
		}
	}

	if request.Method == http.MethodPost {
		wrap.produce(request, responseWriter)
	} else if request.Method == http.MethodGet {
		wrap.consume(request, responseWriter)
	}
}
