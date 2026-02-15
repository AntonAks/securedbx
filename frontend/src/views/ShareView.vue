<template>
  <main class="flex-1 max-w-3xl mx-auto px-6 py-12 w-full">
    <!-- Error -->
    <section v-if="hasError" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <div class="bg-red-500/10 border border-red-500 rounded-lg p-4 mb-6">
        <h2 class="text-red-500 font-semibold text-xl mb-2">Invalid Share Link</h2>
        <p class="text-gray-700 dark:text-slate-300">The share link is missing required information.</p>
      </div>
      <router-link to="/" class="btn-primary block text-center">Go to Upload Page</router-link>
    </section>

    <!-- Share Section -->
    <section v-else class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <!-- Success Header -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-green-500/10 mb-4">
          <svg class="w-8 h-8 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100">Ready to Share!</h2>
        <p class="text-gray-600 dark:text-slate-400 mt-2" v-if="isVault">
          <span class="inline-flex items-center gap-2">
            <svg class="w-5 h-5 text-blue-600 dark:text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
            </svg>
            {{ fileName ? 'Your file has been stored in the vault.' : 'Your text has been stored in the vault.' }}
          </span>
        </p>
        <p v-else class="text-gray-600 dark:text-slate-400 mt-2">
          {{ isTextSecret ? 'Your text secret has been encrypted and stored.' : 'Your file has been encrypted and uploaded.' }}
        </p>
      </div>

      <!-- QR Code -->
      <div class="flex flex-col items-center mb-8">
        <QrCode ref="qrRef" :value="downloadUrl" />
        <p class="text-gray-500 dark:text-slate-400 text-sm mt-3">Scan to download</p>
      </div>

      <!-- QR Action Buttons -->
      <div class="flex justify-center gap-3 mb-8">
        <button @click="handleCopyQr" class="btn-secondary btn-auto flex items-center gap-2">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
          {{ qrCopied ? 'Copied!' : 'Copy QR' }}
        </button>
        <button @click="handleDownloadQr" class="btn-secondary btn-auto flex items-center gap-2">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
          Download QR
        </button>
      </div>

      <!-- Copy Link Section -->
      <div class="mb-8">
        <p class="text-gray-700 dark:text-slate-300 text-sm mb-2">Or copy link:</p>
        <div class="flex gap-2">
          <input type="text" :value="downloadUrl" readonly
            class="flex-1 px-4 py-3 bg-white dark:bg-slate-800 border border-gray-300 dark:border-slate-700 rounded-lg text-gray-900 dark:text-slate-200 font-mono text-sm focus:outline-none focus:border-blue-500">
          <button @click="handleCopyLink" class="btn-secondary btn-auto">
            {{ linkCopied ? 'Copied!' : 'Copy' }}
          </button>
        </div>
      </div>

      <!-- Upload Another Button -->
      <div class="text-center mb-8">
        <router-link to="/" class="btn-outline inline-flex items-center gap-2">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
          </svg>
          Upload Another
        </router-link>
      </div>

      <!-- Warning Section -->
      <div class="bg-yellow-500/10 border border-yellow-500 rounded-lg p-4 mb-6">
        <p class="text-yellow-500 font-semibold mb-2">Warning:</p>
        <ul class="text-gray-700 dark:text-slate-300 text-sm space-y-1 ml-5 list-disc">
          <li v-if="isVault">
            <span class="text-blue-600 dark:text-blue-400 font-medium">Password-protected</span>
            &bull; Can be accessed multiple times until expiry
          </li>
          <template v-else>
            <li>This link works only <strong class="text-gray-900 dark:text-white">ONCE</strong></li>
            <li>{{ isTextSecret ? 'Text deleted after viewing or expiry' : 'File deleted after download or expiry' }}</li>
          </template>
        </ul>
      </div>

      <!-- Expiration Countdown -->
      <div v-if="countdownText" class="text-center">
        <p class="text-gray-500 dark:text-slate-400 text-sm">
          <svg class="w-4 h-4 inline-block mr-1 -mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Expires in: <span class="font-medium text-gray-700 dark:text-slate-300" :class="{ 'text-red-500': countdownText === 'Expired' }">{{ countdownText }}</span>
        </p>
      </div>
    </section>
  </main>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';
import QrCode from '../components/QrCode.vue';
import { useClipboard } from '../composables/useClipboard.js';
import { API_BASE } from '../lib/api.js';
import { formatTimeRemaining } from '../lib/utils.js';

const route = useRoute();
const { copied: linkCopied, copy: copyLink } = useClipboard();

const qrRef = ref(null);
const qrCopied = ref(false);
const countdownText = ref('');
let countdownInterval = null;

const fileId = route.query.id || '';
const key = route.query.key || '';
const salt = route.query.salt || '';
const fileName = route.query.name ? decodeURIComponent(route.query.name) : '';
const isVault = route.query.vault === '1';

const hasError = computed(() => {
    if (!fileId) return true;
    if (isVault && !salt) return true;
    if (!isVault && !key) return true;
    return false;
});

const isTextSecret = computed(() => !isVault && !fileName);

const downloadUrl = computed(() => {
    const base = `${window.location.origin}/#/download`;
    const params = new URLSearchParams();
    params.set('id', fileId);

    if (isVault) {
        params.set('salt', salt);
        if (fileName) params.set('name', encodeURIComponent(fileName));
        params.set('vault', '1');
    } else {
        params.set('key', key);
        if (fileName) params.set('name', encodeURIComponent(fileName));
    }

    return `${base}?${params.toString()}`;
});

onMounted(() => {
    if (!hasError.value) fetchMetadataAndStartCountdown();
});

onUnmounted(() => {
    if (countdownInterval) clearInterval(countdownInterval);
});

async function fetchMetadataAndStartCountdown() {
    try {
        const response = await fetch(`${API_BASE}/files/${fileId}/metadata`);
        if (!response.ok) return;
        const data = await response.json();
        if (data.expires_at) startCountdown(data.expires_at);
    } catch (error) {
        console.warn('Failed to fetch metadata:', error);
    }
}

function startCountdown(expiresAt) {
    const update = () => {
        const remaining = expiresAt - Math.floor(Date.now() / 1000);
        if (remaining <= 0) {
            countdownText.value = 'Expired';
            if (countdownInterval) clearInterval(countdownInterval);
            return;
        }
        countdownText.value = formatTimeRemaining(remaining);
    };
    update();
    countdownInterval = setInterval(update, 60000);
}

function handleCopyLink() {
    copyLink(downloadUrl.value);
}

async function handleCopyQr() {
    const canvas = qrRef.value?.canvasRef;
    if (!canvas) return;

    try {
        const blob = await new Promise((resolve, reject) => {
            canvas.toBlob((b) => b ? resolve(b) : reject(new Error('Failed to create blob')), 'image/png');
        });
        await navigator.clipboard.write([new ClipboardItem({ 'image/png': blob })]);
        qrCopied.value = true;
        setTimeout(() => { qrCopied.value = false; }, 2000);
    } catch (error) {
        console.warn('Copy QR failed:', error);
        alert('Browser does not support copying images. Use "Download QR" instead.');
    }
}

function handleDownloadQr() {
    const canvas = qrRef.value?.canvasRef;
    if (!canvas) return;

    canvas.toBlob((blob) => {
        if (!blob) return;
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = fileName ? `sdbx-qr-${fileName.replace(/\.[^/.]+$/, '')}.png` : 'sdbx-qr-code.png';
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    }, 'image/png');
}
</script>
