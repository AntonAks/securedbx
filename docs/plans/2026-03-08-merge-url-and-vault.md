# Merge URL and Vault Modes Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Fold vault's password-protection and multi-access features into the File and Text URL tabs, then delete the Vault tab entirely.

**Architecture:** Both `LinkUpload.vue` and `TextUpload.vue` get two optional toggle switches (password, multiple downloads). When password is ON the upload uses vault-style PBKDF2 crypto and the share URL uses `?salt=...&vault=1` so the existing `VaultDownload.vue` handles decryption. When multiple-downloads is ON `access_mode: 'multi'` is sent and the TtlSelector is clamped to ≤ 24 h. Backend `upload_init` and `get_metadata` are extended to support salt/encrypted_key for any access_mode (not just multi).

**Tech Stack:** Vue 3, Vite, Tailwind CSS, Web Crypto API, Python 3.12 AWS Lambda, DynamoDB, pytest

---

## Task 1: Add `limitTo24h` prop to `TtlSelector.vue`

**Files:**
- Modify: `frontend/src/components/TtlSelector.vue`

`TtlSelector` delegates state to `useTtl.js` composable. We only need to: hide the Custom radio when `limitTo24h=true`, and snap back to `24h` if currently on custom when the prop turns on.

**Step 1: Add the prop and a watcher**

Replace the `<script setup>` in `TtlSelector.vue`:

```vue
<script setup>
import { watch } from 'vue';
import { useTtl } from '../composables/useTtl.js';

const props = defineProps({
  showCustom: { type: Boolean, default: true },
  limitTo24h: { type: Boolean, default: false },
});

const {
  selectedPreset,
  customUnit,
  customValue,
  isCustom,
  unitOptions,
  expirationPreview,
  getTtl,
} = useTtl();

watch(() => props.limitTo24h, (limited) => {
  if (limited && selectedPreset.value === 'custom') {
    selectedPreset.value = '24h';
  }
});

defineExpose({ getTtl });
</script>
```

**Step 2: Hide Custom radio when `limitTo24h` is true**

In the template, change the custom option's `v-if`:

```html
<label v-if="showCustom && !limitTo24h" class="flex items-center cursor-pointer">
```

**Step 3: Run frontend tests**

```bash
cd frontend && npx vitest run
```

Expected: all 82 tests pass (no test touches TtlSelector directly).

**Step 4: Commit**

```bash
git add frontend/src/components/TtlSelector.vue
git commit -m "feat: add limitTo24h prop to TtlSelector"
```

---

## Task 2: Backend — allow salt/encrypted_key for any access_mode

**Files:**
- Modify: `backend/lambdas/upload_init/handler.py`
- Modify: `backend/lambdas/get_metadata/handler.py`

**Why:** New combinations — password + single download (`access_mode: one_time` + salt + encrypted_key) — need backend support. Currently salt/encrypted_key are only processed for multi mode.

### `upload_init/handler.py`

**Step 1: Change the salt/encrypted_key extraction**

Find this block (lines 117–123):
```python
        salt = None
        encrypted_key = None
        if access_mode == ACCESS_MODE_MULTI:
            salt = body.get("salt")
            encrypted_key = body.get("encrypted_key")
            validate_salt(salt)
            validate_encrypted_key(encrypted_key)
```

Replace with:
```python
        salt = None
        encrypted_key = None
        raw_salt = body.get("salt")
        raw_encrypted_key = body.get("encrypted_key")
        if raw_salt or raw_encrypted_key:
            validate_salt(raw_salt)
            validate_encrypted_key(raw_encrypted_key)
            salt = raw_salt
            encrypted_key = raw_encrypted_key
```

### `get_metadata/handler.py`

**Step 2: Return salt/encrypted_key for any record that has them**

Find this block (lines 69–73):
```python
        # Include vault-specific fields for multi-access mode
        if access_mode == "multi":
            response_data["salt"] = record.get("salt")
            response_data["encrypted_key"] = record.get("encrypted_key")
            response_data["download_count"] = record.get("download_count", 0)
```

Replace with:
```python
        # Include password fields if present (any access mode)
        if record.get("salt"):
            response_data["salt"] = record.get("salt")
            response_data["encrypted_key"] = record.get("encrypted_key")
        if access_mode == "multi":
            response_data["download_count"] = record.get("download_count", 0)
```

