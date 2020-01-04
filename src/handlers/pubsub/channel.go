package pubsub

import (
	"../../utils"
	"io/ioutil"
	"net/http"
	"strconv"
)

type channelWrap struct {
	data       chan *[]byte
	com        chan struct{}
	maxReqSize int64
}

func (channelWrap channelWrap) consume(request *http.Request, responseWriter http.ResponseWriter) {
	if channelWrap.com != nil {
		channelWrap.com <- struct{}{}
	}

	select {
	case bytes := <-channelWrap.data:
		responseWriter.Header().Set("Content-Length", strconv.Itoa(len(*bytes)))
		_, err := responseWriter.Write(*bytes)
		utils.LogError(err, request)

	case <-request.Context().Done():
	}
}

func (channelWrap channelWrap) produce(request *http.Request, responseWriter http.ResponseWriter) {
	if request.ContentLength <= 0 {
		http.Error(responseWriter, "no content", http.StatusBadRequest)
		return
	}

	if request.ContentLength > channelWrap.maxReqSize {
		channelWrap.httpErrorEntityTooLarge(request, responseWriter)
		return
	}

	consumersCount := channelWrap.getConsumerCount()

	if consumersCount == 0 {
		http.Error(responseWriter, "no consumers", http.StatusPreconditionFailed)
		return
	}

	bodyBytes, err := ioutil.ReadAll(request.Body)

	if err != nil {
		utils.LogError(err, request)
		return
	}

	channelWrap.sendDataToConsumers(consumersCount, &bodyBytes, request)
}
