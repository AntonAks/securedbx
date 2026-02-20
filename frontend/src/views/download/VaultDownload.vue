<template>
  <div>
    <!-- Loading -->
    <section v-if="step === 'loading'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-12">
      <div class="flex flex-col items-center">
        <div class="w-12 h-12 border-4 border-gray-300 dark:border-slate-700 border-t-blue-500 rounded-full animate-spin mb-4"></div>
        <p class="text-gray-600 dark:text-slate-300">{{ $t('download.vault.checkingAvailability') }}</p>
      </div>
    </section>

    <!-- Password Entry -->
    <section v-if="step === 'password'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <div class="text-center mb-6">
        <div class="inline-flex items-center justify-center w-16 h-16 bg-blue-100 dark:bg-blue-900/30 rounded-full mb-4">
          <svg class="w-8 h-8 text-blue-600 dark:text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
          </svg>
        </div>
        <h2 class="text-xl font-semibold text-gray-900 dark:text-slate-100 mb-2">{{ $t('download.vault.passwordRequired') }}</h2>
        <p class="text-gray-600 dark:text-slate-400">{{ $t('download.vault.enterPassword') }}</p>
      </div>

      <div class="bg-gray-100 dark:bg-slate-800/50 rounded-lg p-4 mb-6 space-y-2">
        <p class="text-gray-600 dark:text-slate-300">
          <strong class="text-gray-900 dark:text-slate-100">{{ $t('download.vault.size') }}</strong>
          <span class="text-gray-500 dark:text-slate-400 ml-1">{{ formattedSize }}</span>
        </p>
        <p class="text-gray-600 dark:text-slate-300">
          <strong class="text-gray-900 dark:text-slate-100">{{ $t('download.vault.expires') }}</strong>
          <span class="text-gray-500 dark:text-slate-400 ml-1">{{ expiresText }}</span>
        </p>
        <p class="text-gray-600 dark:text-slate-300">
          <strong class="text-gray-900 dark:text-slate-100">{{ $t('download.vault.downloads') }}</strong>
          <span class="text-gray-500 dark:text-slate-400 ml-1">{{ downloadCount }}</span>
        </p>
      </div>

      <div class="mb-6">
        <div class="flex items-center gap-2">
          <input
            :type="showPassword ? 'text' : 'password'"
            v-model="password"
            :placeholder="$t('download.vault.passwordPlaceholder')"
            class="flex-1 px-4 py-3 bg-white dark:bg-slate-800 border border-gray-300 dark:border-slate-700 rounded-lg text-gray-900 dark:text-slate-200 focus:outline-none focus:border-blue-500"
            @keypress.enter="handleVaultDownload"
          />
          <button type="button" @click="showPassword = !showPassword"
            class="px-3 py-3 text-gray-500 hover:text-gray-700 dark:text-slate-400 dark:hover:text-slate-200 border border-gray-300 dark:border-slate-700 rounded-lg">
            {{ showPassword ? $t('download.vault.hide') : $t('download.vault.show') }}
          </button>
        </div>
      </div>

      <button class="btn-primary mb-4" :disabled="!password || isDownloading" @click="handleVaultDownload">
        {{ $t('download.vault.unlockDownload') }}
      </button>

      <ProgressBar :visible="isDownloading" :percent="progress" :text="progressText" />

      <div class="p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg mt-4">
        <p class="text-blue-700 dark:text-blue-300 text-sm text-center">
          {{ $t('download.vault.multiAccessNote') }}
        </p>
      </div>
    </section>

    <!-- Vault Text Result -->
    <section v-if="step === 'text-result'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <div class="bg-green-500/10 border border-green-500 rounded-lg p-4 mb-6">
        <h2 class="text-green-500 font-semibold text-xl mb-2">
          <svg class="w-5 h-5 inline mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
          </svg>
          {{ $t('download.vault.vaultText') }}
        </h2>
        <p class="text-gray-500 dark:text-slate-400 text-sm">{{ $t('download.vault.vaultTextNote') }}</p>
      </div>
      <div class="mb-6">
        <textarea :value="decryptedText" readonly rows="10"
          class="w-full px-4 py-3 bg-gray-100 dark:bg-slate-800 border border-gray-300 dark:border-slate-700 rounded-lg text-gray-900 dark:text-slate-200 font-mono text-sm mb-4 resize-vertical focus:outline-none focus:border-blue-500"></textarea>
        <button @click="copyText" class="btn-secondary btn-auto">{{ textCopied ? $t('download.vault.textCopied') : $t('download.vault.copyText') }}</button>
      </div>
      <div class="p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg mb-6">
        <p class="text-blue-700 dark:text-blue-300 text-sm text-center">
          {{ $t('download.vault.accessAgain') }}
        </p>
      </div>
      <router-link to="/" class="btn-primary block text-center">{{ $t('download.vault.shareAnotherSecret') }}</router-link>
    </section>

    <!-- Vault File Success -->
    <section v-if="step === 'file-success'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <div class="bg-green-500/10 border border-green-500 rounded-lg p-4 mb-6">
        <h2 class="text-green-500 font-semibold text-xl mb-2">
          <svg class="w-5 h-5 inline mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
          </svg>
          {{ $t('download.vault.downloadComplete') }}
        </h2>
        <p class="text-gray-700 dark:text-slate-300 mb-2">{{ $t('download.vault.fileDecrypted') }}</p>
      </div>
      <div class="p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg mb-6">
        <p class="text-blue-700 dark:text-blue-300 text-sm text-center">
          {{ $t('download.vault.downloadAgain') }}
        </p>
      </div>
      <router-link to="/" class="btn-primary block text-center">{{ $t('download.vault.shareAnotherFile') }}</router-link>
    </section>

    <!-- Error -->
    <section v-if="step === 'error'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <div class="bg-red-500/10 border border-red-500 rounded-lg p-4 mb-6">
        <h2 class="text-red-500 font-semibold text-xl mb-2">{{ $t('download.vault.unableToDownload') }}</h2>
        <p class="text-gray-700 dark:text-slate-300">{{ errorMessage }}</p>
      </div>
      <router-link to="/" class="btn-outline block text-center">{{ $t('download.vault.shareNewFile') }}</router-link>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';