**Step 3: Run backend tests**

```bash
cd backend && . venv/bin/activate && pytest tests/ -v --tb=short
```

Expected: all 226 tests pass.

**Step 4: Commit**

```bash
git add backend/lambdas/upload_init/handler.py backend/lambdas/get_metadata/handler.py
git commit -m "feat: allow salt/encrypted_key for any access_mode in upload_init and get_metadata"
```

---

## Task 3: Add i18n strings for the new toggles

**Files:**
- Modify: `frontend/src/i18n/locales/en.json`
- Modify: `frontend/src/i18n/locales/uk.json`

**Step 1: Add to `en.json`** under `"upload"` → add a new `"options"` section (put it after the `"link"` section):

```json
"options": {
  "passwordLabel": "Password protect",
  "passwordPlaceholder": "Enter password...",
  "passwordHint": "Share the password separately from the link",
  "generatePassword": "Generate",
  "copyPassword": "Copy",
  "passwordCopied": "Copied!",
  "multiDownloadLabel": "Allow multiple downloads",
  "multiDownloadHint": "Link stays active until expiry (max 24h). Without this, link works once."
}
```

**Step 2: Add to `uk.json`** (Ukrainian translation, under `"upload"`):

```json
"options": {
  "passwordLabel": "Захист паролем",
  "passwordPlaceholder": "Введіть пароль...",
  "passwordHint": "Поділіться паролем окремо від посилання",
  "generatePassword": "Згенерувати",
  "copyPassword": "Копіювати",
  "passwordCopied": "Скопійовано!",
  "multiDownloadLabel": "Дозволити кілька завантажень",
  "multiDownloadHint": "Посилання діє до закінчення терміну (макс. 24г). Без цього — лише одноразове."
}
```

**Step 3: Commit**

```bash
git add frontend/src/i18n/locales/en.json frontend/src/i18n/locales/uk.json
git commit -m "feat: add i18n strings for password and multi-download toggles"
```

---

## Task 4: Rewrite `LinkUpload.vue` with password + multi-download toggles

**Files:**
- Modify: `frontend/src/views/home/LinkUpload.vue`

This is the file-upload tab. It needs the two new toggle sections injected before the TtlSelector.

**Step 1: Add the toggle UI + updated script**

The new template structure (after `<FileList>` and before `<TtlSelector>`):

```html
<!-- Password Toggle -->
<div class="mb-4">
  <label class="flex items-center cursor-pointer gap-3 mb-2">
    <button type="button" role="switch" :aria-checked="usePassword"
      :class="['relative inline-flex h-6 w-11 shrink-0 rounded-full border-2 border-transparent transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2', usePassword ? 'bg-blue-600' : 'bg-gray-300 dark:bg-slate-600']"
      @click="usePassword = !usePassword">
      <span :class="['pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow transform transition-transform', usePassword ? 'translate-x-5' : 'translate-x-0']"></span>
    </button>
    <span class="text-gray-700 dark:text-slate-300 font-medium">{{ $t('upload.options.passwordLabel') }}</span>
  </label>
  <div v-if="usePassword" class="flex items-center gap-2 mt-2">
    <input
      :type="showPassword ? 'text' : 'password'"
      v-model="password"
      :placeholder="$t('upload.options.passwordPlaceholder')"
      class="flex-1 px-3 py-2 bg-white dark:bg-slate-900 border border-gray-300 dark:border-slate-600 rounded text-gray-900 dark:text-slate-200 focus:outline-none focus:border-blue-500"
    />
    <button type="button" @click="showPassword = !showPassword"
      class="px-3 py-2 text-gray-500 hover:text-gray-700 dark:text-slate-400 dark:hover:text-slate-200 text-sm border border-gray-300 dark:border-slate-600 rounded">
      {{ showPassword ? $t('upload.vault.hide') : $t('upload.vault.show') }}
    </button>
    <button type="button" @click="generatePassword"
      class="px-3 py-2 text-gray-500 hover:text-gray-700 dark:text-slate-400 dark:hover:text-slate-200 text-sm border border-gray-300 dark:border-slate-600 rounded">
      {{ $t('upload.options.generatePassword') }}
    </button>
    <button type="button" @click="copyPasswordToClipboard"
      class="px-3 py-2 text-gray-500 hover:text-gray-700 dark:text-slate-400 dark:hover:text-slate-200 text-sm border border-gray-300 dark:border-slate-600 rounded">
      {{ passwordCopied ? $t('upload.options.passwordCopied') : $t('upload.options.copyPassword') }}
    </button>
  </div>
  <p v-if="usePassword" class="text-gray-500 dark:text-slate-500 text-xs mt-1">{{ $t('upload.options.passwordHint') }}</p>
</div>

<!-- Multiple Downloads Toggle -->
<div class="mb-6">
  <label class="flex items-center cursor-pointer gap-3">
    <button type="button" role="switch" :aria-checked="multiDownload"
      :class="['relative inline-flex h-6 w-11 shrink-0 rounded-full border-2 border-transparent transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2', multiDownload ? 'bg-blue-600' : 'bg-gray-300 dark:bg-slate-600']"
      @click="multiDownload = !multiDownload">
      <span :class="['pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow transform transition-transform', multiDownload ? 'translate-x-5' : 'translate-x-0']"></span>
    </button>
    <span class="text-gray-700 dark:text-slate-300 font-medium">{{ $t('upload.options.multiDownloadLabel') }}</span>
  </label>
  <p class="text-gray-500 dark:text-slate-500 text-xs mt-1 ml-14">{{ $t('upload.options.multiDownloadHint') }}</p>
</div>
```

