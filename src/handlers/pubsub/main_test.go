package pubsub

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var pubSubReqHandler = NewRequestHandler(10)
var testData = "test"

func TestServeHttpPersist(test *testing.T) {
	requestHandler := http.HandlerFunc(pubSubReqHandler.ServeHttp)
	requestRecord := httptest.NewRecorder()

	getRequest, _ := http.NewRequest("GET", "/pubsub/test?persist=true", nil)
	go requestHandler.ServeHTTP(requestRecord, getRequest)

	postRequest0, _ := http.NewRequest("POST", "/pubsub/test", bytes.NewBuffer([]byte(testData)))
	assertRequest("", sendRequestSync(requestHandler, postRequest0), test)

	postRequest1, _ := http.NewRequest("POST", "/pubsub/test", bytes.NewBuffer([]byte(testData)))
	assertRequest("", sendRequestSync(requestHandler, postRequest1), test)

	time.Sleep(time.Millisecond)
	assertRequest(testData+testData, requestRecord, test)
}

func TestServeHttpMulti(test *testing.T) {
	requestHandler := http.HandlerFunc(pubSubReqHandler.ServeHttp)

	numberOfGetRequest := 1000

	requestRecorderChan := make(chan *httptest.ResponseRecorder)
	comChan := make(chan struct{}, numberOfGetRequest)

	log.Println("TestServeHttpMulti", "sending GET reqs", numberOfGetRequest)
	for i := 0; i < numberOfGetRequest; i++ {
		getRequest, _ := http.NewRequest("GET", "/pubsub/test", nil)
		go sendRequest(requestHandler, getRequest, requestRecorderChan, comChan)
	}

	log.Println("TestServeHttpMulti", "waiting for GET reqs to halt")
	for i := 0; i < numberOfGetRequest; i++ {
		select {
		case <-comChan:
		}
	}

	log.Println("TestServeHttpMulti", "sending POST req")
	postRequest, _ := http.NewRequest("POST", "/pubsub/test", bytes.NewBuffer([]byte(testData)))
	requestRecorderPost := sendRequestSync(requestHandler, postRequest)

	if !assertRequest("", requestRecorderPost, test) {
		return
	}

	log.Println("TestServeHttpMulti", "asserting GET reqs")
	for i := 0; i < numberOfGetRequest; i++ {
		if !assertRequest(testData, <-requestRecorderChan, test) {
			return
		}
	}
}

func sendRequestSync(requestHandler http.HandlerFunc, postRequest *http.Request) *httptest.ResponseRecorder {
	// we need to wait a minimum amount of time, or else requests will be faster than writing to channel
	time.Sleep(time.Millisecond)

	requestRecorderPost := httptest.NewRecorder()
	requestHandler.ServeHTTP(requestRecorderPost, postRequest)
	return requestRecorderPost
}

func sendRequest(reqHandler http.HandlerFunc, req *http.Request, reqRecChan chan *httptest.ResponseRecorder, comChan chan struct{}) {
	requestRecord := httptest.NewRecorder()

	comChan <- struct{}{}

	reqHandler.ServeHTTP(requestRecord, req)

	reqRecChan <- requestRecord
}

func assertRequest(expected string, requestRecorder *httptest.ResponseRecorder, test *testing.T) bool {
	if status := requestRecorder.Code; status != http.StatusOK {
		test.Errorf("response has wrong status code: got %v want %v ", status, http.StatusOK)
		return false
	}

	if expected == "" {
		return true
	}

	if actual := requestRecorder.Body.String(); actual != expected {
		test.Errorf("response has unexpected body: got %v want %v", actual, expected)
		return false
	}

	return true
}
