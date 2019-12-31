package handlers

import (
	"../utils"
	"bytes"
	"io"
	"net/http"
)

type HandlerStandard struct {
	data chan bytes.Buffer
}

func (handler HandlerStandard) HandleConsumer(request *http.Request, responseWriter http.ResponseWriter) {
	select {
	case buffer := <-handler.data:
		_, err := io.Copy(responseWriter, &buffer)

		utils.LogError(err, request)

	case <-request.Context().Done():
	}
}

func (handler HandlerStandard) HandleProducer(request *http.Request, responseWriter http.ResponseWriter) {
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(request.Body)

	if err != nil {
		utils.LogError(err, request)
		return
	}

	select {
	case handler.data <- *buffer:
	case <-request.Context().Done():
	}
}