Pass `limitTo24h` to TtlSelector:
```html
<TtlSelector ref="ttlRef" :limitTo24h="multiDownload" />
```

**Step 2: Full new `<script setup>`**

```js
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import DropZone from '../../components/DropZone.vue';
import FileList from '../../components/FileList.vue';
import TtlSelector from '../../components/TtlSelector.vue';
import ProgressBar from '../../components/ProgressBar.vue';
import { useUpload } from '../../composables/useUpload.js';
import { useRecaptcha } from '../../composables/useRecaptcha.js';
import { useClipboard } from '../../composables/useClipboard.js';
import { validateFiles, createBundle } from '../../lib/zip-bundle.js';
import { API_BASE } from '../../lib/api.js';
import * as CryptoModule from '../../lib/crypto.js';

const { t } = useI18n();
const router = useRouter();
const upload = useUpload();
const { getToken } = useRecaptcha();
const { copied: passwordCopied, copy: copyToClipboard } = useClipboard();
const ttlRef = ref(null);
const selectedFiles = ref([]);

// New options
const usePassword = ref(false);
const password = ref('');
const showPassword = ref(false);
const multiDownload = ref(false);

function addFiles(files) {
    const combined = [...selectedFiles.value, ...Array.from(files)];
    const validation = validateFiles(combined);
    if (!validation.valid) { alert(validation.error); return; }
    selectedFiles.value = combined;
}
function removeFile(index) { selectedFiles.value.splice(index, 1); }
function clearFiles() { selectedFiles.value = []; }

function generatePassword() {
    const chars = 'ABCDEFGHJKMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789!@#$%';
    const arr = new Uint8Array(16);
    crypto.getRandomValues(arr);
    password.value = Array.from(arr).map(b => chars[b % chars.length]).join('');
    showPassword.value = true;
}

function copyPasswordToClipboard() {
    if (password.value) copyToClipboard(password.value);
}

async function handleUpload() {
    if (selectedFiles.value.length === 0) return;
    if (usePassword.value && password.value.length < 4) {
        alert('Password must be at least 4 characters');
        return;
    }

    try {
        upload.uploading.value = true;

        let fileToUpload, uploadFileName;
        if (selectedFiles.value.length === 1) {
            fileToUpload = selectedFiles.value[0];
            uploadFileName = selectedFiles.value[0].name;
            upload.updateProgress(0, t('upload.progress.preparingFile'));
        } else {
            upload.updateProgress(0, t('upload.progress.creatingZip'));
            const bundle = await createBundle(selectedFiles.value, (percent, message) => {
                upload.updateProgress(percent * 0.15, message);
            });
            fileToUpload = bundle.blob;
            uploadFileName = bundle.filename;
        }

        const accessMode = multiDownload.value ? 'multi' : 'one_time';
        const ttl = ttlRef.value.getTtl();

        if (usePassword.value) {
            // Vault-style: PBKDF2 double encryption
            upload.updateProgress(15, t('upload.progress.generatingKey'));
            const dataKey = await CryptoModule.generateKey();

            upload.updateProgress(18, t('upload.progress.derivingPasswordKey'));
            const salt = CryptoModule.generateSalt();
            const passwordKey = await CryptoModule.deriveKeyFromPassword(password.value, salt);

            upload.updateProgress(20, t('upload.progress.securingKey'));
            const encryptedKeyData = await CryptoModule.encryptKey(dataKey, passwordKey);
            const encryptedKeyBase64 = CryptoModule.arrayToBase64(encryptedKeyData);
            const saltBase64 = CryptoModule.arrayToBase64(salt);

            upload.updateProgress(25, t('upload.progress.encrypting', { percent: 0 }));
            const encryptedData = await CryptoModule.encryptFile(fileToUpload, dataKey, (progress) => {
                upload.updateProgress(25 + progress * 0.3, t('upload.progress.encrypting', { percent: Math.round(progress) }));
            });

            upload.updateProgress(55, t('upload.progress.initializing'));
            const recaptchaToken = await getToken('upload');
            const fileSize = encryptedData.byteLength;
            const response = await fetch(`${API_BASE}/upload/init`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    content_type: 'file',
                    file_size: fileSize,
                    ttl,
                    access_mode: accessMode,
                    salt: saltBase64,
                    encrypted_key: encryptedKeyBase64,
                    recaptcha_token: recaptchaToken,
                }),
            });
            if (!response.ok) {
                const error = await response.json().catch(() => ({}));
                throw new Error(error.error || `Upload failed: ${response.status}`);
            }
            const data = await response.json();

            upload.updateProgress(65, t('upload.progress.uploadingEncrypted'));
            await upload.uploadToS3(data.upload_url, encryptedData);

            upload.updateProgress(100, t('upload.progress.completeRedirect'));
            const encodedFileName = encodeURIComponent(uploadFileName);
            setTimeout(() => {
                router.push({ path: '/share', query: { id: data.file_id, salt: saltBase64, name: encodedFileName, vault: '1' } });
            }, 500);
        } else {
            // Standard: key in URL
            upload.updateProgress(15, t('upload.progress.generatingKey'));
            const key = await CryptoModule.generateKey();

            upload.updateProgress(20, t('upload.progress.encrypting', { percent: 0 }));
            const encryptedData = await CryptoModule.encryptFile(fileToUpload, key, (progress) => {
                upload.updateProgress(20 + progress * 0.3, t('upload.progress.encrypting', { percent: Math.round(progress) }));
            });

            upload.updateProgress(50, t('upload.progress.initializing'));
            const recaptchaToken = await getToken('upload');
            const fileSize = fileToUpload.size || fileToUpload.byteLength;
            const response = await fetch(`${API_BASE}/upload/init`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    content_type: 'file',
                    file_size: fileSize,
                    ttl,
                    access_mode: accessMode,
                    recaptcha_token: recaptchaToken,
                }),
            });
            if (!response.ok) {
                const error = await response.json().catch(() => ({}));
                throw new Error(error.error || `Upload failed: ${response.status}`);
            }
            const data = await response.json();

            upload.updateProgress(60, t('upload.progress.uploadingEncrypted'));
            await upload.uploadToS3(data.upload_url, encryptedData);

            upload.updateProgress(100, t('upload.progress.completeRedirect'));
            const keyBase64 = await CryptoModule.keyToBase64(key);
            const encodedFileName = encodeURIComponent(uploadFileName);
            setTimeout(() => {
                router.push({ path: '/share', query: { id: data.file_id, key: keyBase64, name: encodedFileName } });
            }, 500);
        }
    } catch (error) {
        console.error('Upload error:', error);
        alert(error.message || t('upload.progress.uploadFailed'));
        upload.uploading.value = false;
    }
}
```

