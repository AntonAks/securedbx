<template>
  <div>
    <!-- Loading -->
    <section v-if="step === 'loading'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-12">
      <div class="flex flex-col items-center">
        <div class="w-12 h-12 border-4 border-gray-300 dark:border-slate-700 border-t-blue-500 rounded-full animate-spin mb-4"></div>
        <p class="text-gray-600 dark:text-slate-300">{{ $t('download.link.checkingAvailability') }}</p>
      </div>
    </section>

    <!-- Available -->
    <section v-if="step === 'available'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-6">
        {{ isText ? $t('download.link.secretTextReady') : $t('download.link.fileReady') }}
      </h2>

      <div class="bg-gray-100 dark:bg-slate-800/50 rounded-lg p-4 mb-6 space-y-2">
        <p class="text-gray-600 dark:text-slate-300">
          <strong class="text-gray-900 dark:text-slate-100">{{ isText ? $t('download.link.textSize') : $t('download.link.fileSize') }}</strong>
          <span class="text-gray-500 dark:text-slate-400 ml-1">{{ formattedSize }}</span>
        </p>
        <p class="text-gray-600 dark:text-slate-300">
          <strong class="text-gray-900 dark:text-slate-100">{{ $t('download.link.expires') }}</strong>
          <span class="text-gray-500 dark:text-slate-400 ml-1">{{ expiresText }}</span>
        </p>
      </div>

      <div class="bg-yellow-500/10 border border-yellow-500 rounded-lg p-4 mb-6">
        <p class="text-yellow-500 font-semibold mb-2">{{ $t('download.link.warning') }}</p>
        <ul class="text-gray-700 dark:text-slate-300 text-sm space-y-1 ml-5 list-disc">
          <li>{{ isText ? $t('download.link.textViewOnce') : $t('download.link.fileDownloadOnce') }}</li>
          <li>{{ $t('download.link.linkInvalidAfter') }}</li>
          <li>{{ $t('download.link.makeReady') }}</li>
        </ul>
      </div>

      <button class="btn-primary mb-6" :disabled="isDownloading" @click="handleDownload">
        {{ isText ? $t('download.link.viewSecret') : $t('download.link.downloadNow') }}
      </button>

      <ProgressBar :visible="isDownloading" :percent="progress" :text="progressText" />
    </section>

    <!-- Text Result -->
    <section v-if="step === 'text-result'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <div class="bg-green-500/10 border border-green-500 rounded-lg p-4 mb-6">
        <h2 class="text-green-500 font-semibold text-xl mb-2">{{ $t('download.link.secretText') }}</h2>
        <p class="text-gray-500 dark:text-slate-400 text-sm">{{ $t('download.link.textDecryptedOnce') }}</p>
      </div>
      <div class="mb-6">
        <textarea :value="decryptedText" readonly rows="10"
          class="w-full px-4 py-3 bg-gray-100 dark:bg-slate-800 border border-gray-300 dark:border-slate-700 rounded-lg text-gray-900 dark:text-slate-200 font-mono text-sm mb-4 resize-vertical focus:outline-none focus:border-blue-500"></textarea>
        <button @click="copyText" class="btn-secondary btn-auto">{{ textCopied ? $t('download.link.textCopied') : $t('download.link.copyText') }}</button>
      </div>
      <div class="bg-yellow-500/10 border border-yellow-500 rounded-lg p-4 mb-6">
        <p class="text-yellow-500 font-semibold mb-2">{{ $t('download.link.important') }}</p>
        <ul class="text-gray-700 dark:text-slate-300 text-sm space-y-1 ml-5 list-disc">
          <li>{{ $t('download.link.linkInvalidForever') }}</li>
          <li>{{ $t('download.link.textDeleted') }}</li>
          <li>{{ $t('download.link.saveTextNow') }}</li>
        </ul>
      </div>
      <router-link to="/" class="btn-primary block text-center">{{ $t('download.link.shareAnotherSecret') }}</router-link>
    </section>

    <!-- File Success -->
    <section v-if="step === 'file-success'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <div class="bg-green-500/10 border border-green-500 rounded-lg p-4 mb-6">
        <h2 class="text-green-500 font-semibold text-xl mb-2">{{ $t('download.link.downloadComplete') }}</h2>
        <p class="text-gray-700 dark:text-slate-300 mb-2">{{ $t('download.link.fileDecrypted') }}</p>
        <p class="text-gray-500 dark:text-slate-400 text-sm">{{ $t('download.link.fileDeletedFromServers') }}</p>
      </div>
      <router-link to="/" class="btn-primary block text-center">{{ $t('download.link.shareAnotherFile') }}</router-link>
    </section>

    <!-- Error -->
    <section v-if="step === 'error'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <div class="bg-red-500/10 border border-red-500 rounded-lg p-4 mb-6">
        <h2 class="text-red-500 font-semibold text-xl mb-2">{{ $t('download.link.unableToDownload') }}</h2>
        <p class="text-gray-700 dark:text-slate-300">{{ errorMessage }}</p>
      </div>
      <router-link to="/" class="btn-outline block text-center">{{ $t('download.link.shareNewFile') }}</router-link>
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
const isText = ref(false);
const isDownloading = ref(false);
const progress = ref(0);
const progressText = ref('');
const formattedSize = ref('');
const expiresText = ref('');
const decryptedText = ref('');
const errorMessage = ref('');

