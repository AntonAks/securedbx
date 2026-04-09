package auth_test

import (
	"testing"

	"github.com/securedbx/securedbx/cli/internal/auth"
	"github.com/stretchr/testify/assert"
)

func TestSolvePoW(t *testing.T) {
	challenge := "testchallenge"
	difficulty := 2 // use 2 in tests for speed

	nonce := auth.SolvePoW(challenge, difficulty)
	assert.GreaterOrEqual(t, nonce, 0)

	// Verify solution
	assert.True(t, auth.VerifyPoW(challenge, nonce, difficulty))
}

func TestVerifyPoWFails(t *testing.T) {
	assert.False(t, auth.VerifyPoW("challenge", 0, 4))
}

func TestSolvePoWDifficulty4(t *testing.T) {
	// Difficulty 4 = 4 leading hex zeros — should solve in <2s
	nonce := auth.SolvePoW("abc123", 4)
	assert.True(t, auth.VerifyPoW("abc123", nonce, 4))
}