**Step 3: Run frontend tests**

```bash
cd frontend && npx vitest run
```

**Step 4: Commit**

```bash
git add frontend/src/views/home/LinkUpload.vue
git commit -m "feat: add password and multi-download options to file URL upload"
```

---

## Task 5: Rewrite `TextUpload.vue` with the same toggles

**Files:**
- Modify: `frontend/src/views/home/TextUpload.vue`

Same two toggles as Task 4, but the upload logic works with text content.

**Step 1: Full new template** — add both toggle blocks (identical HTML as Task 4) between the char-count paragraph and the TtlSelector. Pass `:limitTo24h="multiDownload"` to TtlSelector.

**Step 2: Full new `<script setup>`**

```js
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import TtlSelector from '../../components/TtlSelector.vue';
import ProgressBar from '../../components/ProgressBar.vue';
import { useUpload } from '../../composables/useUpload.js';
import { useRecaptcha } from '../../composables/useRecaptcha.js';
import { useClipboard } from '../../composables/useClipboard.js';
import { API_BASE } from '../../lib/api.js';
import * as CryptoModule from '../../lib/crypto.js';

const { t } = useI18n();
const router = useRouter();
const upload = useUpload();
const { getToken } = useRecaptcha();
const { copied: passwordCopied, copy: copyToClipboard } = useClipboard();
const ttlRef = ref(null);
const text = ref('');

const usePassword = ref(false);
const password = ref('');
const showPassword = ref(false);
const multiDownload = ref(false);

function generatePassword() {
    const chars = 'ABCDEFGHJKMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789!@#$%';
    const arr = new Uint8Array(16);
    crypto.getRandomValues(arr);
    password.value = Array.from(arr).map(b => chars[b % chars.length]).join('');
    showPassword.value = true;
}

function copyPasswordToClipboard() {
    if (password.value) copyToClipboard(password.value);
}

async function handleUpload() {
    const trimmed = text.value.trim();
    if (!trimmed) return;
    if (usePassword.value && password.value.length < 4) {
        alert('Password must be at least 4 characters');
        return;
    }

    try {
        upload.uploading.value = true;
        const accessMode = multiDownload.value ? 'multi' : 'one_time';
        const ttl = ttlRef.value.getTtl();
        const encoder = new TextEncoder();
        const textBuffer = encoder.encode(trimmed);

        if (usePassword.value) {
            // Vault-style: PBKDF2 double encryption
            upload.updateProgress(10, t('upload.progress.generatingKey'));
            const dataKey = await CryptoModule.generateKey();

            upload.updateProgress(15, t('upload.progress.derivingPasswordKey'));
            const salt = CryptoModule.generateSalt();
            const passwordKey = await CryptoModule.deriveKeyFromPassword(password.value, salt);

            upload.updateProgress(20, t('upload.progress.securingKey'));
            const encryptedKeyData = await CryptoModule.encryptKey(dataKey, passwordKey);
            const encryptedKeyBase64 = CryptoModule.arrayToBase64(encryptedKeyData);
            const saltBase64 = CryptoModule.arrayToBase64(salt);

            upload.updateProgress(30, t('upload.progress.encryptingText'));
            const encryptedData = await CryptoModule.encrypt(textBuffer, dataKey);
            const encryptedTextBase64 = CryptoModule.arrayToBase64(encryptedData);

            upload.updateProgress(60, t('upload.progress.verifying'));
            const recaptchaToken = await getToken('upload');

            upload.updateProgress(70, t('upload.progress.initializing'));
            const response = await fetch(`${API_BASE}/upload/init`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    content_type: 'text',
                    encrypted_text: encryptedTextBase64,
                    ttl,
                    access_mode: accessMode,
                    salt: saltBase64,
                    encrypted_key: encryptedKeyBase64,
                    recaptcha_token: recaptchaToken,
                }),
            });
            if (!response.ok) {
                const error = await response.json().catch(() => ({}));
                throw new Error(error.error || 'Upload failed');
            }
            const data = await response.json();

            upload.updateProgress(100, t('upload.progress.completeRedirect'));
            setTimeout(() => {
                router.push({ path: '/share', query: { id: data.file_id, salt: saltBase64, vault: '1' } });
            }, 500);
        } else {
            // Standard: key in URL
            upload.updateProgress(10, t('upload.progress.generatingKey'));
            const key = await CryptoModule.generateKey();

            upload.updateProgress(30, t('upload.progress.encryptingText'));
            const encryptedData = await CryptoModule.encrypt(textBuffer, key);
            const base64Text = CryptoModule.arrayToBase64(encryptedData);
            const keyBase64 = await CryptoModule.keyToBase64(key);

            upload.updateProgress(60, t('upload.progress.verifying'));
            const recaptchaToken = await getToken('upload');

            upload.updateProgress(70, t('upload.progress.initializing'));
            const response = await fetch(`${API_BASE}/upload/init`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    content_type: 'text',
                    encrypted_text: base64Text,
                    ttl,
                    access_mode: accessMode,
                    recaptcha_token: recaptchaToken,
                }),
            });
            if (!response.ok) {
                const error = await response.json().catch(() => ({}));
                throw new Error(error.error || 'Upload failed');
            }
            const data = await response.json();

            upload.updateProgress(100, t('upload.progress.completeRedirect'));
            setTimeout(() => {
                router.push({ path: '/share', query: { id: data.file_id, key: keyBase64 } });
            }, 500);
        }
    } catch (error) {
        console.error('Upload error:', error);
        alert(t('upload.progress.uploadFailedMsg', { message: error.message }));
        upload.uploading.value = false;
    }
}
```

