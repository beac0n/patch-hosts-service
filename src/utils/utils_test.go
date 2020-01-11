package utils

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
)

var request = httptest.NewRequest("POST", "/test/path", bytes.NewBuffer([]byte("test_data")))

func TestLogErrorFalse(test *testing.T) {
	expected := false
	actual := LogError(nil, request)
	assert(test, actual, expected)
}

func TestLogErrorTrue(test *testing.T) {
	expected := true
	actual := LogError(errors.New("test error"), request)
	assert(test, actual, expected)
}

func TestHttpErrorRequestEntityTooLargeTrue(test *testing.T) {
	expected := true
	actual := HttpErrorRequestEntityTooLarge(1, request, httptest.NewRecorder())
	assert(test, actual, expected)

}

func TestHttpErrorRequestEntityTooLargeFalse(test *testing.T) {
	expected := false
	actual := HttpErrorRequestEntityTooLarge(10000, request, httptest.NewRecorder())
	assert(test, actual, expected)
}

func TestLoadAndStore(test *testing.T) {
	key := "test_key"
	m := &sync.Map{}
	chanCreator := func() interface{} { return make(chan struct{}) }

	channel := LoadAndStore(m, key, chanCreator)

	expected := reflect.TypeOf(make(chan struct{}))
	actual := reflect.TypeOf(channel)
	assert(test, actual, expected)
}

func TestNotGetOrPostTrue(test *testing.T) {
	putRequest := httptest.NewRequest("PUT", "/test/path", nil)
	assert(test, NotGetOrPost(putRequest, httptest.NewRecorder()), true)
}

func TestNotGetOrPostFalse(test *testing.T) {
	assert(test, NotGetOrPost(request, httptest.NewRecorder()), false)
}

func assert(test *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		test.Errorf("return value was %v, expected %v", actual, expected)
	}
}
