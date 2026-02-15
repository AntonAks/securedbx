<template>
  <div>
    <!-- Loading -->
    <section v-if="step === 'loading'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-12">
      <div class="flex flex-col items-center">
        <div class="w-12 h-12 border-4 border-gray-300 dark:border-slate-700 border-t-blue-500 rounded-full animate-spin mb-4"></div>
        <p class="text-gray-600 dark:text-slate-300">Checking vault availability...</p>
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
        <h2 class="text-xl font-semibold text-gray-900 dark:text-slate-100 mb-2">Password Required</h2>
        <p class="text-gray-600 dark:text-slate-400">Enter the password to access this content</p>
      </div>

      <div class="bg-gray-100 dark:bg-slate-800/50 rounded-lg p-4 mb-6 space-y-2">
        <p class="text-gray-600 dark:text-slate-300">
          <strong class="text-gray-900 dark:text-slate-100">Size:</strong>
          <span class="text-gray-500 dark:text-slate-400 ml-1">{{ formattedSize }}</span>
        </p>
        <p class="text-gray-600 dark:text-slate-300">
          <strong class="text-gray-900 dark:text-slate-100">Expires:</strong>
          <span class="text-gray-500 dark:text-slate-400 ml-1">{{ expiresText }}</span>
        </p>
        <p class="text-gray-600 dark:text-slate-300">
          <strong class="text-gray-900 dark:text-slate-100">Downloads:</strong>
          <span class="text-gray-500 dark:text-slate-400 ml-1">{{ downloadCount }}</span>
        </p>
      </div>

      <div class="mb-6">
        <div class="flex items-center gap-2">
          <input
            :type="showPassword ? 'text' : 'password'"
            v-model="password"
            placeholder="Enter password"
            class="flex-1 px-4 py-3 bg-white dark:bg-slate-800 border border-gray-300 dark:border-slate-700 rounded-lg text-gray-900 dark:text-slate-200 focus:outline-none focus:border-blue-500"
            @keypress.enter="handleVaultDownload"
          />
          <button type="button" @click="showPassword = !showPassword"
            class="px-3 py-3 text-gray-500 hover:text-gray-700 dark:text-slate-400 dark:hover:text-slate-200 border border-gray-300 dark:border-slate-700 rounded-lg">
            {{ showPassword ? 'Hide' : 'Show' }}
          </button>
        </div>
      </div>

      <button class="btn-primary mb-4" :disabled="!password || isDownloading" @click="handleVaultDownload">
        Unlock &amp; Download
      </button>

      <ProgressBar :visible="isDownloading" :percent="progress" :text="progressText" />

      <div class="p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg mt-4">
        <p class="text-blue-700 dark:text-blue-300 text-sm text-center">
          You can access this multiple times until it expires
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
          Vault Text
        </h2>
        <p class="text-gray-500 dark:text-slate-400 text-sm">This password-protected text can be accessed multiple times until it expires.</p>
      </div>
      <div class="mb-6">
        <textarea :value="decryptedText" readonly rows="10"
          class="w-full px-4 py-3 bg-gray-100 dark:bg-slate-800 border border-gray-300 dark:border-slate-700 rounded-lg text-gray-900 dark:text-slate-200 font-mono text-sm mb-4 resize-vertical focus:outline-none focus:border-blue-500"></textarea>
        <button @click="copyText" class="btn-secondary btn-auto">{{ textCopied ? 'Copied!' : 'Copy Text' }}</button>
      </div>
      <div class="p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg mb-6">
        <p class="text-blue-700 dark:text-blue-300 text-sm text-center">
          You can access this again with the same link and password
        </p>
      </div>
      <router-link to="/" class="btn-primary block text-center">Share Another Secret</router-link>
    </section>

    <!-- Vault File Success -->
    <section v-if="step === 'file-success'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <div class="bg-green-500/10 border border-green-500 rounded-lg p-4 mb-6">
        <h2 class="text-green-500 font-semibold text-xl mb-2">
          <svg class="w-5 h-5 inline mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
          </svg>
          Download Complete
        </h2>
        <p class="text-gray-700 dark:text-slate-300 mb-2">Your file has been decrypted and saved.</p>
      </div>
      <div class="p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg mb-6">
        <p class="text-blue-700 dark:text-blue-300 text-sm text-center">
          You can download this file again with the same link and password until it expires
        </p>
      </div>
      <router-link to="/" class="btn-primary block text-center">Share Another File</router-link>
    </section>

    <!-- Error -->
    <section v-if="step === 'error'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <div class="bg-red-500/10 border border-red-500 rounded-lg p-4 mb-6">
        <h2 class="text-red-500 font-semibold text-xl mb-2">Unable to Download</h2>
        <p class="text-gray-700 dark:text-slate-300">{{ errorMessage }}</p>
      </div>
      <router-link to="/" class="btn-outline block text-center">Share a New File</router-link>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';