**Step 3: Check that `CryptoModule.arrayToBase64` exists**

```bash
grep -n "arrayToBase64\|base64ToArray" frontend/src/lib/crypto.js
```

If it doesn't exist (older name might be different), check what VaultUpload.vue used — it calls `CryptoModule.arrayToBase64`. Confirm the function name is correct before running.

**Step 4: Run frontend tests**

```bash
cd frontend && npx vitest run
```

**Step 5: Commit**

```bash
git add frontend/src/views/home/TextUpload.vue
git commit -m "feat: add password and multi-download options to text URL upload"
```

---

## Task 6: Remove Vault tab from `HomeView.vue`

**Files:**
- Modify: `frontend/src/views/HomeView.vue`

**Step 1: Remove the vault tab button** (lines ~46–51 in HomeView.vue):

```html
<!-- DELETE this entire button block -->
<button :class="['tab-btn', { active: linkTab === 'vault' }]" @click="linkTab = 'vault'">
  <svg class="inline w-4 h-4 mr-2 -mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
    <path d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
  </svg>
  {{ $t('home.tabs.vault') }}
</button>
```

**Step 2: Remove VaultUpload usage** (line ~57):

```html
<!-- DELETE -->
<VaultUpload v-if="linkTab === 'vault'" />
```

**Step 3: Remove imports** in `<script setup>`:

