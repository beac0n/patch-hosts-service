package pubsub

import (
	"../../utils"
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var pubSubReqHandler = NewReqHandler(10)
var testData = "test"

func TestServeHttpPersist(t *testing.T) {
	requestHandler := http.HandlerFunc(pubSubReqHandler.ServeHTTP)
	requestRecord := httptest.NewRecorder()

	getRequest := httptest.NewRequest("GET", "/pubsub/t?persist=true", nil)
	go requestHandler.ServeHTTP(requestRecord, getRequest)

	postRequest0 := httptest.NewRequest("POST", "/pubsub/t", bytes.NewBuffer([]byte(testData)))
	utils.AssertRequest("", utils.SendRequestSync(requestHandler, postRequest0), t)

	postRequest1 := httptest.NewRequest("POST", "/pubsub/t", bytes.NewBuffer([]byte(testData)))
	utils.AssertRequest("", utils.SendRequestSync(requestHandler, postRequest1), t)

	time.Sleep(time.Millisecond)
	utils.AssertRequest(testData+testData, requestRecord, t)
}

func TestServeHttpMulti(t *testing.T) {
	requestHandler := http.HandlerFunc(pubSubReqHandler.ServeHTTP)

	numberOfGetRequest := 1000

	requestRecorderChan := make(chan *httptest.ResponseRecorder)
	comChan := make(chan struct{}, numberOfGetRequest)

	log.Println("TestServeHttpMulti", "sending GET reqs", numberOfGetRequest)
	for i := 0; i < numberOfGetRequest; i++ {
		getRequest := httptest.NewRequest("GET", "/pubsub/t", nil)
		go utils.SendRequestWithCom(requestHandler, getRequest, requestRecorderChan, comChan)
	}

	log.Println("TestServeHttpMulti", "waiting for GET reqs to halt")
	for i := 0; i < numberOfGetRequest; i++ {
		select {
		case <-comChan:
		}
	}

	log.Println("TestServeHttpMulti", "sending POST req")
	postRequest := httptest.NewRequest("POST", "/pubsub/t", bytes.NewBuffer([]byte(testData)))
	requestRecorderPost := utils.SendRequestSync(requestHandler, postRequest)

	if !utils.AssertRequest("", requestRecorderPost, t) {
		return
	}

	log.Println("TestServeHttpMulti", "asserting GET reqs")
	for i := 0; i < numberOfGetRequest; i++ {
		if !utils.AssertRequest(testData, <-requestRecorderChan, t) {
			return
		}
	}
}
