package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/hibare/GoCommon/v2/pkg/errors"
	commonHttp "github.com/hibare/GoCommon/v2/pkg/http"
	"github.com/hibare/GoCommon/v2/pkg/slice"
)

const AuthHeaderName = "Authorization"

func TokenAuth(next http.Handler, tokens []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(AuthHeaderName)

		if apiKey == "" || !slice.SliceContains(apiKey, tokens) {
			commonHttp.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &ResponseRecorder{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		defer func() {
			statusCode := recorder.StatusCode
			path := r.URL.Path
			method := r.Method
			duration := time.Since(start).String()

			switch {
			case statusCode < 400:
				slog.Info("request", "method", method, "path", path, "statusCode", statusCode, "duration", duration)
			case statusCode < 500:
				slog.Warn("request", "method", method, "path", path, "statusCode", statusCode, "duration", duration)
			default:
				slog.Error("request", "method", method, "path", path, "statusCode", statusCode, "duration", duration)
			}

		}()
		// Call the next handler
		next.ServeHTTP(recorder, r)
	})
}

// BasicSecurity adds basic security middleware.
func BasicSecurity(next http.Handler, sizeBytes int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Security headers
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Request size limit
		r.Body = http.MaxBytesReader(w, r.Body, sizeBytes)

		next.ServeHTTP(w, r)
	})
}
