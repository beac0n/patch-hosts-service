package mpmc

import (
	"../../utils"
	"net/http"
	"strconv"
)

func (requestHandler *RequestHandler) consume(request *http.Request, responseWriter http.ResponseWriter, dataChannel chan *[]byte) {
	select {
	case bytes := <-dataChannel:
		responseWriter.Header().Set("Content-Length", strconv.Itoa(len(*bytes)))
		_, err := responseWriter.Write(*bytes)
		utils.LogError(err, request)

	case <-request.Context().Done():
	}
}
