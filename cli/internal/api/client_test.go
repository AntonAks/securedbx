package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/securedbx/securedbx/cli/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) (*api.Client, *httptest.Server) {
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return api.New(srv.URL, "", "test-api-key"), srv
}

func TestUploadInit(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/upload/init", r.URL.Path)
		assert.Equal(t, "test-api-key", r.Header.Get("X-Cli-Api-Key"))

		_ = json.NewEncoder(w).Encode(api.UploadInitResponse{
			FileID:    "abc123",
			UploadURL: "https://s3.example.com/upload",
			ExpiresAt: 9999999999,
		})
	})

	resp, err := client.UploadInit(t.Context(), api.UploadInitRequest{
		ContentType: "file",
		FileSize:    1024,
		TTL:         "1h",
		AccessMode:  "one_time",
	})
	require.NoError(t, err)
	assert.Equal(t, "abc123", resp.FileID)
	assert.Equal(t, "https://s3.example.com/upload", resp.UploadURL)
}

func TestDownload(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/files/abc123/download", r.URL.Path)

		_ = json.NewEncoder(w).Encode(api.DownloadResponse{
			ContentType: "file",
			DownloadURL: "https://s3.example.com/file",
			FileSize:    1024,
			AccessMode:  "one_time",
		})
	})

	resp, err := client.Download(t.Context(), "abc123")
	require.NoError(t, err)
	assert.Equal(t, "file", resp.ContentType)
	assert.Equal(t, "https://s3.example.com/file", resp.DownloadURL)
}

func TestHTTPError(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "File not found"})
	})

	_, err := client.Download(t.Context(), "missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}

func TestPINUpload(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/pin/upload", r.URL.Path)
		_ = json.NewEncoder(w).Encode(api.PINUploadResponse{
			FileID:    "123456",
			Salt:      "a3f8c2d1e4b5a6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1",
			UploadURL: "https://s3.example.com/pin",
			ExpiresAt: 9999999999,
		})
	})

	resp, err := client.PINUpload(t.Context(), api.PINUploadRequest{
		ContentType: "file",
		FileSize:    512,
		PIN:         "1234",
		TTL:         "1h",
	})
	require.NoError(t, err)
	assert.Equal(t, "123456", resp.FileID)
	assert.Len(t, resp.Salt, 64)
}
