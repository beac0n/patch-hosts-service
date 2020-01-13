package reqres

import (
	"net/http"
	"patch-hosts-service/src/constants"
	"patch-hosts-service/src/utils"
	"sync"
)

type httpMethodDataChannels struct {
	httpMethod      string
	dataChanMapFrom *sync.Map
	dataChanMapTo   *sync.Map
}

type ReqHandler struct {
	maxReqSize                int64
	httpMethodDataChannelsMap [len(constants.HttpMethods)]*httpMethodDataChannels
}

func NewReqHandler(maxReqSize int64) http.Handler {
	reqHandler := &ReqHandler{maxReqSize: maxReqSize}

	for i := 0; i < len(constants.HttpMethods); i++ {
		reqHandler.httpMethodDataChannelsMap[i] = &httpMethodDataChannels{
			httpMethod:      constants.HttpMethods[i],
			dataChanMapFrom: &sync.Map{},
			dataChanMapTo:   &sync.Map{},
		}
	}

	return reqHandler
}

func (reqHandler *ReqHandler) ServeHTTP(resWriter http.ResponseWriter, req *http.Request) {
	var dataChannels *httpMethodDataChannels

	for i := 0; i < len(constants.HttpMethods); i++ {
		if req.Method == reqHandler.httpMethodDataChannelsMap[i].httpMethod {
			dataChannels = reqHandler.httpMethodDataChannelsMap[i]
			break
		}
	}

	if dataChannels == nil {
		// this will only happen if the HTTP protocol gets a new method
		http.Error(resWriter, req.Method+(" not implemented yet"), http.StatusNotImplemented)
		return
	}

	dataChanKey := req.URL.Path[4:]
	chanCreator := func() interface{} { return make(chan *reqData) }
	dataChanRes := utils.LoadAndStore(dataChannels.dataChanMapFrom, dataChanKey, chanCreator).(chan *reqData)
	dataChanReq := utils.LoadAndStore(dataChannels.dataChanMapTo, dataChanKey, chanCreator).(chan *reqData)

	reqHandler.produce(req, resWriter, dataChanRes, dataChanReq)
}
