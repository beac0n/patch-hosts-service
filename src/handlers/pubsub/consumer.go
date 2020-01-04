package pubsub

import (
	"../../utils"
	"net/http"
	"strconv"
)

func (requestHandler *RequestHandler) consume(comChannel chan struct{}, dataChannel chan *[]byte, responseWriter http.ResponseWriter, request *http.Request) {
	comChannel <- struct{}{}

	select {
	case bytes := <-dataChannel:
		responseWriter.Header().Set("Content-Length", strconv.Itoa(len(*bytes)))
		_, err := responseWriter.Write(*bytes)
		utils.LogError(err, request)

	case <-request.Context().Done():
	}
}
