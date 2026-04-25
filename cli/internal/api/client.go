package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client is an HTTP client for the sdbx API.
type Client struct {
	baseURL    string
	cfSecret   string
	apiKey     string
	httpClient *http.Client
}

// New creates a new API client.
func New(baseURL, cfSecret, apiKey string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		cfSecret:   cfSecret,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *Client) do(ctx context.Context, method, path string, body any) (*http.Response, error) {
	var buf *bytes.Buffer
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		buf = bytes.NewBuffer(data)
	} else {
		buf = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.cfSecret != "" {
		req.Header.Set("X-Origin-Verify", c.cfSecret)
	}
	if c.apiKey != "" {
		req.Header.Set("X-Cli-Api-Key", c.apiKey)
	}
	return c.httpClient.Do(req)
}

func decode[T any](resp *http.Response) (T, error) {
	var result T
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errBody struct {
			Error string `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		return result, fmt.Errorf("HTTP %d: %s", resp.StatusCode, errBody.Error)
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

// --- Request / Response types ---

type UploadInitRequest struct {
	ContentType   string `json:"content_type"`
	FileSize      int64  `json:"file_size,omitempty"`
	EncryptedText string `json:"encrypted_text,omitempty"`
	TextSize      int    `json:"text_size,omitempty"` // original plaintext byte count for text secrets
	TTL           string `json:"ttl"`
	AccessMode    string `json:"access_mode"`
	Salt          string `json:"salt,omitempty"`
	EncryptedKey  string `json:"encrypted_key,omitempty"`
}

type UploadInitResponse struct {
	FileID    string `json:"file_id"`
	UploadURL string `json:"upload_url,omitempty"`
	ExpiresAt int64  `json:"expires_at"`
}

type DownloadResponse struct {
	ContentType   string `json:"content_type"`
	DownloadURL   string `json:"download_url,omitempty"`
	EncryptedText string `json:"encrypted_text,omitempty"`
	FileSize      int64  `json:"file_size"`
	AccessMode    string `json:"access_mode"`
}

type MetadataResponse struct {
	FileID        string `json:"file_id"`
	ContentType   string `json:"content_type"`
	FileSize      int64  `json:"file_size"`
	Available     bool   `json:"available"`
	ExpiresAt     int64  `json:"expires_at"`
	AccessMode    string `json:"access_mode"`
	Salt          string `json:"salt,omitempty"`
	EncryptedKey  string `json:"encrypted_key,omitempty"`
	DownloadCount int    `json:"download_count,omitempty"`
}

type PINUploadRequest struct {
	ContentType string `json:"content_type"`
	FileSize    int64  `json:"file_size,omitempty"`
	PIN         string `json:"pin"`
	FileName    string `json:"file_name,omitempty"`
	TTL         any    `json:"ttl"`
	AccessMode  string `json:"access_mode,omitempty"`
}

type PINUploadResponse struct {
	FileID    string `json:"file_id"`
	Salt      string `json:"salt"`
	UploadURL string `json:"upload_url,omitempty"`
	ExpiresAt int64  `json:"expires_at"`
}

type PINInitiateRequest struct {
	FileID string `json:"file_id"`
}

type PINInitiateResponse struct {
	Message        string `json:"message"`
	SessionExpires int64  `json:"session_expires"`
	AttemptsLeft   int    `json:"attempts_left"`
}

type PINVerifyRequest struct {
	FileID string `json:"file_id"`
	PIN    string `json:"pin"`
}

type PINVerifyResponse struct {
	ContentType   string `json:"content_type"`
	DownloadURL   string `json:"download_url,omitempty"`
	EncryptedText string `json:"encrypted_text,omitempty"`
	Salt          string `json:"salt"`
	FileSize      int64  `json:"file_size"`
	FileName      string `json:"file_name,omitempty"`
}

type AuthInitResponse struct {
	Challenge  string `json:"challenge"`
	Difficulty int    `json:"difficulty"`
}

type AuthConfirmRequest struct {
	Challenge string `json:"challenge"`
	Nonce     int    `json:"nonce"`
}

type AuthConfirmResponse struct {
	APIKey    string `json:"api_key"`
	ExpiresAt int64  `json:"expires_at"`
}

// --- Endpoint methods ---

func (c *Client) UploadInit(ctx context.Context, req UploadInitRequest) (UploadInitResponse, error) {
	resp, err := c.do(ctx, http.MethodPost, "/upload/init", req)
	if err != nil {
		return UploadInitResponse{}, err
	}
	return decode[UploadInitResponse](resp)
}

func (c *Client) Download(ctx context.Context, fileID string) (DownloadResponse, error) {
	resp, err := c.do(ctx, http.MethodPost, "/files/"+fileID+"/download", map[string]any{})
	if err != nil {
		return DownloadResponse{}, err
	}
	return decode[DownloadResponse](resp)
}

func (c *Client) Confirm(ctx context.Context, fileID string) error {
	resp, err := c.do(ctx, http.MethodPost, "/files/"+fileID+"/confirm", map[string]any{})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("confirm HTTP %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) Metadata(ctx context.Context, fileID string) (MetadataResponse, error) {
	resp, err := c.do(ctx, http.MethodGet, "/files/"+fileID+"/metadata", nil)
	if err != nil {
		return MetadataResponse{}, err
	}
	return decode[MetadataResponse](resp)
}

func (c *Client) PINUpload(ctx context.Context, req PINUploadRequest) (PINUploadResponse, error) {
	resp, err := c.do(ctx, http.MethodPost, "/pin/upload", req)
	if err != nil {
		return PINUploadResponse{}, err
	}
	return decode[PINUploadResponse](resp)
}

func (c *Client) PINInitiate(ctx context.Context, req PINInitiateRequest) (PINInitiateResponse, error) {
	resp, err := c.do(ctx, http.MethodPost, "/pin/initiate", req)
	if err != nil {
		return PINInitiateResponse{}, err
	}
	return decode[PINInitiateResponse](resp)
}

func (c *Client) PINVerify(ctx context.Context, req PINVerifyRequest) (PINVerifyResponse, error) {
	resp, err := c.do(ctx, http.MethodPost, "/pin/verify", req)
	if err != nil {
		return PINVerifyResponse{}, err
	}
	return decode[PINVerifyResponse](resp)
}

func (c *Client) UploadToS3(ctx context.Context, url string, body io.Reader, size int64) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.ContentLength = size
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("S3 upload HTTP %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) DownloadFromS3(ctx context.Context, url string, onProgress func(read, total int64)) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("S3 download HTTP %d", resp.StatusCode)
	}
	total := resp.ContentLength
	var buf bytes.Buffer
	tmp := make([]byte, 32*1024)
	var read int64
	for {
		n, err := resp.Body.Read(tmp)
		if n > 0 {
			buf.Write(tmp[:n])
			read += int64(n)
			if onProgress != nil {
				onProgress(read, total)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (c *Client) AuthInit(ctx context.Context) (AuthInitResponse, error) {
	resp, err := c.do(ctx, http.MethodPost, "/auth/init", map[string]any{})
	if err != nil {
		return AuthInitResponse{}, err
	}
	return decode[AuthInitResponse](resp)
}

func (c *Client) AuthConfirm(ctx context.Context, req AuthConfirmRequest) (AuthConfirmResponse, error) {
	resp, err := c.do(ctx, http.MethodPost, "/auth/confirm", req)
	if err != nil {
		return AuthConfirmResponse{}, err
	}
	return decode[AuthConfirmResponse](resp)
}
