package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/vladikan/url-shortener/config"
	"github.com/vladikan/url-shortener/logger"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"golang.org/x/crypto/acme/autocert"
)

var srv *http.Server

// Start will setup http service
func Start(st *config.HostSettings) error {
	srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", st.Port),
		Handler: newRouter(),
	}

	var err error
	if st.Ssl {
		m := &autocert.Manager{
			Cache:      autocert.DirCache("autocert"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(st.Whitelist...),
		}

		srv.TLSConfig = m.TLSConfig()
		err = srv.ListenAndServeTLS("", "")
	} else {
		err = srv.ListenAndServe()
	}

	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to start server at %d, %s", st.Port, err))
	}

	logger.Info(fmt.Sprintf("Server started at %d", st.Port))
	return err
}

// Shutdown will try to stop server gracefully
func Shutdown() {
	logger.Info("Shutdown server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func newRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(5 * time.Second))

	r.Get("/{code}", GetURI)
	r.Put("/{addr}", PutURI)

	return r
}
