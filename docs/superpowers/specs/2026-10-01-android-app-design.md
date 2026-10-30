# Android App (TWA) — Design

**Date:** 2026-10-01
**Status:** Approved approach: Trusted Web Activity (TWA) wrapping the existing Vue frontend.

## Goal

An Android app for SecureDBX that provides:
1. A home-screen icon (app presence without a browser).
2. A "Share to SecureDBX" entry in the Android share sheet for files and text.
3. securedbx.com links opening in the app instead of messenger in-app browsers (fixes broken link opening from Telegram/WhatsApp).

Distribution: APK via GitHub Releases first; Google Play later as a separate step.

## Non-Goals

- No native Kotlin UI, biometrics, offline queue, or push notifications (YAGNI until a native-only need appears).
- No second codebase: the app renders the production site; app updates ship via site deploys.

## Architecture

The existing Vue app at https://securedbx.com becomes an installable PWA. Bubblewrap generates a thin TWA shell APK that opens the site in full-screen Chrome. Verified Android App Links bind the domain to the app.

### 1. PWA layer (frontend/ — the only code changes)

- `frontend/public/manifest.webmanifest`
  - name/short_name, `display: standalone`, theme/background colors.
  - Icons: 192px and 512px, plus maskable variants (derive from existing `favicon-192.png` / brand assets).
  - `share_target`: `{ action: "/share-target", method: "POST", enctype: "multipart/form-data", params: { title, text, url, files: [{ name: "files", accept: ["*/*"] }] } }`.
- `frontend/src/sw.js` — minimal service worker:
  - Exists only for installability and share-target handling.
  - **Network-only for all requests** — no caching. The service deals in one-time links; a cached API response could burn a download or serve stale state. Static-asset caching can be added later if ever needed.
  - Intercepts `POST /share-target`: stores shared files/text in the Cache API under a known key, responds with a redirect to `/#/share-target` (the router is hash-based — `createWebHashHistory`).
- New Vue route `/#/share-target` (path `/share-target` is distinct from the existing `/share` route):
  - Reads the stashed files/text from the Cache API, clears the stash, and pre-fills the existing upload form flow (HomeView / DropZone) with them.
- `index.html`: link the manifest; register the service worker in `main.js`.

### 2. App Links

- `frontend/public/.well-known/assetlinks.json` with the SHA-256 fingerprint of the APK signing key (`delegate_permission/common.handle_all_urls`).
- Served at `https://securedbx.com/.well-known/assetlinks.json`. Verify the CloudFront/hosting config serves `.well-known` paths with `content-type: application/json` and no redirect.
- Effect: tapping a securedbx.com link in messengers opens the app directly; also removes the TWA URL bar.

### 3. TWA shell (android/)

- Generated with Bubblewrap: `bubblewrap init --manifest https://securedbx.com/manifest.webmanifest`, then `bubblewrap build`.
- Committed to repo: `android/twa-manifest.json` + a short `android/README.md` (how to regenerate/build/sign).
- NOT committed: keystore, generated Android project, APK. Signing key kept locally / in CI secrets.
- APK published to GitHub Releases.

## Data Flow (share path)

Gallery → Android share sheet → SecureDBX → Chrome POSTs multipart to `/share-target` → SW stashes payload in Cache API, redirects to `/#/share-target` → Vue route pulls payload, clears stash → upload form pre-filled → normal existing encrypt/upload flow.

## Error Handling

- Share-target route with empty/missing stash (e.g. page reload): redirect to home, no error dialog.
- Shared files over the existing size limit: handled by the existing form validation — no new logic.
- SW registration failure: site works exactly as today (progressive enhancement).

## Testing

- Vitest: share-target payload parsing/stash-read logic (store files → read & clear → empty-stash case).
- Manual device run: install APK, share a photo from the gallery, open a securedbx.com link from Telegram, complete one send/receive round-trip.

## Build Order

1. PWA layer (manifest + SW + share-target route) — independently valuable: site becomes installable from Chrome.
2. Bubblewrap init/build + assetlinks.json deploy.
3. GitHub Release with APK + install instructions in README.

## Risks / Notes

- TWA requires Chrome (or a Custom-Tabs-capable browser) on the device — present on virtually all Android devices.
- Web Share Target requires the PWA to be installed via Chrome/TWA — expected; that IS the app.
- Hash-based routing means shared/encrypted-key URLs keep working unchanged; App Links match on domain, fragments irrelevant.
