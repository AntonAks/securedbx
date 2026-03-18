<template>
  <div>
    <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-6">{{ $t('upload.link.shareFiles') }}</h2>

    <DropZone v-if="selectedFiles.length === 0" @files="addFiles" />
    <FileList :files="selectedFiles" @remove="removeFile" @clear="clearFiles" />

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

    <TtlSelector ref="ttlRef" :limitTo24h="multiDownload" />

    <button class="btn-primary" :disabled="selectedFiles.length === 0 || upload.uploading.value" @click="handleUpload">
      {{ $t('upload.link.encryptUpload') }}
    </button>

    <ProgressBar :visible="upload.uploading.value" :percent="upload.progress.value" :text="upload.progressText.value" />
  </div>
</template>

<script setup>
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
            const saltBase64 = CryptoModule.arrayToBase64Url(salt);

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
            const rawKey = await CryptoModule.exportKey(key);
            const keyBase64 = CryptoModule.arrayToBase64Url(new Uint8Array(rawKey));
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
</script>
