package pubsub

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServeHttpSingle(test *testing.T) {
	getRequest, _ := http.NewRequest("GET", "/foobar", nil)
	postRequest, _ := http.NewRequest("POST", "/foobar", bytes.NewBuffer([]byte("test")))

	pubSubReqHandler := &RequestHandler{maxReqSizeInMb: 10}

	requestHandler := http.HandlerFunc(pubSubReqHandler.ServeHttp)

	requestRecorderGet := sendRequest(requestHandler, getRequest, true)
	requestRecorderPost := sendRequest(requestHandler, postRequest, false)

	assertRequest("", requestRecorderPost, test)
	assertRequest("test", requestRecorderGet, test)

}

func TestServeHttpMulti(test *testing.T) {
	getRequest, _ := http.NewRequest("GET", "/foobar?pubsub=true", nil)
	postRequest, _ := http.NewRequest("POST", "/foobar?pubsub=true", bytes.NewBuffer([]byte("test")))

	pubSubReqHandler := &RequestHandler{maxReqSizeInMb: 10}

	requestHandler := http.HandlerFunc(pubSubReqHandler.ServeHttp)

	requestRecorderGet0 := sendRequest(requestHandler, getRequest, true)
	requestRecorderGet1 := sendRequest(requestHandler, getRequest, true)
	requestRecorderGet2 := sendRequest(requestHandler, getRequest, true)

	time.Sleep(10 * time.Millisecond)

	requestRecorderPost := sendRequest(requestHandler, postRequest, false)

	assertRequest("", requestRecorderPost, test)

	assertRequest("test", requestRecorderGet0, test)
	assertRequest("test", requestRecorderGet1, test)
	assertRequest("test", requestRecorderGet2, test)
}

func assertRequest(expected string, requestRecorder *httptest.ResponseRecorder, test *testing.T) {
	if status := requestRecorder.Code; status != http.StatusOK {
		test.Errorf("requestHandler returned wrong status code: got %v want %v ", status, http.StatusOK)
	}

	if expected == "" {
		return
	}

	actual := requestRecorder.Body.String()
	if actual != expected {
		test.Errorf("requestHandler returned unexpected body: got %v want %v", actual, expected)
	}
}

func sendRequest(requestHandler http.HandlerFunc, request *http.Request, isAsync bool) *httptest.ResponseRecorder {
	requestRecord := httptest.NewRecorder()

	if isAsync {
		go requestHandler.ServeHTTP(requestRecord, request)
	} else {
		requestHandler.ServeHTTP(requestRecord, request)
	}

	return requestRecord
}
