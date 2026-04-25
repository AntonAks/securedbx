package crypto

import (
	"archive/zip"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	KeySize          = 32 // 256 bits
	NonceSize        = 12 // 96 bits
	pbkdf2Iterations = 100000
)

// Encrypt generates a random 256-bit key and encrypts plaintext with AES-256-GCM.
// Returns: nonce||ciphertext||tag, raw key bytes.
func Encrypt(plaintext []byte) (ciphertext []byte, key []byte, err error) {
	key = make([]byte, KeySize)
	if _, err = io.ReadFull(rand.Reader, key); err != nil {
		return nil, nil, fmt.Errorf("generate key: %w", err)
	}
	ciphertext, err = encryptWithKey(plaintext, key)
	return
}

// Decrypt decrypts nonce||ciphertext||tag with the given key.
func Decrypt(data, key []byte) ([]byte, error) {
	return decryptWithKey(data, key)
}

// EncryptWithKey encrypts plaintext with an existing key.
func EncryptWithKey(plaintext, key []byte) ([]byte, error) {
	return encryptWithKey(plaintext, key)
}

// DecryptWithKey decrypts with an existing key.
func DecryptWithKey(data, key []byte) ([]byte, error) {
	return decryptWithKey(data, key)
}

func encryptWithKey(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}
	nonce := make([]byte, NonceSize)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decryptWithKey(data, key []byte) ([]byte, error) {
	if len(data) < NonceSize {
		return nil, fmt.Errorf("ciphertext too short: %d bytes", len(data))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}
	nonce, ct := data[:NonceSize], data[NonceSize:]
	pt, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, err
	}
	if pt == nil {
		pt = []byte{}
	}
	return pt, nil
}

// KeyToBase64URL encodes a key as base64url without padding (for URL fragment).
func KeyToBase64URL(key []byte) string {
	return base64.RawURLEncoding.EncodeToString(key)
}

// Base64URLToKey decodes a base64url key.
func Base64URLToKey(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

// DeriveKeyFromPIN derives a 256-bit key from a PIN and hex-decoded salt (server-provided as hex string).
func DeriveKeyFromPIN(pin string, salt []byte) []byte {
	return pbkdf2.Key([]byte(pin), salt, pbkdf2Iterations, KeySize, sha256.New)
}

// DeriveKeyFromPassword derives a 256-bit wrap key from a password and base64-decoded salt.
func DeriveKeyFromPassword(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, pbkdf2Iterations, KeySize, sha256.New)
}

// EncryptKey wraps dataKey with wrapKey using AES-256-GCM.
func EncryptKey(dataKey, wrapKey []byte) ([]byte, error) {
	return encryptWithKey(dataKey, wrapKey)
}

// DecryptKey unwraps an encrypted data key.
func DecryptKey(encryptedKey, wrapKey []byte) ([]byte, error) {
	return decryptWithKey(encryptedKey, wrapKey)
}

// GenerateSalt generates a random 16-byte salt for password mode.
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, salt)
	return salt, err
}

// FileEntry is a named file for multi-file bundling.
type FileEntry struct {
	Name string
	Data []byte
}

// BundleFiles creates a ZIP archive in memory from multiple files.
func BundleFiles(entries []FileEntry) ([]byte, error) {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	for _, e := range entries {
		f, err := w.Create(e.Name)
		if err != nil {
			return nil, fmt.Errorf("create zip entry %q: %w", e.Name, err)
		}
		if _, err = f.Write(e.Data); err != nil {
			return nil, fmt.Errorf("write zip entry %q: %w", e.Name, err)
		}
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("close zip: %w", err)
	}
	return buf.Bytes(), nil
}

// UnbundleFiles extracts all files from a ZIP archive.
func UnbundleFiles(data []byte) ([]FileEntry, error) {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}
	entries := make([]FileEntry, 0, len(r.File))
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("open zip entry %q: %w", f.Name, err)
		}
		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("read zip entry %q: %w", f.Name, err)
		}
		entries = append(entries, FileEntry{Name: f.Name, Data: content})
	}
	return entries, nil
}
