package config_test

import (
	"testing"
	"time"

	"github.com/securedbx/securedbx/cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfg := &config.Config{
		APIKey:    "sdbx_cli_abc123",
		ExpiresAt: time.Now().Add(24 * time.Hour).UTC().Truncate(time.Second),
	}
	require.NoError(t, config.Save(cfg))

	loaded, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, cfg.APIKey, loaded.APIKey)
	assert.Equal(t, cfg.ExpiresAt.Unix(), loaded.ExpiresAt.Unix())
}

func TestLoadMissing(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Empty(t, cfg.APIKey)
}

func TestIsExpired(t *testing.T) {
	past := &config.Config{APIKey: "key", ExpiresAt: time.Now().Add(-1 * time.Hour)}
	assert.True(t, past.IsExpired())

	future := &config.Config{APIKey: "key", ExpiresAt: time.Now().Add(1 * time.Hour)}
	assert.False(t, future.IsExpired())
}

func TestEnvOverride(t *testing.T) {
	t.Setenv("SECUREDBX_API_KEY", "env_key")
	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "env_key", cfg.APIKey)
}
