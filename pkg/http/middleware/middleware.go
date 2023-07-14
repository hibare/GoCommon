package middleware

import (
	"log"
	"net/http"

	"github.com/hibare/GoCommon/pkg/errors"
	commonHttp "github.com/hibare/GoCommon/pkg/http"
	"github.com/hibare/GoCommon/pkg/slice"
)

const AuthHeaderName = "Authorization"

func TokenAuth(next http.Handler, tokens []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Client: [%s] %s", r.RemoteAddr, r.RequestURI)

		apiKey := r.Header.Get(AuthHeaderName)

		if apiKey == "" || !slice.SliceContains(apiKey, tokens) {
			commonHttp.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
