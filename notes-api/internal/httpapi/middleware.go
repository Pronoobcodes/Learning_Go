package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

type contextKey string

const requestIDKey contextKey = "request_id"

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(data []byte) (int, error) {
	if r.statusCode == 0 {
		r.statusCode = http.StatusOK
	}
	return r.ResponseWriter.Write(data)
}

func (r *statusRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		// Generate 16 random bytes for this request.
		var id [16]byte

		if _, err := rand.Read(id[:]); err != nil {
			log.Printf("could not generate request ID: %v", err)
			http.Error(w, "could not generate request ID", http.StatusInternalServerError)
			return
		}

		requestID := hex.EncodeToString(id[:])

		// Send the ID back to the client.
		w.Header().Set("X-Request-ID", requestID)

		// Store it in the request context too.
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		start := time.Now()

		recorder := &statusRecorder{
			ResponseWriter: w,
		}

		next.ServeHTTP(recorder, r)

		status := recorder.statusCode

		if status == 0 {
			status = http.StatusOK
		}

		duration := time.Since(start)

		requestID := RequestIDFromContext(r.Context())

		log.Printf(
			"request_id=%s method=%s path=%s status=%d duration=%s remote_addr=%s",
			requestID,
			r.Method,
			r.URL.Path,
			status,
			duration,
			r.RemoteAddr,
		)
	})
}

func RequestIDFromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}

	return requestID
}
