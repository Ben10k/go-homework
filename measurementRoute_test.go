package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestBuildUrlString(t *testing.T) {
	expected := "https://localhost"
	actual := buildUrlString("https", "localhost")

	if actual != expected {
		t.Errorf("Unexpected url: actual [%s], expected [%s]", actual, expected)
	}
}
func TestBuildUrlString_space(t *testing.T) {
	expected := "https://local%20host"
	actual := buildUrlString("https", "local host")

	if actual != expected {
		t.Errorf("Unexpected url: actual [%s], expected [%s]", actual, expected)
	}
}

func TestOpenPage_success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	err := openPage(server.URL)

	if err != nil {
		t.Error("WebPage was not opened")
	}
}
func TestOpenPage_error(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))

	err := openPage(mockServer.URL)

	if err == nil {
		t.Error("openPage should return an error")
	}
}

func TestWork_noSamples(t *testing.T) {
	res, err := Work("cvhjdfhdf", 0, 1*time.Millisecond)
	if err != nil {
		t.Error("0 samples should not return an error")
	}
	if err == nil && len(res.Measurements) != 0 {
		t.Error("0 samples should return 0 measurements")
	}
}

func TestWork_oneSample(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	res, err := Work(server.URL, 1, 1*time.Millisecond)
	if err != nil {
		t.Error("correct server should return correct measurement")
	}
	if err == nil && len(res.Measurements) != 1 {
		t.Error("1 sample should return 1 measurement")
	}
}
func TestWork_oneSampleTimeOut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(15 * time.Second)
		w.WriteHeader(200)
	}))
	_, err := Work(server.URL, 1, 1*time.Millisecond)
	if err == nil || err.Error() != "timed out waiting for function to finish" {
		t.Error("Execution should have timed out")
	}
}

func TestWork_invalidUrl(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	_, err := Work(server.URL, 1, 1*time.Millisecond)
	if err == nil {
		t.Error("invalid server url should return an error")
	}

}
func TestHandler_invalidSamples(t *testing.T) {
	req, err := http.NewRequest("GET", "/measure", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(measureHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("measureHandler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"Message":"Invalid sample count"}`
	if rr.Body.String() != expected {
		t.Errorf("measureHandler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHandler_noHost(t *testing.T) {
	req, err := http.NewRequest("GET", "/measure?samples=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(measureHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("measureHandler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"Message":"Both [host] and [protocol] parameters must be present","Host":"","Protocol":""}`
	if rr.Body.String() != expected {
		t.Errorf("measureHandler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
func TestHandler_invalidHost(t *testing.T) {
	req, err := http.NewRequest("GET", "/measure?samples=1&protocol=qwerty&host=-abcde-", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(measureHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("measureHandler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"Message":"unable to open url [qwerty://-abcde-]"}`
	if rr.Body.String() != expected {
		t.Errorf("measureHandler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHandler_validHost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	parse, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal("this should not happen")
	}
	log.Println(server.URL)
	req, err := http.NewRequest("GET", "/measure?samples=1&protocol="+parse.Scheme+"&host="+parse.Host, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(measureHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("measureHandler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	result := &MeasurementResult{}
	err = json.Unmarshal(rr.Body.Bytes(), result)
	if err != nil {
		t.Fatal(err)
	}
	if result.Host != parse.Host || result.Protocol != parse.Scheme {
		t.Errorf("measureHandler returned unexpected body: got %v",
			rr.Body.String())
	}
}