import ProgressBar from '../../components/ProgressBar.vue';
import { useRecaptcha } from '../../composables/useRecaptcha.js';
import { useClipboard } from '../../composables/useClipboard.js';
import { API_BASE } from '../../lib/api.js';
import { formatFileSize, formatTimeRemaining, getCurrentTimestamp } from '../../lib/utils.js';
import * as CryptoModule from '../../lib/crypto.js';

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
            showError('Invalid vault link');
            return;
        }

        await checkVaultAvailability();
    } catch (error) {
        console.error('Vault initialization error:', error);
        showError('Failed to load vault information');
    }
});

onUnmounted(() => {
    if (countdownInterval) clearInterval(countdownInterval);
});

async function checkVaultAvailability() {
    try {
        const response = await fetch(`${API_BASE}/files/${fileId}/metadata`);

        if (response.status === 404) { showError('Vault content not found'); return; }
        if (response.status === 410) { showError('Vault content has expired'); return; }
        if (!response.ok) throw new Error(`API error: ${response.status}`);

        const metadata = await response.json();
        if (!metadata.available) { showError('Vault content is no longer available'); return; }

        vaultMetadata = metadata;
        formattedSize.value = formatFileSize(metadata.file_size);
        downloadCount.value = metadata.download_count || 0;
        startExpirationCountdown(metadata.expires_at);
        step.value = 'password';
    } catch (error) {
        console.error('Error checking vault availability:', error);
        showError('Failed to check vault availability');
    }
}

function startExpirationCountdown(expiresAt) {
    const update = () => {
        const remaining = expiresAt - getCurrentTimestamp();
        if (remaining <= 0) {
            expiresText.value = 'Expired';
            clearInterval(countdownInterval);
            showError('Vault content has expired');
            return;
        }
        expiresText.value = formatTimeRemaining(remaining);
    };
    update();
    countdownInterval = setInterval(update, 60000);
}

async function handleVaultDownload() {
    if (!password.value) {
        alert('Please enter the password');
        return;
    }

    try {
        isDownloading.value = true;

        progress.value = 5;
        progressText.value = 'Verifying password...';
        const salt = CryptoModule.base64ToArray(vaultSalt);
        const passwordKey = await CryptoModule.deriveKeyFromPassword(password.value, salt);

        progress.value = 15;
        progressText.value = 'Unlocking content...';
        const encryptedKeyData = CryptoModule.base64ToArray(vaultMetadata.encrypted_key);

        let dataKey;
        try {
            dataKey = await CryptoModule.decryptKey(encryptedKeyData, passwordKey);
        } catch (error) {
            console.error('Key decryption failed:', error);
            alert('Incorrect password');
            resetForm();
            return;
        }

        progress.value = 25;
        progressText.value = 'Preparing download...';

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
            progressText.value = 'Decrypting text...';
            const encryptedBytes = CryptoModule.base64ToArray(downloadData.encrypted_text);

            const decryptedData = await CryptoModule.decryptFile(encryptedBytes, dataKey, (p) => {
                progress.value = 40 + p * 0.5;
                progressText.value = `Decrypting... ${Math.round(p)}%`;
            });

            progress.value = 100;
            progressText.value = 'Complete!';
            decryptedText.value = new TextDecoder().decode(decryptedData);
            step.value = 'text-result';
        } else {
            progress.value = 30;
            progressText.value = 'Downloading encrypted file...';
            const encryptedData = await downloadFromS3(downloadData.download_url);

            progress.value = 60;
            progressText.value = 'Decrypting file...';
            const decryptedData = await CryptoModule.decryptFile(encryptedData, dataKey, (p) => {
                progress.value = 60 + p * 0.35;
                progressText.value = `Decrypting... ${Math.round(p)}%`;
            });

            progress.value = 95;
            progressText.value = 'Saving file...';
            saveFile(decryptedData, fileName || 'downloaded-file');

            progress.value = 100;
            progressText.value = 'Download complete!';
            step.value = 'file-success';
        }
    } catch (error) {
        console.error('Vault download error:', error);
        if (error.name === 'OperationError') {
            alert('Incorrect password or corrupted data');
        } else {
            alert(error.message || 'Download failed. Please try again.');
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
                progressText.value = `Downloading... ${pct.toFixed(1)}%`;
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
