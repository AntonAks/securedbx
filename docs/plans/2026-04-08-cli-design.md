# securedbx-cli — Design Spec

> Go CLI client for sdbx.cc · v1.0 · 2026-04-08

---

## 1. Context

SecureDBX (sdbx.cc) is a zero-knowledge file and text sharing service with end-to-end encryption and one-time access. All cryptographic operations run in the browser — the server never sees keys or plaintext.

This spec describes a CLI client for sdbx.cc implemented in Go. The CLI lives in the existing monorepo as a new `cli/` directory alongside `backend/` and `frontend/`. Distributed as a single binary via Homebrew, Scoop, and GitHub Releases.

---

## 2. Repository Structure

```
securedbx/
├── backend/                      # Lambda functions (Python)
├── frontend/                     # Vue 3 SPA
├── cli/                          # ← new directory
│   ├── cmd/securedbx/main.go
│   ├── internal/
│   │   ├── crypto/               # AES-256-GCM, PBKDF2, ZIP
│   │   │   ├── crypto.go
│   │   │   └── crypto_test.go
│   │   ├── api/                  # HTTP client to Lambda
│   │   │   ├── client.go
│   │   │   └── client_test.go
│   │   ├── commands/             # send, receive, status
│   │   │   ├── send.go
│   │   │   ├── receive.go
│   │   │   └── status.go
│   │   ├── auth/                 # API Key management
│   │   │   ├── auth.go
│   │   │   └── auth_test.go
│   │   └── config/config.go
│   ├── go.mod                    # module github.com/securedbx/securedbx/cli
│   ├── go.sum
│   ├── .goreleaser.yaml
│   └── Makefile
├── terraform/
├── scripts/
└── Makefile                      # add cli/* targets
```

### Go Module

Separate Go module inside the monorepo — standard practice for polyglot repos:

```
module github.com/securedbx/securedbx/cli

go 1.22
```

### Root Makefile additions

```makefile
cli-build:
	cd cli && go build ./cmd/securedbx

cli-test:
	cd cli && go test ./...

cli-lint:
	cd cli && golangci-lint run

cli-release:
	cd cli && goreleaser release --clean
```

---

## 3. Dependencies

### External (UI only)

| Package | Version | Purpose |
|---|---|---|
| `github.com/spf13/cobra` | v1.8+ | CLI framework |
| `github.com/charmbracelet/lipgloss` | v0.10+ | Terminal styling |
| `github.com/charmbracelet/bubbles` | v0.18+ | Progress bars, spinners |
| `golang.org/x/crypto` | latest | PBKDF2 |
| `github.com/stretchr/testify` | latest | Test assertions only |

### Stdlib (crypto is 100% stdlib)

`crypto/aes`, `crypto/cipher`, `crypto/rand`, `crypto/sha256`, `archive/zip`, `net/http`, `encoding/json`, `log/slog`, `os`, `path/filepath`

---

## 4. Crypto Layer (`internal/crypto`)

Go implementation mirrors browser `crypto.js` exactly. Files encrypted by the browser can be decrypted by the CLI and vice versa.

### Two modes

**URL mode** — uses `POST /upload/init` and `GET /download/{id}`

- Base: random 256-bit key embedded in URL `#fragment` as base64url (no padding)
- With password (`--password`): `PBKDF2(password, client_salt, 100k iterations)` → wrap key → double encryption; `salt` + `encrypted_key` sent in request body
- `access_mode`: `one_time` (default) or `multi` (with `--multi` flag)

**PIN mode** — uses `POST /pin/initiate` → `POST /pin/upload` → `POST /pin/verify`

- Salt provided by server
- `PBKDF2(pin, server_salt, 100k iterations)` → key
- Recipient knows File ID + PIN only (no URL with key)

### Wire format (identical for both modes)

```
encrypted blob = nonce (12 bytes) || ciphertext || auth tag (16 bytes)
```

### Vault (double encryption) detail

| Step | Operation |
|---|---|
| 1 | Generate random 256-bit data key |
| 2 | `PBKDF2(password, client_salt, 100k)` → wrap key |
| 3 | `AES-256-GCM(data_key, wrap_key)` → encrypted_key |
| 4 | `AES-256-GCM(plaintext, data_key)` → ciphertext |
| 5 | Upload: ciphertext + encrypted_key + salt |

Client generates salt (`crypto/rand`) — unlike PIN mode where server provides it.

### Multi-file

`archive/zip` (stdlib) — ZIP in memory before encryption. Original filenames preserved inside archive. Server-side name always `archive.zip`. Limit: 10 files, 500 MB total.

### Functions

