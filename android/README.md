# SecureDBX Android (TWA)

Thin Trusted Web Activity wrapper around https://securedbx.com — the PWA
(manifest + service worker) lives in `frontend/`, this folder only holds the
build config. The generated Android project, keystore and APKs are NOT
committed.

## Prerequisites

- The PWA layer must be **deployed to securedbx.com first** (bubblewrap reads
  the live web manifest).
- Node 20+, JDK 17 (bubblewrap can download its own JDK/Android SDK on first run).

## One-time setup

```bash
npm i -g @bubblewrap/cli
cd android
bubblewrap init --manifest https://securedbx.com/manifest.webmanifest
# answers are pre-filled from twa-manifest.json; when asked to create a
# signing key, say yes — keep android.keystore OUT of git and back it up:
# losing it means Play Store updates become impossible later.
```

Get the signing key fingerprint (needed for assetlinks.json):

```bash
keytool -list -v -keystore android.keystore -alias securedbx | grep SHA256
```

Put that fingerprint into `frontend/public/.well-known/assetlinks.json` and
deploy the site again — until then Android won't verify the domain link and
the app shows a URL bar.

## Build the APK

```bash
cd android
bubblewrap build
# → app-release-signed.apk
```

Install on a phone: transfer the APK (or `adb install app-release-signed.apk`),
allow installs from unknown sources. Publish via GitHub Releases until the
Play Store step.

## Updating

The app is the website — site deploys ARE app updates. Rebuild the APK only
when this config changes (name, icons, packageId, version for Play).
