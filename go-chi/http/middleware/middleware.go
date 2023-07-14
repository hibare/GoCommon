package middleware

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hibare/GoCommon/errors"
	"github.com/hibare/GoCommon/slice"
	log "github.com/sirupsen/logrus"
)

const AuthHeaderName = "Authorization"

func TokenAuth(next http.Handler, tokens []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Client: [%s] %s", r.RemoteAddr, r.RequestURI)

		apiKey := r.Header.Get(AuthHeaderName)

		if apiKey == "" {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, errors.Error{
				Message: errors.ErrUnauthorized.Error(),
			})
			return
		}

		if slice.SliceContains(apiKey, tokens) {
			next.ServeHTTP(w, r)
			return
		}

		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, errors.Error{
			Message: errors.ErrUnauthorized.Error(),
		})
	})
}
