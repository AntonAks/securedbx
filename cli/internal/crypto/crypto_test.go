package crypto_test

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/securedbx/securedbx/cli/internal/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	tests := []struct {
		name      string
		plaintext []byte
	}{
		{"empty", []byte{}},
		{"hello", []byte("hello world")},
		{"unicode", []byte("Привіт 🔒")},
		{"1kb", bytes.Repeat([]byte("a"), 1024)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct, key, err := crypto.Encrypt(tt.plaintext)
			require.NoError(t, err)
			assert.Len(t, key, 32, "key must be 32 bytes")
			assert.GreaterOrEqual(t, len(ct), 12+16, "ciphertext must be >= nonce+tag")

			pt, err := crypto.Decrypt(ct, key)
			require.NoError(t, err)
			assert.Equal(t, tt.plaintext, pt)
		})
	}
}

func TestDecryptWrongKey(t *testing.T) {
	ct, key, err := crypto.Encrypt([]byte("secret"))
	require.NoError(t, err)
	key[0] ^= 0xFF
	_, err = crypto.Decrypt(ct, key)
	assert.Error(t, err, "wrong key must fail")
}

func TestDecryptTruncated(t *testing.T) {
	_, err := crypto.Decrypt([]byte("tooshort"), make([]byte, 32))
	assert.Error(t, err)
}

func TestKeyEncoding(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef") // 32 bytes
	encoded := crypto.KeyToBase64URL(key)
	assert.NotContains(t, encoded, "+")
	assert.NotContains(t, encoded, "/")
	assert.NotContains(t, encoded, "=")

	decoded, err := crypto.Base64URLToKey(encoded)
	require.NoError(t, err)
	assert.Equal(t, key, decoded)
}

func TestWireFormatLayout(t *testing.T) {
	ct, key, err := crypto.Encrypt([]byte("test"))
	require.NoError(t, err)
	nonce := ct[:12]
	assert.Len(t, nonce, 12)
	assert.Equal(t, 12+4+16, len(ct))
	_ = key
}

// Task 4 tests
func TestDeriveKeyFromPIN(t *testing.T) {
	salt, _ := hex.DecodeString("a3f8c2d1e4b5a6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1")
	key := crypto.DeriveKeyFromPIN("1234", salt)
	assert.Len(t, key, 32)

	key2 := crypto.DeriveKeyFromPIN("1234", salt)
	assert.Equal(t, key, key2, "same inputs must produce same key")

	keyOther := crypto.DeriveKeyFromPIN("5678", salt)
	assert.NotEqual(t, key, keyOther, "different PIN must produce different key")
}

func TestDeriveKeyFromPassword(t *testing.T) {
	salt, _ := base64.StdEncoding.DecodeString("AAAAAAAAAAAAAAAAAAAAAA==") // 16 bytes
	key := crypto.DeriveKeyFromPassword("hunter2", salt)
	assert.Len(t, key, 32)

	key2 := crypto.DeriveKeyFromPassword("hunter2", salt)
	assert.Equal(t, key, key2)
}

func TestEncryptKeyRoundTrip(t *testing.T) {
	dataKey := make([]byte, 32)
	_, _ = rand.Read(dataKey)

	wrapKey := make([]byte, 32)
	_, _ = rand.Read(wrapKey)

	wrapped, err := crypto.EncryptKey(dataKey, wrapKey)
	require.NoError(t, err)

	unwrapped, err := crypto.DecryptKey(wrapped, wrapKey)
	require.NoError(t, err)
	assert.Equal(t, dataKey, unwrapped)
}

func TestEncryptKeyWrongWrapKey(t *testing.T) {
	dataKey := make([]byte, 32)
	_, _ = rand.Read(dataKey)
	wrapKey := make([]byte, 32)
	_, _ = rand.Read(wrapKey)

	wrapped, err := crypto.EncryptKey(dataKey, wrapKey)
	require.NoError(t, err)

	badKey := make([]byte, 32)
	_, _ = rand.Read(badKey)
	_, err = crypto.DecryptKey(wrapped, badKey)
	assert.Error(t, err)
}

// Task 5 tests
func TestBundleAndUnbundle(t *testing.T) {
	entries := []crypto.FileEntry{
		{Name: "hello.txt", Data: []byte("hello")},
		{Name: "world.bin", Data: []byte{0x01, 0x02, 0x03}},
		{Name: "субтитри.srt", Data: []byte("Ukrainian filename")},
	}

	bundle, err := crypto.BundleFiles(entries)
	require.NoError(t, err)
	assert.Greater(t, len(bundle), 0)

	extracted, err := crypto.UnbundleFiles(bundle)
	require.NoError(t, err)
	require.Len(t, extracted, 3)

	for i, e := range entries {
		assert.Equal(t, e.Name, extracted[i].Name)
		assert.Equal(t, e.Data, extracted[i].Data)
	}
}

func TestBundleEmpty(t *testing.T) {
	bundle, err := crypto.BundleFiles(nil)
	require.NoError(t, err)
	assert.Greater(t, len(bundle), 0, "empty ZIP is still valid")
}

func TestUnbundleInvalidData(t *testing.T) {
	_, err := crypto.UnbundleFiles([]byte("not a zip"))
	assert.Error(t, err)
}
