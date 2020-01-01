package handlers

import (
	"../utils"
	"io/ioutil"
	"net/http"
)

func (handler Handler) HandleConsumer(request *http.Request, responseWriter http.ResponseWriter) {
	if handler.com != nil {
		handler.com <- struct{}{}
	}

	select {
	case bodyBytes := <-handler.data:
		_, err := responseWriter.Write(*bodyBytes)
		utils.LogError(err, request)

	case <-request.Context().Done():
	}
}

func (handler Handler) HandleProducer(request *http.Request, responseWriter http.ResponseWriter) {
	if request.ContentLength <= 0 {
		http.Error(responseWriter, "no content", http.StatusBadRequest)
		return
	}

	if request.ContentLength > 10*1000*1000 {
		http.Error(responseWriter, "", http.StatusRequestEntityTooLarge)
		return
	}

	consumersCount := handler.getConsumerCount()

	if consumersCount == 0 {
		http.Error(responseWriter, "no consumers", http.StatusPreconditionFailed)
		return
	}

	bodyBytes, err := ioutil.ReadAll(request.Body)

	if err != nil {
		utils.LogError(err, request)
		return
	}

	handler.sendDataToConsumers(consumersCount, bodyBytes, request)
}

func (handler Handler) sendDataToConsumers(consumersCount uint64, bodyBytes []byte, request *http.Request) {
	for ; consumersCount > 0; consumersCount-- {
		select {
		case handler.data <- &bodyBytes:
		case <-request.Context().Done():
			return
		}
	}
}

func (handler Handler) getConsumerCount() uint64 {
	if handler.com == nil {
		return 1
	}

	consumersCount := uint64(0)
ComLoop:
	for {
		select {
		case <-handler.com:
			consumersCount++
		default:
			break ComLoop
		}
	}

	return consumersCount
}
