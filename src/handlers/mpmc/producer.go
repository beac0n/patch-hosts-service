package mpmc

import (
	"../../utils"
	"io/ioutil"
	"net/http"
)

func (requestHandler *RequestHandler) produce(request *http.Request, responseWriter http.ResponseWriter, dataChannel chan *[]byte) {
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

	select {
	case dataChannel <- &bodyBytes:
	case <-request.Context().Done():
	}

}
