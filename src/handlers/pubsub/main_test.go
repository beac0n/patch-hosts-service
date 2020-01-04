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

	pubSubReqHandler := &RequestHandler{maxReqSize: 10}

	requestHandler := http.HandlerFunc(pubSubReqHandler.ServeHttp)

	requestRecorderChan := make(chan *httptest.ResponseRecorder)

	go sendRequest(requestHandler, getRequest, requestRecorderChan)

	requestRecorderPost := httptest.NewRecorder()
	requestHandler.ServeHTTP(requestRecorderPost, postRequest)

	assertRequest("", requestRecorderPost, test)

	requestRecorderGet := <-requestRecorderChan
	assertRequest("test", requestRecorderGet, test)

}

func TestServeHttpMulti(test *testing.T) {
	data := `test`

	getRequest, _ := http.NewRequest("GET", "/foobar?pubsub=true", nil)
	postRequest, _ := http.NewRequest("POST", "/foobar?pubsub=true", bytes.NewBuffer([]byte(data)))

	pubSubReqHandler := &RequestHandler{maxReqSize: 10}

	requestHandler := http.HandlerFunc(pubSubReqHandler.ServeHttp)

	numberOfGetRequest := 10000

	requestRecorderChan := make(chan *httptest.ResponseRecorder)


	for i := 0; i < numberOfGetRequest; i++ {
		go sendRequest(requestHandler, getRequest, requestRecorderChan)
	}


	time.Sleep(100 * time.Millisecond)

	requestRecorderPost := httptest.NewRecorder()
	requestHandler.ServeHTTP(requestRecorderPost, postRequest)

	assertRequest("", requestRecorderPost, test)

	for i := 0; i < numberOfGetRequest; i++ {
		assertRequest(data, <-requestRecorderChan, test)
	}
}

func assertRequest(expected string, requestRecorder *httptest.ResponseRecorder, test *testing.T) {
	requestRecorder.Flush()
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

func sendRequest(reqHandler http.HandlerFunc, req *http.Request, reqRecChan chan *httptest.ResponseRecorder) {
	requestRecord := httptest.NewRecorder()

	reqHandler.ServeHTTP(requestRecord, req)

	reqRecChan <- requestRecord
}