```go
func Encrypt(plaintext []byte) (ciphertext, key []byte, err error)
func Decrypt(ciphertext, key []byte) (plaintext []byte, err error)
func DeriveKeyFromPassword(password string, salt []byte) []byte
func DeriveKeyFromPIN(pin string, salt []byte) []byte
func EncryptKey(dataKey, wrapKey []byte) ([]byte, error)
func DecryptKey(encryptedKey, wrapKey []byte) ([]byte, error)
```

---

## 5. API Client (`internal/api`)

### Authentication

| Mechanism | Browser | CLI |
|---|---|---|
| CloudFront Origin Header | Automatic | `X-CloudFront-Secret` from `SECUREDBX_CF_SECRET` |
| reCAPTCHA v3 | Automatic | `X-CLI-API-Key` instead |

### Endpoint mapping

| Command | Endpoint | Purpose |
|---|---|---|
| `send` (URL) | `POST /upload/init` | Upload encrypted file/text |
| `receive` phase 1 | `POST /files/{id}/download` | Reserve (10-min window for one-time) |
| `receive` phase 2 | `POST /files/{id}/confirm` | Confirm — marks downloaded (one-time only) |
| `receive --password` pre-step | `GET /files/{id}/metadata` | Fetch salt + encrypted_key |
| `send --pin` | `POST /pin/upload` | Upload with PIN (server generates salt, returns it) |
| `receive --pin` step 1 | `POST /pin/initiate` | Start 60-second session |
| `receive --pin` step 2 | `POST /pin/verify` | Verify PIN + get download URL |
| `status` | `GET /files/{id}/metadata` | File metadata |
| `auth` | `POST /auth/init` | Get PoW challenge |
| `auth` | `POST /auth/confirm` | Verify PoW, get API Key |

> Two-phase download for one-time mode: `POST /files/{id}/download` (reserve) → download from S3 → decrypt → `POST /files/{id}/confirm`.
> Confirm is called **only** after successful decryption. Multi-access (vault) mode skips confirm entirely.

---

## 6. CLI Commands

### `send`

```bash
securedbx send <file> [<file> ...]
securedbx send --text "secret text"
```

| Flag | Default | Description |
|---|---|---|
| `--text TEXT` | — | Send text (up to 1000 chars) |
| `--pin` | false | PIN mode (prompts for PIN) |
| `--pin-value N` | — | PIN explicitly 4-8 digits (for scripts) |
| `--password` | false | Password protection (prompts) |
| `--password-value P` | — | Password explicitly (for scripts) |
| `--multi` | false | Multi-access mode (for --password) |
| `--ttl DURATION` | `1h` | `5m / 1h / 12h / 24h / 3d / 7d` |
| `--json` | false | Output result as JSON |

### `receive`

```bash
securedbx receive "https://sdbx.cc/#key#filename"
securedbx receive --pin <file-id>
securedbx receive --password <file-id>
```

| Flag | Default | Description |
|---|---|---|
| `--pin FILE-ID` | — | PIN mode |
| `--pin-value N` | — | PIN (prompts if omitted) |
| `--password` | false | Password-protected URL mode |
| `--password-value P` | — | Password (prompts if omitted) |
| `--output PATH` | `.` | Directory or path to save |
| `--json` | false | Output metadata as JSON |

### `status`

```bash
securedbx status <FILE-ID>
```

Shows metadata: size, expiry, download status. Useful for scripting.

---

## 7. Auth (`internal/auth`)

### Concept

No accounts, no registration — consistent with sdbx.cc philosophy. Automatic anonymous API Key generation via Proof-of-Work.

### Flow

```
1. CLI → POST /auth/init
   ← { challenge: "a3f8...", difficulty: 4 }

2. CLI solves PoW (~0.5s):
   for nonce := 0; ; nonce++ {
       h := sha256(challenge + strconv.Itoa(nonce))
       if strings.HasPrefix(hex(h), "0000") { break }
   }

3. CLI → POST /auth/confirm { challenge, nonce }
   ← { api_key: "sdbx_cli_xxx...", expires_at }

4. CLI saves: ~/.config/securedbx/config.json
```

### Auto-renewal

CLI checks key expiry before every command. If expired or missing → silently runs auth flow. User sees no interruption.

### Key parameters

| Parameter | Value |
|---|---|
| TTL | 24 hours |
| Storage | `~/.config/securedbx/config.json` |
| Format | `sdbx_cli_<random_hex>` |
| Request header | `X-CLI-API-Key` |

### New Lambda functions required

| Lambda | Endpoint | Purpose |
|---|---|---|
| `auth_init` | `POST /auth/init` | Generate PoW challenge, store in DynamoDB (TTL 5 min) |
| `auth_confirm` | `POST /auth/confirm` | Verify PoW, generate API Key, store in DynamoDB (TTL 24h) |

---

## 8. Backend Changes

### New decorator: `require_cloudfront_and_auth`

