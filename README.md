# SecureDBX

> Share files and secrets that self-destruct. End-to-end encrypted, zero-knowledge, no account needed.

**Everything is encrypted on your device. The server never has access to your data.**

## Features

### 🔐 Security First
- **Client-side encryption** — Files and text are encrypted in your browser before upload (AES-256-GCM)
- **Zero-knowledge architecture** — The server never sees your encryption keys or unencrypted content
- **One-time access** — Each file/text can only be viewed once, then it's automatically deleted
- **Network failure protection** — 10-minute retry window if download is interrupted

### 📝 Four Sharing Modes
- **URL Link** — Share files or text via encrypted URL. Key is embedded in the URL fragment (one-time access)
- **PIN Code** — Share files or text protected by a 4-8 digit PIN. Uses PBKDF2 key derivation with server-side salt (configurable: one-time or keep until expiry)
- **Text Secrets** — Share encrypted text snippets (up to 1000 characters) via URL or PIN without uploading files
- **Vault** — Password-protected sharing with unlimited access until expiry (multi-access)

### 🚫 No Tracking, No Accounts
- **No registration** — No accounts, no email, no cookies
- **Privacy-first** — No analytics, no IP logging, no user tracking
- **Self-destructing** — Content expires after 1-24 hours or custom time (5 min - 7 days)

### 🎨 Modern Interface
- **Clean design** — Built with Tailwind CSS for a modern, professional look
- **Dark mode** — Automatic dark/light theme with manual toggle
- **Responsive** — Works seamlessly on desktop and mobile devices
- **Accessible** — WCAG 2.1 AA compliant

### 📦 Convenient Sharing
- **Multi-file bundles** — Upload up to 10 files at once, automatically bundled into an encrypted ZIP
- **QR codes** — Generate QR codes for easy mobile sharing
- **Flexible expiry** — One-time delete after download or keep until expiry (PIN mode)

### 🛡️ Bot Protection
- **reCAPTCHA v3** — Invisible bot detection on all POST endpoints
- **CloudFront origin verification** — Blocks direct API access
- **Rate limiting** — Prevents abuse without compromising legitimate use

---

## This Service Is For You If...

- You want to share a file or secret and be certain it disappears forever after viewing
- You don't trust services that store your files on their servers indefinitely
- You need to securely share passwords, API keys, or confidential documents without registration
- You value privacy and full control over your data — no logs, no tracking
- You need a simple "view once → gone" approach, with no history or repeated access

---

## How It Works

### Uploading Files or Secrets

1. Select a file or enter text and choose expiration time (1h, 12h, 24h, or custom 5min-7days)
2. Your browser generates a random encryption key
3. Content is encrypted locally using AES-256-GCM
4. Encrypted data is uploaded to secure storage (server never sees the key or original filename)
5. You receive a shareable link with the decryption key in the URL fragment (`#key#filename`)
6. **The key never leave your browser or reach the server**

### Downloading/Viewing

1. Open the shared link
2. Your browser extracts the decryption key from the URL fragment
3. Click to view/download (warning: this works only once!)
4. **Download reserved** - Server reserves the content for your download (10-minute window)
5. Content is downloaded and decrypted in your browser:
   - **Files**: Downloaded and saved with original filename
   - **Text**: Displayed in a copyable text area
6. **Download confirmed** - Your browser confirms successful download to the server
7. Content is automatically deleted from the server
8. The link becomes invalid forever

**Network Failure Protection**: If the download is interrupted (network error, browser crash), you can retry within 10 minutes using the same link. This ensures that you don't lose access due to temporary technical problems while maintaining the security of one-time downloads.

### Vault (Password-Protected Multi-Access)

1. Select the **Vault** tab and choose File or Text
2. Enter your content and set a password (min 4 characters)
3. Your browser:
   - Generates a random data encryption key
   - Derives a password key using PBKDF2 (100,000 iterations, SHA-256)
   - Encrypts the data key with the password-derived key
   - Encrypts content with the data key (AES-256-GCM)
4. Encrypted content + encrypted key + salt are uploaded
5. You receive a link (password NOT included - share separately!)
6. Recipients enter the password to decrypt and access content
7. **Multi-access**: Can be downloaded unlimited times until expiry

### PIN Code Sharing

1. Select the **PIN** method and choose File or Text
2. Enter a 4-8 digit PIN code and choose expiration time
3. Your browser requests a unique salt from the server
4. A key is derived from PIN + salt using PBKDF2 (100,000 iterations, SHA-256)
5. Content is encrypted with the derived key (AES-256-GCM)
6. Encrypted content is uploaded to secure storage
7. You receive a short File ID to share along with the PIN
8. **Toggle option**: Delete after first download, or keep until expiry
9. Recipient enters the File ID and PIN to decrypt and access content

