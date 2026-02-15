<template>
  <div>
    <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-6">Share Text Secret</h2>

    <div class="mb-6">
      <textarea
        v-model="text"
        placeholder="Enter your secret text (max 1000 characters)..."
        rows="8"
        maxlength="1000"
        class="w-full px-4 py-3 bg-white dark:bg-slate-800 border border-gray-300 dark:border-slate-700 rounded-lg text-gray-900 dark:text-slate-200 resize-vertical focus:outline-none focus:border-blue-500 mb-2"
      ></textarea>
      <p class="text-gray-600 dark:text-slate-400 text-sm text-right">
        <span class="text-gray-700 dark:text-slate-300">{{ text.length }}</span> / 1000 characters
      </p>
    </div>

    <TtlSelector ref="ttlRef" />

    <button class="btn-primary" :disabled="text.trim().length === 0 || upload.uploading.value" @click="handleUpload">
      Encrypt & Share
    </button>

    <ProgressBar :visible="upload.uploading.value" :percent="upload.progress.value" :text="upload.progressText.value" />
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import TtlSelector from '../../components/TtlSelector.vue';
import ProgressBar from '../../components/ProgressBar.vue';
import { useUpload } from '../../composables/useUpload.js';
import { useRecaptcha } from '../../composables/useRecaptcha.js';
import { API_BASE } from '../../lib/api.js';
import * as CryptoModule from '../../lib/crypto.js';

const router = useRouter();
const upload = useUpload();
const { getToken } = useRecaptcha();
const ttlRef = ref(null);
const text = ref('');

async function handleUpload() {
    const trimmed = text.value.trim();
    if (!trimmed) return;

    try {
        upload.uploading.value = true;

        upload.updateProgress(10, 'Generating encryption key...');
        const key = await CryptoModule.generateKey();

        upload.updateProgress(30, 'Encrypting text...');
        const encoder = new TextEncoder();
        const textBuffer = encoder.encode(trimmed);
        const encryptedData = await CryptoModule.encrypt(textBuffer, key);

        upload.updateProgress(50, 'Preparing upload...');
        const base64Text = btoa(String.fromCharCode(...encryptedData));
        const keyBase64 = await CryptoModule.keyToBase64(key);

        upload.updateProgress(60, 'Verifying...');
        const recaptchaToken = await getToken('upload');

        upload.updateProgress(70, 'Initializing...');
        const ttl = ttlRef.value.getTtl();
        const response = await fetch(`${API_BASE}/upload/init`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                content_type: 'text',
                encrypted_text: base64Text,
                ttl,
                recaptcha_token: recaptchaToken,
            }),
        });

        if (!response.ok) {
            const error = await response.json().catch(() => ({}));
            throw new Error(error.error || 'Upload failed');
        }

        const data = await response.json();
        upload.updateProgress(100, 'Complete! Redirecting...');

        setTimeout(() => {
            router.push({ path: '/share', query: { id: data.file_id, key: keyBase64 } });
        }, 500);
    } catch (error) {
        console.error('Upload error:', error);
        alert(`Upload failed: ${error.message}`);
        upload.uploading.value = false;
    }
}
</script>
