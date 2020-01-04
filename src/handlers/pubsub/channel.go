package pubsub

import (
	"../../utils"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

type dataChanel struct {
	bytes  *[]byte
	length int64
}

type channelWrap struct {
	data           chan dataChanel
	com            chan struct{}
	maxReqSizeInMb int64
}

func (channelWrap channelWrap) consume(request *http.Request, responseWriter http.ResponseWriter) {
	if channelWrap.com != nil {
		channelWrap.com <- struct{}{}
	}

	select {
	case dataChanel := <-channelWrap.data:
		newBytes := make([]byte, dataChanel.length)
		copy(newBytes, *dataChanel.bytes)
		buffer := bytes.NewBuffer(newBytes)
		responseWriter.Header().Set("Content-Length", strconv.FormatInt(dataChanel.length, 10))
		responseWriter.Header().Set("Content-Type", "application/octet-stream")
		_, err := io.Copy(responseWriter, buffer)
		utils.LogError(err, request)

	case <-request.Context().Done():
	}
}

func (channelWrap channelWrap) produce(request *http.Request, responseWriter http.ResponseWriter) {
	if request.ContentLength <= 0 {
		http.Error(responseWriter, "no content", http.StatusBadRequest)
		return
	}

	maxReqSizeInByte := channelWrap.maxReqSizeInMb * 1000 * 1000
	if request.ContentLength > maxReqSizeInByte {
		maxReqSizeInByteStr := strconv.FormatInt(maxReqSizeInByte, 10)
		reqContentLenStr := strconv.FormatInt(request.ContentLength, 10)
		errorMsg := "max. request size is " + maxReqSizeInByteStr + ", got " + reqContentLenStr
		http.Error(responseWriter, errorMsg, http.StatusRequestEntityTooLarge)
		return
	}

	consumersCount := channelWrap.getConsumerCount()

	if consumersCount == 0 {
		http.Error(responseWriter, "no consumers", http.StatusPreconditionFailed)
		return
	}

	bodyBytes, err := ioutil.ReadAll(request.Body)

	if err != nil {
		utils.LogError(err, request)
		return
	}

	channelWrap.sendDataToConsumers(consumersCount, bodyBytes, request.ContentLength, request)
}