---

## Security

### Encryption
- **Algorithm**: AES-256-GCM with 256-bit keys
- **Key generation**: Cryptographically secure random values (`crypto.getRandomValues`)
- **Zero-knowledge**: Decryption keys and filenames stay in URL fragments, never sent to server
- **HTTPS only**: All communication encrypted in transit (TLS 1.2+)

### Protection Layers

**Layer 1 - CloudFront Origin Verification** ✅ Deployed
- Custom secret header blocks direct API calls
- Only requests through CloudFront are accepted
- Prevents automated script abuse
- Zero cost protection

**Layer 2 - reCAPTCHA v3** ✅ Deployed
- Bot detection and prevention on all POST endpoints
- Invisible to legitimate users (score-based verification)
- Protects upload, download, and abuse report endpoints
- Minimum score threshold: 0.5 (adjustable)

**Layer 3 - Two-Phase Download with Reservation System** ✅ Deployed
- **Phase 1 - Reserve**: Atomic DynamoDB update reserves content for download
  - 10-minute reservation window
  - Prevents race conditions when multiple users try to download
  - Allows retry if network fails during download
- **Phase 2 - Confirm**: Browser confirms successful download
  - Content is permanently deleted only after confirmation
  - Statistics are updated atomically
  - No data loss from temporary network issues

### What We DON'T Store
- ❌ Email addresses or user names
- ❌ File contents (only encrypted blobs)
- ❌ Original filenames
- ❌ Encryption keys
- ❌ IP addresses (only hashed for abuse prevention)

### What We DO Store
- ✅ File size (metadata)
- ✅ Upload timestamp
- ✅ Expiration time
- ✅ Download status (reserved/confirmed)
- ✅ IP address hash (HMAC-SHA256 with secret salt, for abuse prevention only)

---

## Limitations

- Maximum file size: 500 MB
- Files expire after max 7 days (custom) or 24 hours (preset)
- URL mode: each file can only be downloaded once (with 10-minute retry window). PIN mode: configurable one-time or keep until expiry
- Multi-file upload: up to 10 files per bundle
- Maximum text secret size: 1,000 characters
- Desktop browsers recommended (mobile support available)

---

## Technology

### High level architecture

![High level architecture](architecture.jpeg)

### Frontend
- **Vue 3** - Component-based SPA with Vue Router (hash mode) and Pinia stores
- **Vite** - Fast build tooling with HMR dev server
- **Tailwind CSS** - Modern utility-first CSS framework
- **Web Crypto API** - AES-256-GCM encryption, PBKDF2 key derivation
- **Dark Mode** - CSS class-based theme switching with localStorage persistence
- **Web Workers** - Background encryption/decryption for large files
- **Responsive Design** - Mobile-first approach

### Backend
- **AWS Lambda** - 10 serverless functions (Python 3.12)
- **API Gateway** - RESTful API with request validation
- **Lambda Layers** - Shared dependencies (boto3, requests)

### Storage
- **Amazon S3** - Encrypted file storage with lifecycle policies
- **DynamoDB** - Metadata with TTL and atomic operations
- **Global Statistics** - Privacy-safe aggregate metrics

### CDN & Security
- **CloudFront** - Global CDN with custom security headers
- **Origin Access Identity** - Secure S3 access
- **reCAPTCHA v3** - Invisible bot protection

### Infrastructure
- **Terraform** - Infrastructure as Code
- **GitHub Actions** - CI/CD pipeline (planned)
- **CloudWatch** - Monitoring and alerting

---

## CLI

`sdbx` is a command-line client for securedbx.com. It encrypts your files and secrets locally before uploading — the server never sees your data.

### Installation

