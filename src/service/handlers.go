package service

import (
	"fmt"
	"net/http"

	"github.com/vladikan/url-shortener/db"

	"github.com/go-chi/chi"
)

// GetURI will resolve stored address by the code
func GetURI(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	idx, _ := Decode(code)
	addr := db.Read(idx)

	w.Write([]byte(fmt.Sprintf("Your addr is: %s", addr)))
}

// PutURI will store new address and return it's code
func PutURI(w http.ResponseWriter, r *http.Request) {
	addr := chi.URLParam(r, "addr")
	idx, _ := db.Write(addr)
	code := Encode(idx)

	w.Write([]byte(fmt.Sprintf("Your code is: %s", code)))
}
