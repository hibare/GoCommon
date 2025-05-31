package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealhCheckHandler(t *testing.T) {
	testCases := []struct {
		Name string
		URL  string
	}{
		{
			Name: "URL without trailing slash",
			URL:  "/ping",
		}, {
			Name: "URL with trailing slash",
			URL:  "/ping/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, err := http.NewRequestWithContext(t.Context(), http.MethodGet, tc.URL, nil)

			require.NoError(t, err)

			HealthCheck(w, r)

			require.Equal(t, http.StatusOK, w.Code)

			expectedBody := map[string]bool{"ok": true}
			responseBody := map[string]bool{}

			err = json.NewDecoder(w.Body).Decode(&responseBody)
			require.NoError(t, err)
			require.Equal(t, expectedBody, responseBody)
		})
	}
}
