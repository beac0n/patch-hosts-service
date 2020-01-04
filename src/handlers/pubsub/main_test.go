package pubsub

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var pubSubReqHandler = NewRequestHandler(10)
var testData = "test"

func TestServeHttpMulti(test *testing.T) {
	getRequest, _ := http.NewRequest("GET", "/foobar", nil)
	postRequest, _ := http.NewRequest("POST", "/foobar", bytes.NewBuffer([]byte(testData)))

	requestHandler := http.HandlerFunc(pubSubReqHandler.ServeHttp)

	numberOfGetRequest := 10000

	requestRecorderChan := make(chan *httptest.ResponseRecorder, numberOfGetRequest)

	for i := 0; i < numberOfGetRequest; i++ {
		go sendRequest(requestHandler, getRequest, requestRecorderChan)
	}

	requestRecorderPost := sendRequestSync(requestHandler, postRequest)

	assertRequest("", requestRecorderPost, test)

	for i := 0; i < numberOfGetRequest; i++ {
		assertRequest(testData, <-requestRecorderChan, test)
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
