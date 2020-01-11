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

func TestReqHandler_ServeHTTP_wrongMethod(t *testing.T) {
	postReq := httptest.NewRequest("PUT", "/foobar", bytes.NewBuffer([]byte("")))

	recorder := utils.SendRequestSync(pubSubReqHandler, postReq)
	utils.Assert(t, recorder.Code, 400)
}

func TestReqHandler_ServeHTTP_noConsumers(t *testing.T) {
	postReq := httptest.NewRequest("POST", "/foobar", bytes.NewBuffer([]byte(testData)))

	recorder := utils.SendRequestSync(pubSubReqHandler, postReq)
	utils.Assert(t, recorder.Code, http.StatusPreconditionFailed)
}


func TestReqHandler_ServeHTTP_contentLengthError(t *testing.T) {
	postReq := httptest.NewRequest("POST", "/foobar", bytes.NewBuffer([]byte("")))

	recorder := utils.SendRequestSync(pubSubReqHandler, postReq)
	utils.Assert(t, recorder.Code, http.StatusBadRequest)
}

func TestReqHandler_ServeHTTP_reqEntityTooLarge(t *testing.T) {
	pubSubReqHandler := NewReqHandler(1)
	postReq := httptest.NewRequest("POST", "/foobar", bytes.NewBuffer([]byte("toolarge")))

	recorder := utils.SendRequestSync(pubSubReqHandler, postReq)
	utils.Assert(t, recorder.Code, http.StatusRequestEntityTooLarge)
}

func TestReqHandler_ServeHTTP_persist(t *testing.T) {
	requestRecord := httptest.NewRecorder()

	getRequest := httptest.NewRequest("GET", "/pubsub/t?persist=true", nil)
	go pubSubReqHandler.ServeHTTP(requestRecord, getRequest)

	postRequest0 := httptest.NewRequest("POST", "/pubsub/t", bytes.NewBuffer([]byte(testData)))
	utils.AssertRequest(nil, utils.SendRequestSync(pubSubReqHandler, postRequest0), t)

	postRequest1 := httptest.NewRequest("POST", "/pubsub/t", bytes.NewBuffer([]byte(testData)))
	utils.AssertRequest(nil, utils.SendRequestSync(pubSubReqHandler, postRequest1), t)

	time.Sleep(time.Millisecond)
	utils.AssertRequest(testData+testData, requestRecord, t)
}

func TestReqHandler_ServeHTTP_multi(t *testing.T) {
	numberOfGetRequest := 1000

	requestRecorderChan := make(chan *httptest.ResponseRecorder)
	comChan := make(chan struct{}, numberOfGetRequest)

	log.Println("TestReqHandler_ServeHTTP_multi", "sending GET reqs", numberOfGetRequest)
	for i := 0; i < numberOfGetRequest; i++ {
		getRequest := httptest.NewRequest("GET", "/pubsub/t", nil)
		go utils.SendRequestWithCom(pubSubReqHandler, getRequest, requestRecorderChan, comChan)
		time.Sleep(time.Millisecond)
	}

	log.Println("TestReqHandler_ServeHTTP_multi", "waiting for GET reqs to halt")
	for i := 0; i < numberOfGetRequest; i++ {
		select {
		case <-comChan:
		}
	}

	log.Println("TestReqHandler_ServeHTTP_multi", "sending POST req")
	postRequest := httptest.NewRequest("POST", "/pubsub/t", bytes.NewBuffer([]byte(testData)))
	requestRecorderPost := utils.SendRequestSync(pubSubReqHandler, postRequest)

	if !utils.AssertRequest(nil, requestRecorderPost, t) {
		return
	}

	log.Println("TestReqHandler_ServeHTTP_multi", "asserting GET reqs")
	for i := 0; i < numberOfGetRequest; i++ {
		if !utils.AssertRequest(testData, <-requestRecorderChan, t) {
			return
		}
	}
}
