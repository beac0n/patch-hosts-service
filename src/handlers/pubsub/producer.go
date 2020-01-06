package pubsub

import (
	"../../utils"
	"io/ioutil"
	"net/http"
	"sync"
)

func (reqHandler *ReqHandler) produce(req *http.Request, resWriter http.ResponseWriter, dataChan chan *[]byte, comChan chan struct{}, mux *sync.Mutex) {
	if req.ContentLength <= 0 {
		http.Error(resWriter, "no content", http.StatusBadRequest)
		return
	}

	if utils.HttpErrorRequestEntityTooLarge(reqHandler.maxReqSize, req, resWriter) {
		return
	}

	bodyBytes, err := ioutil.ReadAll(req.Body)

	if utils.LogError(err, req) {
		return
	}

	// only one producer is allowed to send to the currently listening consumers
	mux.Lock()
	defer mux.Unlock()
	if consumersCount := getConsumerCount(comChan); consumersCount > 0 {
		sendDataToConsumers(consumersCount, &bodyBytes, dataChan, req)
	} else {
		http.Error(resWriter, "no consumers", http.StatusPreconditionFailed)
	}
}

func getConsumerCount(comChan chan struct{}) uint64 {
	consumersCount := uint64(0)

	for {
		select {
		case <-comChan:
			consumersCount++
		default:
			return consumersCount
		}
	}
}

func sendDataToConsumers(consumersCount uint64, bodyBytes *[]byte, dataChan chan *[]byte, req *http.Request) {
	for ; consumersCount > 0; consumersCount-- {
		select {
		case dataChan <- bodyBytes:
		case <-req.Context().Done():
			return
		}
	}
}
