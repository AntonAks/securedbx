package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Config holds CLI authentication state.
type Config struct {
	APIKey    string    `json:"api_key"`
	ExpiresAt time.Time `json:"expires_at"`
}

// IsExpired returns true if the API key is missing or past its expiry.
func (c *Config) IsExpired() bool {
	return c.APIKey == "" || time.Now().After(c.ExpiresAt)
}

// BaseURL returns the API base URL from env or default.
func BaseURL() string {
	if v := os.Getenv("SECUREDBX_BASE_URL"); v != "" {
		return v
	}
	return "https://api.securedbx.com"
}

// CFSecret returns the CloudFront origin secret from env.
func CFSecret() string {
	return os.Getenv("SECUREDBX_CF_SECRET")
}

// FrontendURL returns the frontend base URL from env or default.
func FrontendURL() string {
	if v := os.Getenv("SECUREDBX_FRONTEND_URL"); v != "" {
		return v
	}
	return "https://securedbx.com"
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "securedbx", "config.json"), nil
}

// Load reads the config file. If missing, returns empty config.
// If SECUREDBX_API_KEY env var is set, it overrides the file value.
func Load() (*Config, error) {
	cfg := &Config{}

	path, err := configPath()
	if err != nil {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// no file yet — check env override below
	} else if err != nil {
		return nil, err
	} else {
		if err = json.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	// Env var takes precedence over file (never expires from CLI's perspective)
	if key := os.Getenv("SECUREDBX_API_KEY"); key != "" {
		cfg.APIKey = key
		cfg.ExpiresAt = time.Now().Add(365 * 24 * time.Hour)
	}

	return cfg, nil
}

// Save writes the config to disk with 0600 permissions.
func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
