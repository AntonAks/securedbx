package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/securedbx/securedbx/cli/internal/api"
	"github.com/securedbx/securedbx/cli/internal/config"
)

// SolvePoW finds a nonce such that sha256(challenge+nonce) has `difficulty` leading zero hex chars.
func SolvePoW(challenge string, difficulty int) int {
	prefix := strings.Repeat("0", difficulty)
	for nonce := 0; ; nonce++ {
		h := sha256.Sum256([]byte(challenge + strconv.Itoa(nonce)))
		if strings.HasPrefix(hex.EncodeToString(h[:]), prefix) {
			return nonce
		}
	}
}

// VerifyPoW checks a PoW solution.
func VerifyPoW(challenge string, nonce, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	h := sha256.Sum256([]byte(challenge + strconv.Itoa(nonce)))
	return strings.HasPrefix(hex.EncodeToString(h[:]), prefix)
}

// EnsureValidKey checks if the stored API key is valid; if not, runs the PoW auth flow.
// Returns the valid API key.
func EnsureValidKey(ctx context.Context, apiClient *api.Client) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", fmt.Errorf("load config: %w", err)
	}
	if !cfg.IsExpired() {
		return cfg.APIKey, nil
	}
	return obtainNewKey(ctx, apiClient)
}

func obtainNewKey(ctx context.Context, apiClient *api.Client) (string, error) {
	initResp, err := apiClient.AuthInit(ctx)
	if err != nil {
		return "", fmt.Errorf("auth init: %w", err)
	}

	nonce := SolvePoW(initResp.Challenge, initResp.Difficulty)

	confirmResp, err := apiClient.AuthConfirm(ctx, api.AuthConfirmRequest{
		Challenge: initResp.Challenge,
		Nonce:     nonce,
	})
	if err != nil {
		return "", fmt.Errorf("auth confirm: %w", err)
	}

	cfg := &config.Config{
		APIKey:    confirmResp.APIKey,
		ExpiresAt: time.Unix(confirmResp.ExpiresAt, 0),
	}
	if err = config.Save(cfg); err != nil {
		return "", fmt.Errorf("save config: %w", err)
	}
	return cfg.APIKey, nil
}
