package service

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vladikan/url-shortener/config"
	"github.com/vladikan/url-shortener/db"
	"github.com/vladikan/url-shortener/logger"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"golang.org/x/crypto/acme/autocert"
)

var srv *http.Server
var srvDb db.Database

// Start will setup http service
func Start(st *config.HostSettings) {
	// Hook for system signal
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		logger.Warn("System interrupt signal")
		cancel()
	}()

	// Call shutdown service and db explicitly
	go func() {
		<-ctx.Done()
		shutdown()
	}()

	// Configure and start db
	srvDb = db.Open()

	// Configure and start service
	srv = &http.Server{Addr: st.Addr, Handler: newRouter()}
	logger.Infof("Server starting at %s", srv.Addr)

	var err error
	if st.Ssl {
		m := &autocert.Manager{
			Cache:      autocert.DirCache("autocert"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(st.Whitelist...),
		}

		srv.TLSConfig = m.TLSConfig()
		logger.Infof("Setting up SSL for the whitelist")

		// serve HTTP, which will redirect automatically to HTTPS
		go http.ListenAndServe(":http", m.HTTPHandler(nil))
		err = srv.ListenAndServeTLS("", "")
	} else {
		err = srv.ListenAndServe()
	}

	logger.Warnf("Server was terminated or failed to start, %s", err)
}

func shutdown() {
	logger.Info("Shutdown server")

	// Close service
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Warnf("Server shutdown error, %s", err)
	}

	// Close db
	db.Close(srvDb)
}

func newRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(loggerMiddleware)
	r.Use(middleware.Timeout(5 * time.Second))

	r.Get("/", InfoURI)
	r.Get("/{code}", GetURI)
	r.Put("/", PutURI)

	return r
}
