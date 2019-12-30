package handlers

import (
	"../utils"
	"bytes"
	"io"
	"net/http"
)

func HandleConsumerPubSub(dataChannel chan bytes.Buffer, comChannel chan struct{}, request *http.Request, responseWriter http.ResponseWriter) {
	comChannel <- struct{}{}

	select {
	case buffer := <-dataChannel:
		_, err := io.Copy(responseWriter, &buffer)

		utils.LogError(err, request)

	case <-request.Context().Done():
	}
}

func HandleProducerPubSub(dataChannel chan bytes.Buffer, comChannel chan struct{}, request *http.Request, responseWriter http.ResponseWriter) {
	consumersCount := uint64(0)

	ComLoop:
	for {
		select {
		case <-comChannel:
			consumersCount++
		default:
			break ComLoop
		}
	}

	if consumersCount == 0 {
		responseWriter.WriteHeader(http.StatusPreconditionFailed)
		_, err := responseWriter.Write([]byte("no consumers"))
		utils.LogError(err, request)
		return
	}

	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(request.Body)

	if err != nil {
		utils.LogError(err, request)
		return
	}

	for i := consumersCount; i > 0; i-- {
		select {
		case dataChannel <- *buffer:
		case <-request.Context().Done():
			return
		}
	}
}
