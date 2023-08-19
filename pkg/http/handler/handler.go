package handler

import (
	"net/http"

	commonHttp "github.com/hibare/GoCommon/v2/pkg/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	commonHttp.WriteJsonResponse(w, http.StatusOK, map[string]bool{"ok": true})
}
