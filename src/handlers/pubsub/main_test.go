package pubsub

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var pubSubReqHandler = NewRequestHandler(10)
var testData = "test"

func TestServeHttpMulti(test *testing.T) {
	getRequest, _ := http.NewRequest("GET", "/pubsub/test", nil)
	postRequest, _ := http.NewRequest("POST", "/pubsub/test", bytes.NewBuffer([]byte(testData)))

	requestHandler := http.HandlerFunc(pubSubReqHandler.ServeHttp)

	numberOfGetRequest := 1000000

	requestRecorderChan := make(chan *httptest.ResponseRecorder)
	comChan := make(chan struct{}, numberOfGetRequest)

	log.Println("TestServeHttpMulti", "sending GET reqs", numberOfGetRequest)
	for i := 0; i < numberOfGetRequest; i++ {
		go sendRequest(requestHandler, getRequest, requestRecorderChan, comChan)
	}

	log.Println("TestServeHttpMulti", "waiting for GET reqs to halt")
	for i := 0; i < numberOfGetRequest; i++ {
		select {
		case <-comChan:
		}
	}

	log.Println("TestServeHttpMulti", "sending POST req")
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
	requestRecorderPost := httptest.NewRecorder()
	requestHandler.ServeHTTP(requestRecorderPost, postRequest)
	return requestRecorderPost
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

func sendRequest(reqHandler http.HandlerFunc, req *http.Request, reqRecChan chan *httptest.ResponseRecorder, comChan chan struct{}) {
	requestRecord := httptest.NewRecorder()

	comChan <- struct{}{}

	reqHandler.ServeHTTP(requestRecord, req)

	reqRecChan <- requestRecord
}
