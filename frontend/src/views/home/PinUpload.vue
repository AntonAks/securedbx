<template>
  <div>
    <!-- File / Text Tab Navigation -->
    <div class="flex justify-center mb-8">
      <div class="inline-flex rounded-full bg-gray-200 dark:bg-slate-800 p-1">
        <button :class="['tab-btn', { active: activeTab === 'file' }]" @click="activeTab = 'file'">
          <svg class="inline w-4 h-4 mr-2 -mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
          </svg>
          File
        </button>
        <button :class="['tab-btn', { active: activeTab === 'text' }]" @click="activeTab = 'text'">
          <svg class="inline w-4 h-4 mr-2 -mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
          </svg>
          Text
        </button>
      </div>
    </div>

    <!-- File Tab -->
    <div v-if="activeTab === 'file'">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-6">Share Files</h2>
      <DropZone v-if="selectedFiles.length === 0" @files="setFiles" />

      <!-- Single file info -->
      <div v-if="selectedFiles.length === 1" class="bg-gray-100 dark:bg-slate-800/50 rounded-lg p-4 mb-6 flex items-center justify-between">
        <div>
          <span class="text-gray-900 dark:text-slate-100 font-medium">{{ selectedFiles[0].name }}</span>
          <span class="text-gray-500 dark:text-slate-400 text-sm ml-2">{{ formatFileSize(selectedFiles[0].size) }}</span>
        </div>
        <button @click="clearFiles" class="text-red-500 hover:text-red-600 dark:text-red-400 dark:hover:text-red-300 text-sm">&times; Remove</button>
      </div>

      <!-- Multi file list -->
      <div v-if="selectedFiles.length > 1" class="bg-gray-100 dark:bg-slate-800/50 rounded-lg p-4 mb-6">
        <div class="flex items-center justify-between mb-2">
          <span class="text-gray-900 dark:text-slate-100 font-medium">{{ selectedFiles.length }} files ({{ formattedTotalSize }})</span>
          <button @click="clearFiles" class="text-red-500 hover:text-red-600 dark:text-red-400 dark:hover:text-red-300 text-sm">&times; Clear all</button>
        </div>
        <ul class="text-sm text-gray-600 dark:text-slate-400 space-y-1 max-h-32 overflow-y-auto">
          <li v-for="f in selectedFiles" :key="f.name + f.size">{{ f.name }} ({{ formatFileSize(f.size) }})</li>
        </ul>
      </div>
    </div>

    <!-- Text Tab -->
    <div v-if="activeTab === 'text'">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-6">Share Text Secret</h2>
      <div class="mb-6">
        <textarea v-model="textInput" placeholder="Enter your secret text (max 1000 characters)..." rows="8" maxlength="1000"
          class="w-full px-4 py-3 bg-white dark:bg-slate-800 border border-gray-300 dark:border-slate-700 rounded-lg text-gray-900 dark:text-slate-200 resize-vertical focus:outline-none focus:border-blue-500 mb-2"></textarea>
        <p class="text-gray-600 dark:text-slate-400 text-sm text-right"><span class="text-gray-700 dark:text-slate-300">{{ textInput.length }}</span> / 1000 characters</p>
      </div>
    </div>

    <!-- PIN Input -->
    <div class="mb-6">
      <label class="block text-gray-700 dark:text-slate-300 font-medium mb-2">Set your PIN (4 characters):</label>
      <div class="flex items-center gap-3">
        <input
          v-model="pin"
          type="text"
          maxlength="4"
          autocomplete="off"
          autocorrect="off"
          autocapitalize="off"
          spellcheck="false"
          class="w-32 px-4 py-3 bg-white dark:bg-slate-800 border border-gray-300 dark:border-slate-700 rounded-lg text-gray-900 dark:text-slate-200 text-center text-lg font-mono tracking-widest focus:outline-none focus:border-blue-500"
        />
        <span class="text-gray-400 dark:text-slate-500 text-sm">{{ pin.length }}/4</span>
      </div>
      <div class="text-sm mt-1 min-h-[1.25rem]" :class="pinValidationClass">{{ pinValidationMsg }}</div>
      <p class="text-gray-500 dark:text-slate-500 text-xs mt-1">Letters & numbers only (case-sensitive)</p>
    </div>

    <!-- Access Mode Toggle -->
    <div class="mb-6">
      <label class="flex items-center cursor-pointer gap-3">
        <button type="button" role="switch" :aria-checked="oneTime"
          :class="['relative inline-flex h-6 w-11 shrink-0 rounded-full border-2 border-transparent transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2', oneTime ? 'bg-blue-600' : 'bg-gray-300 dark:bg-slate-600']"
          @click="oneTime = !oneTime">
          <span :class="['pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow transform transition-transform', oneTime ? 'translate-x-5' : 'translate-x-0']"></span>
        </button>
        <span class="text-gray-700 dark:text-slate-300 font-medium">Delete after first download</span>
      </label>
    </div>

    <!-- TTL -->
    <div class="mb-6">
      <label class="block text-gray-700 dark:text-slate-300 font-medium mb-3">Delete after:</label>
      <div class="flex flex-wrap gap-4">
        <label v-for="opt in [{ v: '1h', l: '1 hour' }, { v: '12h', l: '12 hours' }, { v: '24h', l: '24 hours' }]" :key="opt.v" class="flex items-center cursor-pointer">
          <input type="radio" :value="opt.v" v-model="ttl" class="mr-2 text-blue-600 focus:ring-blue-500">
          <span class="text-gray-700 dark:text-slate-300">{{ opt.l }}</span>
        </label>
      </div>
    </div>

    <!-- Upload Button -->
    <button class="btn-primary" :disabled="!canUpload || isUploading" @click="handleUpload">
      {{ activeTab === 'text' ? 'Encrypt & Share' : 'Encrypt & Upload' }}
    </button>

    <ProgressBar :visible="isUploading" :percent="progress" :text="progressText" />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import DropZone from '../../components/DropZone.vue';
