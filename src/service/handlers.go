package service

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

// GetURI will resolve stored address by the code
func GetURI(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	w.Write([]byte(fmt.Sprintf("Your code is: %s", code)))
}

// PutURI will store new address and return it's code
func PutURI(w http.ResponseWriter, r *http.Request) {
	addr := chi.URLParam(r, "addr")
	w.Write([]byte(fmt.Sprintf("Your addr is: %s", addr)))
}
