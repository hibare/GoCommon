package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testToken = "test-key"
)

func TestTokenAuthSuccess(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(AuthHeaderName, testToken)

	rr := httptest.NewRecorder()

	mw := TokenAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), []string{testToken})

	mw.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestTokenAuthNoKeyFailure(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mw := TokenAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), []string{testToken})

	mw.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestTokenAuthWrongKeyFailure(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(AuthHeaderName, "adfafs")

	rr := httptest.NewRecorder()

	mw := TokenAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), []string{testToken})

	mw.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestRequestLogger(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(log.Writer())
		}()

		// Create a mock server
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello, world!"))
		})

		// Create a test server with the RequestLogger middleware
		ts := httptest.NewServer(RequestLogger(handler))
		defer ts.Close()

		// Make a request to the test server to trigger logging
		client := &http.Client{}
		resp, err := client.Get(ts.URL)
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		logString := buf.String()
		assert.NotEmpty(t, logString)
		assert.Contains(t, logString, "INFO request method=GET path=/ statusCode=200 duration")

	})

	t.Run("Warning", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(log.Writer())
		}()

		// Create a mock server
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		})

		// Create a test server with the RequestLogger middleware
		ts := httptest.NewServer(RequestLogger(handler))
		defer ts.Close()

		// Make a request to the test server to trigger logging
		client := &http.Client{}
		resp, err := client.Get(ts.URL)
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		logString := buf.String()
		assert.NotEmpty(t, logString)
		assert.Contains(t, logString, "WARN request method=GET path=/ statusCode=400 duration=")
	})

	t.Run("Error", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(log.Writer())
		}()

		// Create a mock server
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})

		// Create a test server with the RequestLogger middleware
		ts := httptest.NewServer(RequestLogger(handler))
		defer ts.Close()

		// Make a request to the test server to trigger logging
		client := &http.Client{}
		resp, err := client.Get(ts.URL)
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		logString := buf.String()
		assert.NotEmpty(t, logString)
		assert.Contains(t, logString, "ERROR request method=GET path=/ statusCode=500 duration=")
	})

}

func TestBasicSecurity(t *testing.T) {
	handler := BasicSecurity(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, response.Code)
	}
}
