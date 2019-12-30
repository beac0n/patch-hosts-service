package handlers

import (
	"../utils"
	"bytes"
	"io"
	"net/http"
)

func HandleConsumerDefault(dataChannel chan bytes.Buffer, comChannel chan struct{}, request *http.Request, responseWriter http.ResponseWriter) {
	select {
	case buffer := <-dataChannel:
		_, err := io.Copy(responseWriter, &buffer)

		utils.LogError(err, request)

	case <-request.Context().Done():
	}
}

func HandleProducerDefault(dataChannel chan bytes.Buffer, comChannel chan struct{}, request *http.Request, responseWriter http.ResponseWriter) {
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(request.Body)

	if err != nil {
		utils.LogError(err, request)
		return
	}

	select {
	case dataChannel <- *buffer:
	case <-request.Context().Done():
	}
}