```js
// DELETE these two lines:
import VaultUpload from './home/VaultUpload.vue';
// and remove 'VaultUpload' from the import list
```

**Step 4: Run frontend tests**

```bash
cd frontend && npx vitest run
```

**Step 5: Commit**

```bash
git add frontend/src/views/HomeView.vue
git commit -m "feat: remove vault tab from HomeView"
```

---

## Task 7: Delete `VaultUpload.vue`

**Files:**
- Delete: `frontend/src/views/home/VaultUpload.vue`

**Step 1: Delete the file**

```bash
rm frontend/src/views/home/VaultUpload.vue
```

**Step 2: Run frontend tests to confirm nothing imports it**

```bash
cd frontend && npx vitest run
```

**Step 3: Commit**

```bash
git add -A
git commit -m "chore: delete VaultUpload.vue (merged into LinkUpload and TextUpload)"
```

---

## Task 8: Update `VaultDownload.vue` messaging for one-time password links

**Files:**
- Modify: `frontend/src/views/download/VaultDownload.vue`

Currently VaultDownload always shows "you can access this again with the same link" messaging. Now it can also serve password-protected single-download files, so messaging must reflect `access_mode`.

**Step 1: Add `accessMode` reactive variable** in `<script setup>`:

```js
const accessMode = ref('multi'); // default; updated from metadata
```

**Step 2: Set it from metadata** in `checkVaultAvailability()`:

```js
vaultMetadata = metadata;
accessMode.value = metadata.access_mode || 'multi';
```

**Step 3: Conditionally show "access again" note** — replace the hardcoded note in the password-entry section:

```html
<!-- Replace: -->
<div class="p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg mt-4">
  <p class="text-blue-700 dark:text-blue-300 text-sm text-center">
    {{ $t('download.vault.multiAccessNote') }}
  </p>
</div>

<!-- With: -->
<div v-if="accessMode === 'multi'" class="p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg mt-4">
  <p class="text-blue-700 dark:text-blue-300 text-sm text-center">
    {{ $t('download.vault.multiAccessNote') }}
  </p>
</div>
```

**Step 4: Conditionally show "access again" in text-result and file-success sections** — same pattern, wrap the "access again" / "download again" blue boxes with `v-if="accessMode === 'multi'"`.

