package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hibare/GoCommon/errors"
	"github.com/hibare/GoCommon/slice"
)

const AuthHeaderName = "Authorization"

func TokenAuth(next http.Handler, tokens []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Client: [%s] %s", r.RemoteAddr, r.RequestURI)

		apiKey := r.Header.Get(AuthHeaderName)

		if apiKey == "" || !slice.SliceContains(apiKey, tokens) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)

			e := errors.Error{
				Code:    http.StatusUnauthorized,
				Message: errors.ErrUnauthorized.Error(),
			}

			jsonData, err := json.Marshal(e)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "something went wrong, failed to authenticate")
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(jsonData)
			return
		}
		next.ServeHTTP(w, r)
	})
}
