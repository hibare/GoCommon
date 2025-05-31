// Package handler provides HTTP handlers for the application.
package handler

import (
	"net/http"

	commonHttp "github.com/hibare/GoCommon/v2/pkg/http"
)

// HealthCheck is a simple health check handler.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	commonHttp.WriteJsonResponse(w, http.StatusOK, map[string]bool{"ok": true})
}