import ProgressBar from '../../components/ProgressBar.vue';
import { useRecaptcha } from '../../composables/useRecaptcha.js';
import { useClipboard } from '../../composables/useClipboard.js';
import { API_BASE } from '../../lib/api.js';
import { formatFileSize, formatTimeRemaining, getCurrentTimestamp } from '../../lib/utils.js';
import * as CryptoModule from '../../lib/crypto.js';

const { t } = useI18n();
const route = useRoute();
const { getToken } = useRecaptcha();
const { copied: textCopied, copy } = useClipboard();

const step = ref('loading');
const password = ref('');
const showPassword = ref(false);
const isDownloading = ref(false);
const progress = ref(0);
const progressText = ref('');
const formattedSize = ref('');
const expiresText = ref('');
const downloadCount = ref(0);
const decryptedText = ref('');
const errorMessage = ref('');

let fileId = '';
let vaultSalt = '';
let fileName = '';
let vaultMetadata = null;
let countdownInterval = null;

onMounted(async () => {
    try {
        fileId = route.query.id;
        vaultSalt = route.query.salt;
        fileName = route.query.name ? decodeURIComponent(route.query.name) : '';

        if (!fileId || !vaultSalt) {
            showError(t('download.vault.invalidVaultLink'));
            return;
        }

        await checkVaultAvailability();
    } catch (error) {
        console.error('Vault initialization error:', error);
        showError(t('download.vault.failedToLoad'));
    }
});

onUnmounted(() => {
    if (countdownInterval) clearInterval(countdownInterval);
});

async function checkVaultAvailability() {
    try {
        const response = await fetch(`${API_BASE}/files/${fileId}/metadata`);

        if (response.status === 404) { showError(t('download.vault.vaultNotFound')); return; }
        if (response.status === 410) { showError(t('download.vault.vaultExpired')); return; }
        if (!response.ok) throw new Error(`API error: ${response.status}`);

        const metadata = await response.json();
        if (!metadata.available) { showError(t('download.vault.vaultUnavailable')); return; }

        vaultMetadata = metadata;
        formattedSize.value = formatFileSize(metadata.file_size);
        downloadCount.value = metadata.download_count || 0;
        startExpirationCountdown(metadata.expires_at);
        step.value = 'password';
    } catch (error) {
        console.error('Error checking vault availability:', error);
        showError(t('download.vault.failedToCheck'));
    }
}

