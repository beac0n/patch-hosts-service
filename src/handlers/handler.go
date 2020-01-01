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

	bodyBytes, err := ioutil.ReadAll(request.Body)

	if err != nil {
		utils.LogError(err, request)
		return
	}

	for i := consumersCount; i > 0; i-- {
		select {
		case handler.data <- &bodyBytes:
		case <-request.Context().Done():
			return
		}
	}
}
