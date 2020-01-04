package pubsub

import (
	"../../utils"
	"net/http"
	"strconv"
)

func (requestHandler *RequestHandler) consume(request *http.Request, responseWriter http.ResponseWriter, dataChannel chan *[]byte, comChannel chan struct{}, persist bool) {
	comChannel <- struct{}{}

	if persist {
		requestHandler.consumePersist(dataChannel, comChannel, responseWriter, request)
	} else {
		requestHandler.consumeNormal(dataChannel, comChannel, responseWriter, request)
	}
}

func (requestHandler *RequestHandler) consumePersist(dataChannel chan *[]byte, comChannel chan struct{}, responseWriter http.ResponseWriter, request *http.Request) {
	for {
		select {
		case bytes := <-dataChannel:
			comChannel <- struct{}{}
			_, err := responseWriter.Write(*bytes)
			responseWriter.(http.Flusher).Flush()
			if utils.LogError(err, request) {
				return
			}

		case <-request.Context().Done():
			<-comChannel
			return
		}
	}
}

func (requestHandler *RequestHandler) consumeNormal(dataChannel chan *[]byte, comChannel chan struct{}, responseWriter http.ResponseWriter, request *http.Request) {
	select {
	case bytes := <-dataChannel:
		responseWriter.Header().Set("Content-Length", strconv.Itoa(len(*bytes)))
		_, err := responseWriter.Write(*bytes)
		utils.LogError(err, request)

	case <-request.Context().Done():
		<-comChannel
	}
}
