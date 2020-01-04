package pubsub

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var pubSubReqHandler = &RequestHandler{maxReqSize: 10}
var testData0 = "test"
var testData1 = "test2"

func TestServeHttpSingle(test *testing.T) {
	getRequest, _ := http.NewRequest("GET", "/foobar", nil)
	postRequest, _ := http.NewRequest("POST", "/foobar", bytes.NewBuffer([]byte(testData0)))

	requestHandler := http.HandlerFunc(pubSubReqHandler.ServeHttp)

	requestRecorderChan := make(chan *httptest.ResponseRecorder)

	go sendRequest(requestHandler, getRequest, requestRecorderChan)

	assertRequest("", sendRequestSync(requestHandler, postRequest), test)
	assertRequest(testData0, <-requestRecorderChan, test)
}

func TestServeHttpSingleParallel(test *testing.T) {
	getRequest0, _ := http.NewRequest("GET", "/foobar", nil)
	postRequest0, _ := http.NewRequest("POST", "/foobar", bytes.NewBuffer([]byte(testData0)))

	getRequest1, _ := http.NewRequest("GET", "/barfoo", nil)
	postRequest1, _ := http.NewRequest("POST", "/barfoo", bytes.NewBuffer([]byte(testData1)))

	requestHandler := http.HandlerFunc(pubSubReqHandler.ServeHttp)

	requestRecorderChan0 := make(chan *httptest.ResponseRecorder)
	requestRecorderChan1 := make(chan *httptest.ResponseRecorder)

	go sendRequest(requestHandler, getRequest0, requestRecorderChan0)
	go sendRequest(requestHandler, getRequest1, requestRecorderChan1)

	assertRequest("", sendRequestSync(requestHandler, postRequest0), test)
	assertRequest(testData0, <-requestRecorderChan0, test)

	assertRequest("", sendRequestSync(requestHandler, postRequest1), test)
	assertRequest(testData1, <-requestRecorderChan1, test)
}


func TestServeHttpMulti(test *testing.T) {
	getRequest, _ := http.NewRequest("GET", "/foobar?pubsub=true", nil)
	postRequest, _ := http.NewRequest("POST", "/foobar?pubsub=true", bytes.NewBuffer([]byte(testData0)))

	requestHandler := http.HandlerFunc(pubSubReqHandler.ServeHttp)

	numberOfGetRequest := 10000

	requestRecorderChan := make(chan *httptest.ResponseRecorder, numberOfGetRequest)

	for i := 0; i < numberOfGetRequest; i++ {
		go sendRequest(requestHandler, getRequest, requestRecorderChan)
	}

	requestRecorderPost := sendRequestSync(requestHandler, postRequest)

	assertRequest("", requestRecorderPost, test)

	for i := 0; i < numberOfGetRequest; i++ {
		assertRequest(testData0, <-requestRecorderChan, test)
	}
}

func sendRequestSync(requestHandler http.HandlerFunc, postRequest *http.Request) *httptest.ResponseRecorder {
	requestRecorderPost := httptest.NewRecorder()
	requestHandler.ServeHTTP(requestRecorderPost, postRequest)
	return requestRecorderPost
}

func assertRequest(expected string, requestRecorder *httptest.ResponseRecorder, test *testing.T) {
	if status := requestRecorder.Code; status != http.StatusOK {
		test.Errorf("response has wrong status code: got %v want %v ", status, http.StatusOK)
	}

	if expected == "" {
		return
	}

	if actual := requestRecorder.Body.String(); actual != expected {
		test.Errorf("response has unexpected body: got %v want %v", actual, expected)
	}
}

func sendRequest(reqHandler http.HandlerFunc, req *http.Request, reqRecChan chan *httptest.ResponseRecorder) {
	requestRecord := httptest.NewRecorder()

	reqHandler.ServeHTTP(requestRecord, req)

	reqRecChan <- requestRecord
}
