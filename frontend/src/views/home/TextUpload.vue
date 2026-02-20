<template>
  <div>
    <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-6">{{ $t('upload.text.shareText') }}</h2>

    <div class="mb-6">
      <textarea
        v-model="text"
        :placeholder="$t('upload.text.placeholder')"
        rows="8"
        maxlength="1000"
        class="w-full px-4 py-3 bg-white dark:bg-slate-800 border border-gray-300 dark:border-slate-700 rounded-lg text-gray-900 dark:text-slate-200 resize-vertical focus:outline-none focus:border-blue-500 mb-2"
      ></textarea>
      <p class="text-gray-600 dark:text-slate-400 text-sm text-right">
        {{ $t('upload.text.charCount', { count: text.length }) }}
      </p>
    </div>

    <TtlSelector ref="ttlRef" />

    <button class="btn-primary" :disabled="text.trim().length === 0 || upload.uploading.value" @click="handleUpload">
      {{ $t('upload.text.encryptShare') }}
    </button>

    <ProgressBar :visible="upload.uploading.value" :percent="upload.progress.value" :text="upload.progressText.value" />
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import TtlSelector from '../../components/TtlSelector.vue';
import ProgressBar from '../../components/ProgressBar.vue';
import { useUpload } from '../../composables/useUpload.js';
import { useRecaptcha } from '../../composables/useRecaptcha.js';
import { API_BASE } from '../../lib/api.js';
import * as CryptoModule from '../../lib/crypto.js';

const { t } = useI18n();
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

        upload.updateProgress(10, t('upload.progress.generatingKey'));
        const key = await CryptoModule.generateKey();

        upload.updateProgress(30, t('upload.progress.encryptingText'));
        const encoder = new TextEncoder();
        const textBuffer = encoder.encode(trimmed);
        const encryptedData = await CryptoModule.encrypt(textBuffer, key);

        upload.updateProgress(50, t('upload.progress.preparingUpload'));
        const base64Text = btoa(String.fromCharCode(...encryptedData));
        const keyBase64 = await CryptoModule.keyToBase64(key);

        upload.updateProgress(60, t('upload.progress.verifying'));
        const recaptchaToken = await getToken('upload');

        upload.updateProgress(70, t('upload.progress.initializing'));
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
        upload.updateProgress(100, t('upload.progress.completeRedirect'));

        setTimeout(() => {
            router.push({ path: '/share', query: { id: data.file_id, key: keyBase64 } });
        }, 500);
    } catch (error) {
        console.error('Upload error:', error);
        alert(t('upload.progress.uploadFailedMsg', { message: error.message }));
        upload.uploading.value = false;
    }
}
</script>