Added to `backend/shared/security.py`. The **old** `require_cloudfront_and_recaptcha` decorator is not touched.

Logic:
- Always verify CloudFront origin
- If `X-CLI-API-Key` header present → validate key against DynamoDB
- Else → validate `recaptcha_token` via Google (existing flow)

All handlers replace `@require_cloudfront_and_recaptcha` with `@require_cloudfront_and_auth`:
- `upload_init`, `download`, `confirm_download`, `pin_initiate`, `pin_upload_init`, `pin_verify`

### Terraform

Two new Lambda blocks in `terraform/modules/api/main.tf`: `auth_init`, `auth_confirm`.
New API Gateway routes: `POST /auth/init`, `POST /auth/confirm`.

---

## 9. Configuration

| Variable | Purpose |
|---|---|
| `SECUREDBX_API_KEY` | Manual key override — skips PoW flow |
| `SECUREDBX_CF_SECRET` | CloudFront origin header value |
| `SECUREDBX_BASE_URL` | Default: `https://api.sdbx.cc` |
| `SECUREDBX_DEFAULT_TTL` | Default: `1h` |
| `SECUREDBX_OUTPUT_DIR` | Default: `.` |

---

## 10. Anti-spam Strategy

| Layer | Priority | Implementation |
|---|---|---|
| CLI API Key | MUST | Replaces reCAPTCHA. Rate-limited via API Gateway Usage Plans. |
| Rate Limiting | MUST | 10 uploads/min, 50/hour per key. |
| PoW (key issuance) | MUST | ~0.5s cost per key generation. Blocks mass key creation. |
| Short TTL | MUST | Default 1h. Spam files self-delete. |
| Size limit | MUST | 500 MB — checked client-side before upload. |
| PIN throttling | MUST | Already implemented in `pin_verify`. No changes. |

**PoW is only for API Key generation, not per-upload.** The key + rate limiting is sufficient for upload protection.

---

## 11. Testing

### Unit tests

- `go test` + `testify` (only external test dependency)
- Table-driven tests throughout
- Zero mocks — consistent with existing project standard (300+ tests)

Coverage:
- `crypto`: encrypt/decrypt round-trip for both modes, vault double-encryption, ZIP bundling, rejection on wrong key/PIN/password
- `api`: response parsing, error handling, retry
- `auth`: PoW solving, key save/load, TTL check
- `commands`: argument parsing, JSON output

### Browser ↔ Go compatibility tests (most critical)

```
Go encrypt → Node.js/Vitest decrypt
Browser encrypt (fixture) → Go decrypt
```

Run in CI on every PR touching `internal/crypto`.

### Integration tests

- `send` + `receive` in all modes against dev environment
- Multi-file: 3 files → receive → unzip → verify all 3
- PIN: wrong PIN → verify error
- Password: wrong password → verify error

---

## 12. CI/CD and Release

### GitHub Actions (`.github/workflows/cli.yml`)

Triggers only on changes to `cli/**`.

| Trigger | Actions |
|---|---|
| PR → main | `go test ./...` + `go vet` + `golangci-lint` |
| Push → main | Tests + cross-compile for all platforms |
| Tag `cli/v*.*.*` | GoReleaser v2: binaries + GitHub Release + Homebrew + Scoop |

### Linting (`.golangci.yml`)

`golangci-lint` v1.57+ with: `errcheck`, `gosimple`, `govet`, `ineffassign`, `staticcheck`, `gosec`, `gofmt`.

### Release platforms

| OS | Arch | Format |
|---|---|---|
| macOS | amd64 + arm64 | `.tar.gz` + Homebrew |
| Linux | amd64 + arm64 | `.tar.gz` |
| Windows | amd64 | `.zip` + Scoop |

---

## 13. Milestones

| M# | Name | Deliverables |
|---|---|---|
| M1 | Crypto layer | `internal/crypto` — URL + PIN + password modes. Unit tests. Browser compatibility tests. |
| M2 | API client | `internal/api` — all endpoints. Config. Error handling and retry. |
| M3 | Auth | `internal/auth`: PoW, key generation, storage, auto-renewal. Lambda `auth_init` + `auth_confirm`. |
| M4 | send + receive (URL + PIN) | Commands for URL and PIN modes. Progress bars. `--json` flag. |
| M5 | Password + multi-file | Password mode. ZIP bundling. `securedbx status`. |
| M6 | Release | GoReleaser. Homebrew Formula. Scoop. README. |

---

## 14. Open Questions

- **PoW difficulty**: starting with 4 leading zeros (~0.5s). Adjust after load testing.
- **Shell completions**: cobra supports bash/zsh/fish auto-generation — add to M6?
- **QR code**: `securedbx send --qr` — print QR to terminal. Low priority.

---

*SecureDBX · Zero-knowledge file & text sharing · sdbx.cc*