function startExpirationCountdown(expiresAt) {
    const update = () => {
        const remaining = expiresAt - getCurrentTimestamp();
        if (remaining <= 0) {
            expiresText.value = t('download.vault.vaultExpired');
            clearInterval(countdownInterval);
            showError(t('download.vault.vaultExpired'));
            return;
        }
        expiresText.value = formatTimeRemaining(remaining);
    };
    update();
    countdownInterval = setInterval(update, 60000);
}

async function handleVaultDownload() {
    if (!password.value) {
        alert(t('download.vault.enterPasswordAlert'));
        return;
    }

    try {
        isDownloading.value = true;

        progress.value = 5;
        progressText.value = t('download.progress.verifyingPassword');
        const salt = CryptoModule.base64ToArray(vaultSalt);
        const passwordKey = await CryptoModule.deriveKeyFromPassword(password.value, salt);

        progress.value = 15;
        progressText.value = t('download.progress.unlockingContent');
        const encryptedKeyData = CryptoModule.base64ToArray(vaultMetadata.encrypted_key);

        let dataKey;
        try {
            dataKey = await CryptoModule.decryptKey(encryptedKeyData, passwordKey);
        } catch (error) {
            console.error('Key decryption failed:', error);
            alert(t('download.vault.incorrectPassword'));
            resetForm();
            return;
        }

        progress.value = 25;
        progressText.value = t('download.progress.preparingDownload');

        const recaptchaToken = await getToken('download');
        const response = await fetch(`${API_BASE}/files/${fileId}/download`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ recaptcha_token: recaptchaToken }),
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.error || `Download failed: ${response.status}`);
        }

        const downloadData = await response.json();

        if (downloadData.content_type === 'text') {
            progress.value = 40;
            progressText.value = t('download.progress.decryptingText');
            const encryptedBytes = CryptoModule.base64ToArray(downloadData.encrypted_text);

            const decryptedData = await CryptoModule.decryptFile(encryptedBytes, dataKey, (p) => {
                progress.value = 40 + p * 0.5;
                progressText.value = t('download.progress.decrypting', { percent: Math.round(p) });
            });

            progress.value = 100;
            progressText.value = t('download.progress.complete');
            decryptedText.value = new TextDecoder().decode(decryptedData);
            step.value = 'text-result';
        } else {
            progress.value = 30;
            progressText.value = t('download.progress.downloadingFile');
            const encryptedData = await downloadFromS3(downloadData.download_url);

            progress.value = 60;
            progressText.value = t('download.progress.decryptingFile', { percent: 0 });
            const decryptedData = await CryptoModule.decryptFile(encryptedData, dataKey, (p) => {
                progress.value = 60 + p * 0.35;
                progressText.value = t('download.progress.decrypting', { percent: Math.round(p) });
            });

            progress.value = 95;
            progressText.value = t('download.progress.savingFile');
            saveFile(decryptedData, fileName || 'downloaded-file');

            progress.value = 100;
            progressText.value = t('download.progress.downloadComplete');
            step.value = 'file-success';
        }
    } catch (error) {
        console.error('Vault download error:', error);
        if (error.name === 'OperationError') {
            alert(t('download.vault.incorrectOrCorrupt'));
        } else {
            alert(error.message || t('download.link.downloadFailed'));
        }
        resetForm();
    }
}

function downloadFromS3(url) {
    return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        xhr.responseType = 'arraybuffer';

        xhr.onprogress = (event) => {
            if (event.lengthComputable) {
                const pct = (event.loaded / event.total) * 100;
                progress.value = 30 + pct * 0.3;
                progressText.value = t('download.progress.downloading', { percent: pct.toFixed(1) });
            }
        };

        xhr.onload = () => {
            if (xhr.status >= 200 && xhr.status < 300) resolve(new Uint8Array(xhr.response));
            else reject(new Error(`Download failed: ${xhr.status}`));
        };

        xhr.onerror = () => reject(new Error('Network error during download'));
        xhr.ontimeout = () => reject(new Error('Download timeout'));

        xhr.open('GET', url);
        xhr.send();
    });
}

function saveFile(data, filename) {
    const blob = new Blob([data]);
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}

function resetForm() {
    isDownloading.value = false;
    progress.value = 0;
    progressText.value = '';
}

function copyText() {
    copy(decryptedText.value);
}

function showError(msg) {
    errorMessage.value = msg;
    step.value = 'error';
}
</script>