import ProgressBar from '../../components/ProgressBar.vue';
import { useRecaptcha } from '../../composables/useRecaptcha.js';
import { API_BASE, RECAPTCHA_SITE_KEY } from '../../lib/api.js';
import { formatFileSize } from '../../lib/utils.js';
import { createBundle } from '../../lib/zip-bundle.js';
import * as CryptoModule from '../../lib/crypto.js';

const PIN_REGEX = /^[a-zA-Z0-9]{4}$/;
const TTL_LABELS = { '1h': '1 hour', '12h': '12 hours', '24h': '24 hours' };

const emit = defineEmits(['result']);
const { getToken } = useRecaptcha();

const activeTab = ref('file');
const selectedFiles = ref([]);
const textInput = ref('');
const pin = ref('');
const oneTime = ref(true);
const ttl = ref('24h');
const isUploading = ref(false);
const progress = ref(0);
const progressText = ref('');

const formattedTotalSize = computed(() => {
    const total = selectedFiles.value.reduce((sum, f) => sum + f.size, 0);
    return formatFileSize(total);
});

const pinValid = computed(() => PIN_REGEX.test(pin.value));

const pinValidationMsg = computed(() => {
    if (pin.value.length === 0) return '';
    if (pin.value.length < 4) return `${4 - pin.value.length} more character${pin.value.length === 3 ? '' : 's'} needed`;
    if (!PIN_REGEX.test(pin.value)) return 'Only letters and numbers allowed';
    return 'Valid PIN';
});

const pinValidationClass = computed(() => {
    if (pin.value.length === 0) return '';
    if (pin.value.length < 4) return 'text-gray-500 dark:text-slate-400';
    if (!PIN_REGEX.test(pin.value)) return 'text-red-500 dark:text-red-400';
    return 'text-green-600 dark:text-green-400';
});

const canUpload = computed(() => {
    if (!pinValid.value) return false;
    if (activeTab.value === 'file') return selectedFiles.value.length > 0;
    return textInput.value.trim().length > 0 && textInput.value.trim().length <= 1000;
});

function setFiles(files) {
    const arr = Array.from(files);
    if (arr.length > 10) { alert(`Maximum 10 files allowed. You selected ${arr.length}.`); return; }
    const totalSize = arr.reduce((sum, f) => sum + f.size, 0);
    if (totalSize > 500 * 1024 * 1024) { alert(`Total size exceeds 500 MB limit.`); return; }
    if (arr.some(f => f.size === 0)) { alert('Cannot upload empty files'); return; }
    selectedFiles.value = arr;
}

function clearFiles() { selectedFiles.value = []; }

function showProgress(pct, txt) { progress.value = pct; progressText.value = txt; }

function hexToUint8Array(hex) {
    const bytes = new Uint8Array(hex.length / 2);
    for (let i = 0; i < hex.length; i += 2) bytes[i / 2] = parseInt(hex.substr(i, 2), 16);
    return bytes;
}

