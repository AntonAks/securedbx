<template>
  <div>
    <!-- Loading -->
    <section v-if="step === 'loading'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-12">
      <div class="flex flex-col items-center">
        <div class="w-12 h-12 border-4 border-gray-300 dark:border-slate-700 border-t-blue-500 rounded-full animate-spin mb-4"></div>
        <p class="text-gray-600 dark:text-slate-300">Checking file availability...</p>
      </div>
    </section>

    <!-- Available -->
    <section v-if="step === 'available'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-6">
        {{ isText ? 'Secret Text Ready' : 'File Ready to Download' }}
      </h2>

      <div class="bg-gray-100 dark:bg-slate-800/50 rounded-lg p-4 mb-6 space-y-2">
        <p class="text-gray-600 dark:text-slate-300">
          <strong class="text-gray-900 dark:text-slate-100">{{ isText ? 'Text Size:' : 'File Size:' }}</strong>
          <span class="text-gray-500 dark:text-slate-400 ml-1">{{ formattedSize }}</span>
        </p>
        <p class="text-gray-600 dark:text-slate-300">
          <strong class="text-gray-900 dark:text-slate-100">Expires:</strong>
          <span class="text-gray-500 dark:text-slate-400 ml-1">{{ expiresText }}</span>
        </p>
      </div>

      <div class="bg-yellow-500/10 border border-yellow-500 rounded-lg p-4 mb-6">
        <p class="text-yellow-500 font-semibold mb-2">Warning:</p>
        <ul class="text-gray-700 dark:text-slate-300 text-sm space-y-1 ml-5 list-disc">
          <li>{{ isText ? 'This text can only be viewed ONCE' : 'This file can only be downloaded ONCE' }}</li>
          <li>After viewing, the link becomes invalid forever</li>
          <li>Make sure you're ready before proceeding</li>
        </ul>
      </div>

      <button class="btn-primary mb-6" :disabled="isDownloading" @click="handleDownload">
        {{ isText ? 'View Secret' : 'Download Now' }}
      </button>

      <ProgressBar :visible="isDownloading" :percent="progress" :text="progressText" />
    </section>

    <!-- Text Result -->
    <section v-if="step === 'text-result'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <div class="bg-green-500/10 border border-green-500 rounded-lg p-4 mb-6">
        <h2 class="text-green-500 font-semibold text-xl mb-2">Secret Text</h2>
        <p class="text-gray-500 dark:text-slate-400 text-sm">This text has been decrypted and can only be viewed once.</p>
      </div>
      <div class="mb-6">
        <textarea :value="decryptedText" readonly rows="10"
          class="w-full px-4 py-3 bg-gray-100 dark:bg-slate-800 border border-gray-300 dark:border-slate-700 rounded-lg text-gray-900 dark:text-slate-200 font-mono text-sm mb-4 resize-vertical focus:outline-none focus:border-blue-500"></textarea>
        <button @click="copyText" class="btn-secondary btn-auto">{{ textCopied ? 'Copied!' : 'Copy Text' }}</button>
      </div>
      <div class="bg-yellow-500/10 border border-yellow-500 rounded-lg p-4 mb-6">
        <p class="text-yellow-500 font-semibold mb-2">Important:</p>
        <ul class="text-gray-700 dark:text-slate-300 text-sm space-y-1 ml-5 list-disc">
          <li>This link is now <strong class="text-gray-900 dark:text-white">invalid forever</strong></li>
          <li>The text has been deleted from our servers</li>
          <li>Save the text now if you need it</li>
        </ul>
      </div>
      <router-link to="/" class="btn-primary block text-center">Share Another Secret</router-link>
    </section>

    <!-- File Success -->
    <section v-if="step === 'file-success'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <div class="bg-green-500/10 border border-green-500 rounded-lg p-4 mb-6">
        <h2 class="text-green-500 font-semibold text-xl mb-2">Download Complete</h2>
        <p class="text-gray-700 dark:text-slate-300 mb-2">Your file has been decrypted and saved.</p>
        <p class="text-gray-500 dark:text-slate-400 text-sm">This link is now invalid and the file has been deleted from our servers.</p>
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
            showError('Invalid download link');
            return;
        }

        encryptionKey = await CryptoModule.base64ToKey(keyBase64);
        await checkFileAvailability();
    } catch (error) {
        console.error('Initialization error:', error);
        showError('Failed to load file information');
    }
});

onUnmounted(() => {
    if (countdownInterval) clearInterval(countdownInterval);
});

async function checkFileAvailability() {
    try {
        const response = await fetch(`${API_BASE}/files/${fileId}/metadata`);

        if (response.status === 404) { showError('File not found'); return; }
        if (response.status === 410) { showError('File has expired or was already downloaded'); return; }
        if (!response.ok) throw new Error(`API error: ${response.status}`);

        const metadata = await response.json();
        if (!metadata.available) { showError('File has already been downloaded'); return; }

        isText.value = metadata.content_type === 'text';
        formattedSize.value = formatFileSize(metadata.file_size);
        startExpirationCountdown(metadata.expires_at);
        step.value = 'available';
    } catch (error) {
        console.error('Error checking file availability:', error);
        showError('Failed to check file availability');
    }
}

function startExpirationCountdown(expiresAt) {
    const update = () => {
        const remaining = expiresAt - getCurrentTimestamp();
        if (remaining <= 0) {
            expiresText.value = 'Expired';
            clearInterval(countdownInterval);
            showError('File has expired');
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
        progressText.value = 'Preparing download...';

        const recaptchaToken = await getToken('download');

        const response = await fetch(`${API_BASE}/files/${fileId}/download`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ recaptcha_token: recaptchaToken }),
        });

        if (response.status === 404) throw new Error('File not found');
        if (response.status === 410) throw new Error('File already downloaded or expired');
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
        showError('Download failed. The file may have already been downloaded or expired.');
    }
}

async function handleTextDownload(downloadData) {
    progress.value = 50;
    progressText.value = 'Decrypting text... 0%';

    const encryptedBytes = Uint8Array.from(atob(downloadData.encrypted_text), c => c.charCodeAt(0));
    const decryptedData = await CryptoModule.decryptFile(encryptedBytes, encryptionKey, (p) => {
        progress.value = 50 + p * 0.4;
        progressText.value = `Decrypting text... ${Math.round(p)}%`;
    });

    progress.value = 100;
    progressText.value = 'Complete!';
    decryptedText.value = new TextDecoder().decode(decryptedData);
    step.value = 'text-result';
}

async function handleFileDownload(downloadData) {
    progress.value = 30;
    progressText.value = 'Downloading encrypted file...';

    const encryptedData = await downloadFromS3(downloadData.download_url);

    progress.value = 70;
    progressText.value = 'Decrypting file... 0%';
    const decryptedData = await CryptoModule.decryptFile(encryptedData, encryptionKey, (p) => {
        progress.value = 70 + p * 0.2;
        progressText.value = `Decrypting file... ${Math.round(p)}%`;
    });

    progress.value = 90;
    progressText.value = 'Saving file...';
    saveFile(decryptedData, fileName || 'downloaded-file');

    progress.value = 100;
    progressText.value = 'Download complete!';
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
                progressText.value = `Downloading... ${pct.toFixed(1)}%`;
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
