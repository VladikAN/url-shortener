package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/vladikan/url-shortener/config"
)

var r *chi.Mux

type dbMock struct {
	store string
	idx   uint64
	err   error
}

func TestMain(m *testing.M) {
	config.Init("../")
	srvDb = &dbMock{}

	r = chi.NewRouter()
	r.Get("/", InfoURI)
	r.Get("/{code}", GetURI)
	r.Put("/", PutURI)

	code := m.Run()
	os.Exit(code)
}

func TestGetServiceInfo(t *testing.T) {
	rq, _ := http.NewRequest("GET", "/", nil)
	rsp := httptest.NewRecorder()
	r.ServeHTTP(rsp, rq)

	if status := rsp.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got '%v' want '%v'", status, http.StatusOK)
	}
}

func TestGetWithInvalidChar(t *testing.T) {
	rq, _ := http.NewRequest("GET", "/test+", nil)
	rsp := httptest.NewRecorder()
	r.ServeHTTP(rsp, rq)

	if status := rsp.Code; status != http.StatusBadRequest {
		t.Errorf("Wrong status code: got '%v' want '%v'", status, http.StatusBadRequest)
	}

	expected := "Validation error, invalid character: +"
	if body := rsp.Body.String(); body != expected {
		t.Errorf("Wrong body: got '%s' want '%s'", body, expected)
	}
}

func TestGetUnknownCode(t *testing.T) {
	rq, _ := http.NewRequest("GET", "/unknown", nil)
	rsp := httptest.NewRecorder()
	r.ServeHTTP(rsp, rq)

	if status := rsp.Code; status != http.StatusBadRequest {
		t.Errorf("Wrong status code: got '%v' want '%v'", status, http.StatusBadRequest)
	}

	expected := "No record for this code `unknown`"
	if body := rsp.Body.String(); body != expected {
		t.Errorf("Wrong body: got '%s' want '%s'", body, expected)
	}
}

func TestGetKnownCode(t *testing.T) {
	srvDb = &dbMock{store: "http://example.com/"}
	rq, _ := http.NewRequest("GET", "/known", nil)
	rsp := httptest.NewRecorder()
	r.ServeHTTP(rsp, rq)

	if status := rsp.Code; status != http.StatusMovedPermanently {
		t.Errorf("Wrong status code: got '%v' want '%v'", status, http.StatusMovedPermanently)
	}

	expected := "http://example.com/"
	lct := rsp.HeaderMap["Location"][0]
	if lct != expected {
		t.Errorf("Wrong location: got '%s' want '%s'", lct, expected)
	}
}

func TestPutWithEmptyAddr(t *testing.T) {
	rq, _ := http.NewRequest("PUT", "/", nil)
	rsp := httptest.NewRecorder()
	r.ServeHTTP(rsp, rq)

	if status := rsp.Code; status != http.StatusBadRequest {
		t.Errorf("Wrong status code: got '%v' want '%v'", status, http.StatusBadRequest)
	}

	expected := "Missed request `addr` parameter"
	if body := rsp.Body.String(); body != expected {
		t.Errorf("Wrong body: got '%s' want '%s'", body, expected)
	}
}

func TestPutGivesDbError(t *testing.T) {
	srvDb = &dbMock{err: errors.New("test error")}
	rq, _ := http.NewRequest("PUT", "/", strings.NewReader("addr=http://example.com"))
	rq.Header.Set("content-type", "application/x-www-form-urlencoded")

	rsp := httptest.NewRecorder()
	r.ServeHTTP(rsp, rq)

	if status := rsp.Code; status != http.StatusInternalServerError {
		t.Errorf("Wrong status code: got '%v' want '%v'", status, http.StatusInternalServerError)
	}
}

func TestPutGivesNewCode(t *testing.T) {
	srvDb = &dbMock{idx: 10241024}
	rq, _ := http.NewRequest("PUT", "/", strings.NewReader("addr=http://example.com"))
	rq.Header.Set("content-type", "application/x-www-form-urlencoded")

	rsp := httptest.NewRecorder()
	r.ServeHTTP(rsp, rq)

	if status := rsp.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got '%v' want '%v'", status, http.StatusOK)
	}

	expected := "Q8jY"
	if body := rsp.Body.String(); body != expected {
		t.Errorf("Wrong body: got '%s' want '%s'", body, expected)
	}
}

func (db dbMock) Read(key uint64) string {
	return db.store
}

func (db dbMock) Write(value string) (uint64, error) {
	return db.idx, db.err
}
