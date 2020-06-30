package service

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/vladikan/url-shortener/logger"
)

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var fields []interface{}

		// Enrich log message
		ctx := r.Context()
		if reqID := middleware.GetReqID(ctx); reqID != "" {
			fields = append(fields, "request-id", reqID)
		}
		fields = append(fields, "proto", r.Proto)
		fields = append(fields, "uri", r.RequestURI)
		fields = append(fields, "verb", r.Method)

		logger.Debugw("Request started", fields...)
		next.ServeHTTP(w, r)
		logger.Debugw("Request completed", fields...)
	})
}
