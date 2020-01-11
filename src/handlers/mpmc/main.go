package mpmc

import (
	"../../utils"
	"net/http"
	"sync"
)

type ReqHandler struct {
	maxReqSize  int64
	dataChanMap *sync.Map
}

func NewReqHandler(maxReqSize int64) http.Handler {
	return &ReqHandler{maxReqSize: maxReqSize, dataChanMap: &sync.Map{}}
}

func (reqHandler *ReqHandler) ServeHTTP(resWriter http.ResponseWriter, req *http.Request) {
	if utils.NotGetOrPost(req, resWriter) {
		return
	}

	dataChanCreator := func() interface{} { return make(chan *[]byte) }
	dataChan := utils.LoadAndStore(reqHandler.dataChanMap, req.URL.Path, dataChanCreator).(chan *[]byte)

	if req.Method == http.MethodPost {
		reqHandler.produce(req, resWriter, dataChan)
	} else if req.Method == http.MethodGet {
		reqHandler.consume(req, resWriter, dataChan)
	}
}