Download the binary for your platform from [Releases](https://github.com/AntonAks/securedbx/releases):

### Linux (amd64)
```bash
curl -L https://github.com/AntonAks/securedbx/releases/latest/download/sdbx_linux_amd64.tar.gz | tar xz
sudo mv sdbx /usr/local/bin/
```
### macOS (Apple Silicon)
```bash
curl -L https://github.com/AntonAks/securedbx/releases/latest/download/sdbx_darwin_arm64.tar.gz | tar xz
sudo mv sdbx /usr/local/bin/
```
### macOS (Intel)
```bash
curl -L https://github.com/AntonAks/securedbx/releases/latest/download/sdbx_darwin_amd64.tar.gz | tar xz
sudo mv sdbx /usr/local/bin/
```

### Usage

**Send a file:**
```bash
sdbx send document.pdf
# {"file_id": "482973", "pin": "aB3x", "expires_at": 1234567890}
```

**Send a text secret:**
```bash
sdbx send --text "my secret password"
# {"file_id": "736104", "pin": "Tz9k", "expires_at": 1234567890}
```

**Receive a file or secret:**
```bash
sdbx receive 482973
# Enter PIN: ****
# {"type": "file", "path": "./document.pdf"}
```

**Options for `send`:**
```
--text "..."       Send text instead of a file
--ttl 1-24         Expiry in hours (default: 1)
--pin-value xxxx   Custom PIN (4 chars, a-z A-Z 0-9); auto-generated if omitted
```

### How it works

1. CLI generates a random 4-character PIN (or uses `--pin-value`)
2. A key is derived from PIN + server salt using PBKDF2-HMAC-SHA256 (100k iterations)
3. Content is encrypted locally with AES-256-GCM
4. Encrypted data is uploaded — the server never sees the original content or the PIN
5. Recipient uses `file_id` + PIN to decrypt on their side

The recipient can download via CLI or the web UI at [securedbx.com](https://securedbx.com).

### Configuration
API keys are obtained automatically via Proof-of-Work challenge and cached in `~/.config/securedbx/config.json` (24h TTL).


---

## Development

### Prerequisites

- AWS Account
- Terraform >= 1.5
- Python 3.12
- Node.js >= 18 (for frontend build)
- AWS CLI configured (`aws configure`)

### Quick Start - Deploy to AWS

```bash
# 1. Bootstrap Terraform backend (once per AWS account)
./scripts/bootstrap-terraform-backend.sh

# 2. Build Lambda packages
make build-lambdas-dev

# 3. Build frontend CSS
make build-frontend

# 4. Initialize Terraform
make init-dev

# 5. Deploy infrastructure
make deploy-dev

# 6. Get your CloudFront URL
cd terraform/environments/dev && terraform output cloudfront_domain
```

### Salt Management

IP addresses are hashed using HMAC-SHA256 with a secret salt stored in AWS Systems Manager Parameter Store. This prevents rainbow table attacks against the IP hash database.

```bash
# Initialize salt (once per environment)
make init-salt-dev
make init-salt-prod

# Check if salt exists
make check-salt-dev
make check-salt-prod
```

The salt is:
- KMS encrypted in Parameter Store (never in Terraform state or code)
- Cached in Lambda memory (one API call per cold start)
- Audited via CloudTrail
- Free (within AWS free tier)

### Local Frontend Development

```bash
# Install dependencies
cd frontend && npm install

# Start Vite dev server (with HMR)
make dev-frontend
# Visit http://localhost:5173

# Build for production
make build-frontend

# The dev server proxies API requests to your deployed backend.
# Set VITE_API_BASE=/dev for dev environment builds.
```

### Project Structure

```
securedbx/
├── cli/                   # Go CLI client (sdbx)
│   ├── cmd/sdbx/          # Entry point
│   ├── internal/
│   │   ├── api/           # API client
│   │   ├── auth/          # PoW auth + API key lifecycle
│   │   ├── commands/      # send, receive commands
│   │   ├── config/        # Config and env vars
│   │   └── crypto/        # AES-256-GCM, PBKDF2, ZIP bundling
│   └── go.mod
├── backend/
│   ├── lambdas/           # 10 Lambda functions
│   │   ├── upload_init/
│   │   ├── download/
│   │   ├── confirm_download/
│   │   ├── get_metadata/
│   │   ├── get_stats/
│   │   ├── cleanup/
│   │   ├── report_abuse/
│   │   ├── pin_initiate/    # PIN: generate salt
│   │   ├── pin_upload_init/ # PIN: initialize upload
│   │   └── pin_verify/      # PIN: verify & download
│   ├── shared/            # Shared utilities
│   │   ├── constants.py
│   │   ├── dynamo.py      # DynamoDB operations
│   │   ├── s3.py          # S3 operations
│   │   ├── security.py    # Security decorators
│   │   ├── validation.py
│   │   ├── response.py
│   │   └── ...
│   └── tests/             # 226 backend tests
├── frontend/              # Vue 3 SPA (Vite)
│   ├── src/
│   │   ├── main.js        # App entry point
│   │   ├── App.vue        # Root component
│   │   ├── router/        # Vue Router (hash mode)
│   │   ├── stores/        # Pinia stores (theme, stats)
│   │   ├── lib/
│   │   │   ├── crypto.js  # AES-256-GCM encryption
│   │   │   ├── zip-bundle.js  # Multi-file ZIP bundling
│   │   │   ├── utils.js   # Shared utilities
│   │   │   └── api.js     # API base URL config
│   │   ├── workers/       # Web Workers for large files
│   │   ├── composables/   # Vue composables (upload, download, etc.)
│   │   ├── components/    # Reusable Vue components
│   │   └── views/         # Page views (Home, Download, Share, FAQ, About)
│   ├── tests/             # 82 frontend tests (Vitest)
│   ├── index.html         # SPA entry point
│   └── vite.config.js     # Vite + Vue plugin config
├── terraform/
│   ├── environments/
│   │   ├── dev/
│   │   └── prod/
│   └── modules/
│       ├── storage/       # S3 + DynamoDB
│       ├── api/           # API Gateway + Lambdas
│       ├── cdn/           # CloudFront
│       └── monitoring/    # CloudWatch
└── scripts/               # Deployment scripts
```

### Deployment Status

- ✅ Development environment deployed and operational
- ✅ Production environment deployed and operational
- ✅ Modern UI with Tailwind CSS and dark mode
- ✅ Network failure protection implemented

---

## Testing

securedbx has a comprehensive test suite with **300+ tests** and **ZERO mocks**:

```bash
# Run all tests (backend + frontend)
make test

# Run backend tests only (Python)
make test-backend

# Run frontend tests only (JavaScript)
make test-frontend

# Run backend tests with coverage report
make test-backend-cov
```

**Test Coverage:**
- **Backend**: 226 tests covering validation, response formatting, JSON encoding, security, PIN flows
- **Frontend**: 82 tests (Vitest) covering AES-256-GCM encryption, PBKDF2 key derivation, PIN crypto, multi-file bundles
- **Coverage**: ~80% of shared modules

**What's Tested:**
- ✅ AES-256-GCM encryption/decryption round-trips
- ✅ Key generation and export/import
- ✅ Large file handling (1 MB)
- ✅ Unicode text preservation
- ✅ Wrong key rejection
- ✅ Input validation (file ID, size, TTL)
- ✅ Response formatting and CORS headers
- ✅ Decimal encoding for DynamoDB
- ✅ CloudFront origin verification
- ✅ PBKDF2 key derivation (PIN mode)
- ✅ PIN upload/download encryption round-trips
- ✅ Multi-file ZIP bundle encryption
- ✅ Edge cases and error handling

---

## Privacy Policy

### Simple Version: We Collect Nothing

- We don't store unencrypted files
- We don't store or have access to filenames
- We don't have access to your encryption keys
- We don't log IP addresses for file access
- We don't use cookies, tracking, or analytics
- Files are automatically deleted after the first download or expiration

### What We Store

The only data we store is:
- Encrypted file blobs (which we cannot decrypt)
- File size, upload timestamp, and expiration time
- Download reservation status
- IP address hash (HMAC-SHA256 with secret salt, for abuse prevention only, not linked to files)

**That's it.** No personal information, no tracking, no analytics.

---

## Roadmap

Planned features for securedbx:

- ✅ **Multiple Files / Zip Bundle** - Upload multiple files as encrypted bundle
- ✅ **Custom Expiration Times** - Precise expiration (5 min - 7 days) with real-time preview
- ✅ **Vault (Password Protection)** - Password-protected multi-access sharing with PBKDF2 key derivation
- ✅ **PIN Code Sharing** - Share via numeric PIN with PBKDF2 key derivation and configurable access mode
- ✅ **Short URLs** - Shorter file IDs for cleaner links
- ✅ **CLI Client** - Command-line tool (`sdbx`) for sending and receiving encrypted files/secrets
- 📋 **IP/Geo Restriction** - Restrict downloads by country or IP
- 📋 **Self-destructing Voice Message** - Encrypted audio messages
- 📋 **Dead Man's Switch** - Auto-share if user doesn't check in

See [FEATURE_ROADMAP.md](./FEATURE_ROADMAP.md) for full details and implementation plans.

---

## Contributing

Contributions welcome! Please read our [contribution guidelines](./CONTRIBUTING.md) first.

### How to Contribute

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Areas Where We Need Help

- Security audits
- Performance optimization
- Documentation improvements
- UI/UX enhancements
- Additional test coverage
- Feature implementations (see [FEATURE_ROADMAP.md](./FEATURE_ROADMAP.md))

---

## License

MIT License - See [LICENSE](./LICENSE) file for details

---

## Disclaimer

This service is provided as-is. While we implement strong encryption and security practices, users are responsible for their own data. Do not use securedbx for illegal content.

For sensitive files, we recommend additional security measures like password-protecting archives before upload.

---

## Acknowledgments

Built with privacy in mind 🔒

Special thanks to:
- Web Crypto API for browser-native encryption
- AWS for reliable serverless infrastructure
- Tailwind CSS for beautiful utility-first styling
- The open-source community

---

**Questions or Issues?**

- 🐛 [Report a Bug](https://github.com/antonaks/securedbx/issues)
- 💡 [Request a Feature](https://github.com/antonaks/securedbx/issues)
- 🔒 [Report Security Issue](./SECURITY.md)
