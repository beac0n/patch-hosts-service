package reqres

import (
	"../../constants"
	"../../utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type reqData struct {
	data   *[]byte
	header http.Header
	host   string
	uri    string
}

func (reqHandler *ReqHandler) produce(req *http.Request, resWriter http.ResponseWriter, dataChanRes chan *reqData, dataChanReq chan *reqData) {
	bodyBytes, ok := reqHandler.getBodyBytes(req, resWriter)
	if !ok {
		return
	}

	if strings.HasPrefix(req.URL.Path, constants.Res) {
		getData(dataChanRes, resWriter, req)
		sendData(dataChanReq, bodyBytes, req)
	} else if strings.HasPrefix(req.URL.Path, constants.Req) {
		sendData(dataChanRes, bodyBytes, req)
		getData(dataChanReq, resWriter, req)
	}
}

func (reqHandler *ReqHandler) getBodyBytes(req *http.Request, resWriter http.ResponseWriter) ([]byte, bool) {
	if utils.HttpErrorRequestEntityTooLarge(reqHandler.maxReqSize, req, resWriter) {
		return nil, false
	}

	bodyBytes, err := ioutil.ReadAll(req.Body)

	if utils.LogError(err, req) {
		return nil, false
	}

	return bodyBytes, true
}

func sendData(dataChan chan *reqData, bodyBytes []byte, req *http.Request) {
	select {
	case dataChan <- &reqData{data: &bodyBytes, header: req.Header, host: req.Host, uri: req.RequestURI}:
	case <-req.Context().Done():
	}
}

func getData(dataChan chan *reqData, resWriter http.ResponseWriter, req *http.Request) {
	select {
	case reqData := <-dataChan:
		resWriter.Header().Set(constants.HeaderPrefix+"0-host", reqData.host)
		resWriter.Header().Set(constants.HeaderPrefix+"0-uri", reqData.uri)
		for key, value := range reqData.header {
			for index, value := range value {
				resWriter.Header().Set(constants.HeaderPrefix+strconv.Itoa(index)+"-"+key, value)
			}
		}

		resWriter.Header().Set(constants.ContentLength, strconv.Itoa(len(*reqData.data)))
		_, err := resWriter.Write(*reqData.data)
		resWriter.(http.Flusher).Flush()
		utils.LogError(err, req)

	case <-req.Context().Done():
	}
}
