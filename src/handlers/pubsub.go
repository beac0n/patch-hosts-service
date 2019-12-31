package handlers

import (
	"../utils"
	"bytes"
	"io"
	"net/http"
)

type HandlerPubSub struct {
	data chan bytes.Buffer
	com  chan struct{}
}

func (handler HandlerPubSub) HandleConsumer(request *http.Request, responseWriter http.ResponseWriter) {
	if handler.com != nil {
		handler.com <- struct{}{}
	}

	select {
	case buffer := <-handler.data:
		_, err := io.Copy(responseWriter, &buffer)

		utils.LogError(err, request)

	case <-request.Context().Done():
	}
}

func (handler HandlerPubSub) HandleProducer(request *http.Request, responseWriter http.ResponseWriter) {
	consumersCount := uint64(0)

	if handler.com != nil {
	ComLoop:
		for {
			select {
			case <-handler.com:
				consumersCount++
			default:
				break ComLoop
			}
		}
	} else {
		consumersCount = 1
	}

	if consumersCount == 0 {
		utils.HttpError(request, responseWriter, http.StatusPreconditionFailed, "no consumers")
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
		case handler.data <- *buffer:
		case <-request.Context().Done():
			return
		}
	}
}