The text-result "access again" block (line ~83–87):
```html
<div v-if="accessMode === 'multi'" class="p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg mb-6">
  <p class="text-blue-700 dark:text-blue-300 text-sm text-center">
    {{ $t('download.vault.accessAgain') }}
  </p>
</div>
```

The file-success "download again" block (line ~101–105):
```html
<div v-if="accessMode === 'multi'" class="p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg mb-6">
  <p class="text-blue-700 dark:text-blue-300 text-sm text-center">
    {{ $t('download.vault.downloadAgain') }}
  </p>
</div>
```

**Step 5: Run frontend tests**

```bash
cd frontend && npx vitest run
```

**Step 6: Commit**

```bash
git add frontend/src/views/download/VaultDownload.vue
git commit -m "feat: update VaultDownload messaging for one-time password links"
```

---

## Task 9: Update `ShareView.vue` messaging

**Files:**
- Modify: `frontend/src/views/ShareView.vue`

Currently `isVault` shows "stored in the vault" text. Since we've merged vault into URL mode, update these labels to say "password protected" instead.

**Step 1: Check current vault messages in `en.json`**

```bash
grep -n "vaultFileStored\|vaultTextStored\|passwordProtected\|multiAccessUntilExpiry" frontend/src/i18n/locales/en.json
```

**Step 2: Update i18n keys** in `en.json` and `uk.json`:

```json
"passwordFileStored": "Your file is encrypted and password-protected.",
"passwordTextStored": "Your text is encrypted and password-protected.",
"passwordProtectedNote": "Password protected",
"multiAccessUntilExpiry": "multiple downloads until expiry"
```

**Step 3: Update `ShareView.vue`** — replace the vault subtitle lines (~23–29):

```html
<p class="text-gray-600 dark:text-slate-400 mt-2" v-if="isVault">
  <span class="inline-flex items-center gap-2">
    <svg class="w-5 h-5 text-blue-600 dark:text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
      <path d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
    </svg>
    {{ fileName ? $t('share.passwordFileStored') : $t('share.passwordTextStored') }}
  </span>
</p>
```

And the warning section vault message (line ~83–86):
```html
<li v-if="isVault">
  <span class="text-blue-600 dark:text-blue-400 font-medium">{{ $t('share.passwordProtectedNote') }}</span>
  &bull; {{ $t('share.multiAccessUntilExpiry') }}
</li>
```

**Step 4: Run frontend tests**

```bash
cd frontend && npx vitest run
```

**Step 5: Commit**

```bash
git add frontend/src/views/ShareView.vue frontend/src/i18n/locales/en.json frontend/src/i18n/locales/uk.json
git commit -m "feat: update ShareView messaging for password-protected links"
```

---

## Task 10: Manual smoke test

With dev server running (`cd frontend && npm run dev`):

1. **File upload, no options** → upload a file → share page shows link → download works once → link dead.
2. **File upload + password** → toggle password ON → generate password → copy it → upload → share page (no key in URL, has salt+vault=1) → download page asks for password → enter wrong → error → enter correct → file downloads → link dead (one_time).
3. **File upload + multiple downloads** → toggle multi ON (TTL clamps to max 24h, Custom option hidden) → upload → share page → download works → download same link again → works again.
4. **File upload + password + multiple downloads** → both ON → upload → share page → download page asks password → enter correct → downloads → download again with same link + password → works again.
5. **Text upload** — repeat all 4 combos for text tab.
6. **Vault tab gone** — confirm the tab bar only shows File + Text.
7. **TtlSelector custom option** — with multi OFF: custom option visible. With multi ON: custom option hidden.

---

## Quick Reference

| Combo | `access_mode` | URL has | Download flow |
|---|---|---|---|
| No password + single | `one_time` | `key` | LinkDownload |
| No password + multi | `multi` | `key` | LinkDownload (multi access) |
| Password + single | `one_time` | `salt` + `vault=1` | VaultDownload (one-time) |
| Password + multi | `multi` | `salt` + `vault=1` | VaultDownload (multi access) |

> Note: `LinkDownload.vue` already handles `access_mode: multi` for the no-password multi case — it shows a "file available multiple times" message when `access_mode === 'multi'` in the download response. Verify this works or update it in a follow-up.