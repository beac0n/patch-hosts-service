package mpmc

import (
	"../../utils"
	"net/http"
	"strconv"
)

func (requestHandler *RequestHandler) consume(dataChannel chan *[]byte, responseWriter http.ResponseWriter, request *http.Request) {
	select {
	case bytes := <-dataChannel:
		responseWriter.Header().Set("Content-Length", strconv.Itoa(len(*bytes)))
		_, err := responseWriter.Write(*bytes)
		utils.LogError(err, request)

	case <-request.Context().Done():
	}
}
