package service

import (
	"fmt"
	"net/http"

	"github.com/vladikan/url-shortener/config"
	"github.com/vladikan/url-shortener/db"
	"github.com/vladikan/url-shortener/logger"

	"github.com/go-chi/chi"
)

// GetURI will resolve stored address by the code
func GetURI(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if len(code) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missed request uri code, domain.com/[code]"))
		return
	}

	cfg := config.Service()
	idx, err := Decode(code, cfg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Validation error, %s", err)))
		return
	}

	addr := db.Read(idx)
	if len(addr) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("No record for this code `%s`", code)))
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

	idx, err := db.Write(addr)
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