async function uploadToS3(presignedUrl, data) {
    return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        xhr.upload.onprogress = (event) => {
            if (event.lengthComputable) {
                const pct = (event.loaded / event.total) * 100;
                showProgress(65 + pct * 0.30, `Uploading... ${pct.toFixed(1)}%`);
            }
        };
        xhr.onload = () => xhr.status >= 200 && xhr.status < 300 ? resolve() : reject(new Error(`S3 upload failed: ${xhr.status}`));
        xhr.onerror = () => reject(new Error('Network error during upload'));
        xhr.ontimeout = () => reject(new Error('Upload timeout'));
        xhr.open('PUT', presignedUrl);
        xhr.setRequestHeader('Content-Type', 'application/octet-stream');
        xhr.send(data);
    });
}

async function callPinUploadApi(body) {
    const response = await fetch(`${API_BASE}/pin/upload`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
    });
    if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || `Upload initialization failed: ${response.status}`);
    }
    return response.json();
}

async function handleUpload() {
    if (!canUpload.value || isUploading.value) return;

    try {
        isUploading.value = true;
        const accessMode = oneTime.value ? 'one_time' : 'multi';

        if (activeTab.value === 'text') {
            await handleTextUpload(accessMode);
        } else {
            await handleFileUpload(accessMode);
        }
    } catch (error) {
        console.error('PIN upload error:', error);
        alert(error.message || 'Upload failed. Please try again.');
        isUploading.value = false;
    }
}

async function handleFileUpload(accessMode) {
    showProgress(0, 'Preparing...');

    let fileToUpload, fileName;
    if (selectedFiles.value.length > 1) {
        showProgress(2, 'Creating ZIP bundle...');
        const bundle = await createBundle(selectedFiles.value, (pct) => showProgress(2 + pct * 0.08, `Bundling files... ${Math.round(pct)}%`));
        fileToUpload = new File([bundle.blob], bundle.filename, { type: 'application/zip' });
        fileName = bundle.filename;
    } else {
        fileToUpload = selectedFiles.value[0];
        fileName = selectedFiles.value[0].name;
    }

    showProgress(10, 'Verifying...');
    const recaptchaToken = await getToken('pin_upload');

    showProgress(15, 'Initializing upload...');
    const initResponse = await callPinUploadApi({
        content_type: 'file', file_size: fileToUpload.size, file_name: fileName,
        pin: pin.value, ttl: ttl.value, access_mode: accessMode, recaptcha_token: recaptchaToken,
    });

    const { file_id, upload_url, salt, expires_at } = initResponse;

    showProgress(20, 'Deriving encryption key...');
    const saltBytes = hexToUint8Array(salt);
    const encryptionKey = await CryptoModule.deriveKeyFromPassword(pin.value, saltBytes, true);

    showProgress(25, 'Encrypting... 0%');
    const encryptedData = await CryptoModule.encryptFile(fileToUpload, encryptionKey, (p) => {
        showProgress(25 + p * 0.35, `Encrypting... ${Math.round(p)}%`);
    });

    showProgress(65, 'Uploading encrypted file...');
    await uploadToS3(upload_url, encryptedData);

    showProgress(100, 'Upload complete!');
    emit('result', { fileId: file_id, pin: pin.value, ttl: ttl.value, expiresAt: expires_at });
}

async function handleTextUpload(accessMode) {
    showProgress(0, 'Preparing...');
    const text = textInput.value.trim();

    showProgress(10, 'Verifying...');
    const recaptchaToken = await getToken('pin_upload');

    showProgress(20, 'Initializing...');
    const textBlob = new Blob([new TextEncoder().encode(text)], { type: 'text/plain' });

    const initResponse = await callPinUploadApi({
        content_type: 'file', file_size: textBlob.size, file_name: 'secret.txt',
        pin: pin.value, ttl: ttl.value, access_mode: accessMode, recaptcha_token: recaptchaToken,
    });

    const { file_id, upload_url, salt, expires_at } = initResponse;

    showProgress(40, 'Deriving encryption key...');
    const saltBytes = hexToUint8Array(salt);
    const encryptionKey = await CryptoModule.deriveKeyFromPassword(pin.value, saltBytes, true);

    showProgress(50, 'Encrypting text...');
    const textFile = new File([textBlob], 'secret.txt', { type: 'text/plain' });
    const encryptedData = await CryptoModule.encryptFile(textFile, encryptionKey, (p) => showProgress(50 + p * 0.2, `Encrypting... ${Math.round(p)}%`));

    showProgress(75, 'Uploading...');
    await uploadToS3(upload_url, encryptedData);

    showProgress(100, 'Upload complete!');
    emit('result', { fileId: file_id, pin: pin.value, ttl: ttl.value, expiresAt: expires_at });
}
</script>
