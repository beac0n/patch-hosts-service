package pubsub

import (
	"../../utils"
	"io/ioutil"
	"net/http"
)

func (requestHandler *RequestHandler) produce(request *http.Request, responseWriter http.ResponseWriter, dataChannel chan *[]byte, comChannel chan struct{}) {
	if request.ContentLength <= 0 {
		http.Error(responseWriter, "no content", http.StatusBadRequest)
		return
	}

	if utils.HttpErrorRequestEntityTooLarge(requestHandler.maxReqSize, request, responseWriter) {
		return
	}

	bodyBytes, err := ioutil.ReadAll(request.Body)

	if utils.LogError(err, request) {
		return
	}

	if consumersCount := getConsumerCount(comChannel); consumersCount > 0 {
		sendDataToConsumers(consumersCount, &bodyBytes, dataChannel, request)
	} else {
		http.Error(responseWriter, "no consumers", http.StatusPreconditionFailed)
	}
}

func getConsumerCount(comChannel chan struct{}) uint64 {
	consumersCount := uint64(0)

	for {
		select {
		case <-comChannel:
			consumersCount++
		default:
			return consumersCount
		}
	}
}

func sendDataToConsumers(consumersCount uint64, bodyBytes *[]byte, dataChannel chan *[]byte, request *http.Request) {
	for ; consumersCount > 0; consumersCount-- {
		select {
		case dataChannel <- bodyBytes:
		case <-request.Context().Done():
			return
		}
	}
}
