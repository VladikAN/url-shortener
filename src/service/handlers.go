package service

import (
	"fmt"
	"net/http"

	"github.com/vladikan/url-shortener/config"
	"github.com/vladikan/url-shortener/logger"

	"github.com/go-chi/chi"
)

// InfoURI will print service information
func InfoURI(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
Url shortener service.
GET: domain.com/ - to pring this output.
GET: domain.com/testcode - to resolve URI stored by code 'testcode'.
PUT: domain.com addr=https://example.com/ - to store specified address and get its code.`))
}

// GetURI will resolve stored address by the code
func GetURI(w http.ResponseWriter, r *http.Request) {
	cfg := config.Service()

	code := chi.URLParam(r, "code")
	idx, err := Decode(code, cfg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Validation error, %s", err)))
		return
	}

	addr := srvDb.Read(idx)
	if len(addr) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("No record for this code `%s`", code)))
		return
	}

	logger.Debugf("Incomming request for `%s` code resolved by `%s` and `%d` index", code, addr, idx)
	http.Redirect(w, r, addr, http.StatusMovedPermanently)
}

// PutURI will store new address and return it's code
func PutURI(w http.ResponseWriter, r *http.Request) {
	addr := r.FormValue("addr")
	if len(addr) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missed request `addr` parameter"))
		return
	}

	idx, err := srvDb.Write(addr)
	if err != nil {
		logger.Warnf("Error occured while performing write operation, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cfg := config.Service()
	code := Encode(idx, cfg)
	logger.Debugf("Put request for addr `%s` completed with `%s` code and `%d` index", addr, code, idx)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(code))
}
