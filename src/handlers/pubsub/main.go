package pubsub

import (
	"../../utils"
	"math"
	"net/http"
	"sync"
)

type RequestHandler struct {
	maxReqSize  int64
	dataChanMap *sync.Map
	comChanMap  *sync.Map
	muxMap      *sync.Map
}

func NewRequestHandler(maxReqSize int64) *RequestHandler {
	return &RequestHandler{maxReqSize, &sync.Map{}, &sync.Map{}, &sync.Map{}}
}

func (reqHandler *RequestHandler) ServeHttp(resWriter http.ResponseWriter, req *http.Request) {
	if (req.Method != http.MethodGet) && (req.Method != http.MethodPost) {
		http.Error(resWriter, "wrong http method", http.StatusBadRequest)
		return
	}

	persistKeys, persistOk := req.URL.Query()["persist"]
	persist := persistOk && len(persistKeys) == 1 && persistKeys[0] == "true"

	dataChanCreator := func() interface{} { return make(chan *[]byte) }
	dataChan :=  utils.LoadAndStore(reqHandler.dataChanMap, req.URL.Path, dataChanCreator).(chan *[]byte)

	comChanCreator := func() interface{} { return make(chan struct{}, math.MaxInt64) }
	comChan :=  utils.LoadAndStore(reqHandler.comChanMap, req.URL.Path, comChanCreator).(chan struct{})

	muxCreator := func() interface{} { return &sync.Mutex{} }
	mux := utils.LoadAndStore(reqHandler.muxMap, req.URL.Path, muxCreator).(*sync.Mutex)

	if req.Method == http.MethodPost {
		reqHandler.produce(req, resWriter, dataChan, comChan, mux)
	} else if req.Method == http.MethodGet {
		reqHandler.consume(req, resWriter, dataChan, comChan, persist)
	}
}
