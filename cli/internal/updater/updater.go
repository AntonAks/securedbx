package updater

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	releasesURL = "https://api.github.com/repos/AntonAks/sdbx/releases/latest"
	cacheMaxAge = 24 * time.Hour
)

type cache struct {
	CheckedAt     time.Time `json:"checked_at"`
	LatestVersion string    `json:"latest_version"`
}

// Check returns the latest version if newer than currentVersion, otherwise "".
// Uses a 24h cache — refreshes in background when stale so there's no latency.
func Check(currentVersion string) string {
	if currentVersion == "dev" {
		return ""
	}

	c, err := readCache()
	stale := err != nil || time.Since(c.CheckedAt) >= cacheMaxAge

	if stale {
		go func() {
			if v, err := fetchLatest(); err == nil {
				_ = writeCache(cache{CheckedAt: time.Now(), LatestVersion: v})
			}
		}()
	}

	if err == nil && c.LatestVersion != "" && c.LatestVersion != currentVersion {
		return c.LatestVersion
	}
	return ""
}

func fetchLatest() (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(releasesURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var payload struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	return strings.TrimPrefix(payload.TagName, "v"), nil
}

func cachePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "securedbx", "update_check.json"), nil
}

func readCache() (*cache, error) {
	path, err := cachePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c cache
	return &c, json.Unmarshal(data, &c)
}

func writeCache(c cache) error {
	path, err := cachePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