let fileId = '';
let encryptionKey = null;
let fileName = '';
let countdownInterval = null;

onMounted(async () => {
    try {
        fileId = route.query.id;
        const keyBase64 = route.query.key;
        fileName = route.query.name ? decodeURIComponent(route.query.name) : '';

        if (!fileId || !keyBase64) {
            showError(t('download.link.invalidLink'));
            return;
        }

        encryptionKey = await CryptoModule.base64ToKey(keyBase64);
        await checkFileAvailability();
    } catch (error) {
        console.error('Initialization error:', error);
        showError(t('download.link.failedToLoad'));
    }
});

onUnmounted(() => {
    if (countdownInterval) clearInterval(countdownInterval);
});

async function checkFileAvailability() {
    try {
        const response = await fetch(`${API_BASE}/files/${fileId}/metadata`);

        if (response.status === 404) { showError(t('download.link.fileNotFound')); return; }
        if (response.status === 410) { showError(t('download.link.fileExpired')); return; }
        if (!response.ok) throw new Error(`API error: ${response.status}`);

        const metadata = await response.json();
        if (!metadata.available) { showError(t('download.link.fileAlreadyDownloaded')); return; }

        isText.value = metadata.content_type === 'text';
        formattedSize.value = formatFileSize(metadata.file_size);
        startExpirationCountdown(metadata.expires_at);
        step.value = 'available';
    } catch (error) {
        console.error('Error checking file availability:', error);
        showError(t('download.link.failedToCheck'));
    }
}

function startExpirationCountdown(expiresAt) {
    const update = () => {
        const remaining = expiresAt - getCurrentTimestamp();
        if (remaining <= 0) {
            expiresText.value = t('download.link.expired');
            clearInterval(countdownInterval);
            showError(t('download.link.fileExpired'));
            return;
        }
        expiresText.value = formatTimeRemaining(remaining);
    };
    update();
    countdownInterval = setInterval(update, 60000);
}

async function handleDownload() {
    try {
        isDownloading.value = true;
        progress.value = 0;
        progressText.value = t('download.progress.preparingDownload');

        const recaptchaToken = await getToken('download');

        const response = await fetch(`${API_BASE}/files/${fileId}/download`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ recaptcha_token: recaptchaToken }),
        });

        if (response.status === 404) throw new Error(t('download.link.fileNotFound'));
        if (response.status === 410) throw new Error(t('download.link.fileExpired'));
        if (!response.ok) throw new Error(`Download request failed: ${response.status}`);

        const downloadData = await response.json();

        if (downloadData.content_type === 'text') {
            await handleTextDownload(downloadData);
        } else {
            await handleFileDownload(downloadData);
        }

        confirmDownload();
    } catch (error) {
        console.error('Download error:', error);
        showError(t('download.link.downloadFailed'));
    }
}

async function handleTextDownload(downloadData) {
    progress.value = 50;
    progressText.value = t('download.progress.decryptingTextProgress', { percent: 0 });

    const encryptedBytes = Uint8Array.from(atob(downloadData.encrypted_text), c => c.charCodeAt(0));
    const decryptedData = await CryptoModule.decryptFile(encryptedBytes, encryptionKey, (p) => {
        progress.value = 50 + p * 0.4;
        progressText.value = t('download.progress.decryptingTextProgress', { percent: Math.round(p) });
    });

    progress.value = 100;
    progressText.value = t('download.progress.complete');
    decryptedText.value = new TextDecoder().decode(decryptedData);
    step.value = 'text-result';
}

async function handleFileDownload(downloadData) {
    progress.value = 30;
    progressText.value = t('download.progress.downloadingFile');

    const encryptedData = await downloadFromS3(downloadData.download_url);

    progress.value = 70;
    progressText.value = t('download.progress.decryptingFile', { percent: 0 });
    const decryptedData = await CryptoModule.decryptFile(encryptedData, encryptionKey, (p) => {
        progress.value = 70 + p * 0.2;
        progressText.value = t('download.progress.decryptingFile', { percent: Math.round(p) });
    });

    progress.value = 90;
    progressText.value = t('download.progress.savingFile');
    saveFile(decryptedData, fileName || 'downloaded-file');

    progress.value = 100;
    progressText.value = t('download.progress.downloadComplete');
    step.value = 'file-success';
}

function downloadFromS3(url) {
    return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        xhr.responseType = 'arraybuffer';

        xhr.onprogress = (event) => {
            if (event.lengthComputable) {
                const pct = (event.loaded / event.total) * 100;
                progress.value = 30 + pct * 0.4;
                progressText.value = t('download.progress.downloading', { percent: pct.toFixed(1) });
            }
        };

        xhr.onload = () => {
            if (xhr.status >= 200 && xhr.status < 300) resolve(new Uint8Array(xhr.response));
            else reject(new Error(`File download failed: ${xhr.status}`));
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

function confirmDownload() {
    fetch(`${API_BASE}/files/${fileId}/confirm`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
    }).catch(err => console.warn('Confirm download error:', err));
}

function copyText() {
    copy(decryptedText.value);
}

function showError(msg) {
    errorMessage.value = msg;
    step.value = 'error';
}
</script>
