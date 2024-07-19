package tests

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modylegi/service/internal/api/transport/http/handlers"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	server := httptest.NewServer(handlers.Make(handlers.HealthHandler))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got: %d", resp.StatusCode)
	}

	expected := `{"message":"alive"}`
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	trimmedBody := bytes.TrimSpace(body)

	assert.Equal(t, expected, string(trimmedBody))

}
