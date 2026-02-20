<template>
  <div>
    <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-2">{{ $t('upload.vault.title') }}</h2>
    <p class="text-gray-600 dark:text-slate-400 mb-6">{{ $t('upload.vault.subtitle') }}</p>

    <!-- Content Type Toggle -->
    <div class="flex justify-center mb-6">
      <div class="inline-flex rounded-lg bg-gray-200 dark:bg-slate-800 p-1">
        <label class="cursor-pointer">
          <input type="radio" value="file" v-model="contentType" class="sr-only peer">
          <span class="px-4 py-2 rounded-lg text-sm font-medium peer-checked:bg-white dark:peer-checked:bg-slate-700 peer-checked:shadow text-gray-600 dark:text-slate-400 peer-checked:text-gray-900 dark:peer-checked:text-slate-100 transition-all">{{ $t('upload.vault.file') }}</span>
        </label>
        <label class="cursor-pointer">
          <input type="radio" value="text" v-model="contentType" class="sr-only peer">
          <span class="px-4 py-2 rounded-lg text-sm font-medium peer-checked:bg-white dark:peer-checked:bg-slate-700 peer-checked:shadow text-gray-600 dark:text-slate-400 peer-checked:text-gray-900 dark:peer-checked:text-slate-100 transition-all">{{ $t('upload.vault.text') }}</span>
        </label>
      </div>
    </div>

    <!-- File Section -->
    <div v-if="contentType === 'file'">
      <DropZone v-if="selectedFiles.length === 0" @files="addFiles" />
      <FileList :files="selectedFiles" @remove="removeFile" @clear="clearFiles" />
    </div>

    <!-- Text Section -->
    <div v-else class="mb-6">
      <textarea v-model="textInput" :placeholder="$t('upload.vault.textPlaceholder')" rows="6" maxlength="10000"
        class="w-full px-4 py-3 bg-white dark:bg-slate-800 border border-gray-300 dark:border-slate-700 rounded-lg text-gray-900 dark:text-slate-200 resize-vertical focus:outline-none focus:border-blue-500 mb-2"></textarea>
      <p class="text-gray-600 dark:text-slate-400 text-sm text-right">{{ $t('upload.vault.textCharCount', { count: textInput.length.toLocaleString() }) }}</p>
    </div>

    <!-- Password -->
    <div class="mb-6">
      <label class="block text-gray-700 dark:text-slate-300 font-medium mb-2">
        {{ $t('upload.vault.password') }} <span class="text-red-500">*</span>
      </label>
      <div class="flex items-center gap-2">
        <input :type="showPassword ? 'text' : 'password'" v-model="password" :placeholder="$t('upload.vault.passwordPlaceholder')"
          class="flex-1 px-3 py-2 bg-white dark:bg-slate-900 border border-gray-300 dark:border-slate-600 rounded text-gray-900 dark:text-slate-200 focus:outline-none focus:border-blue-500">
        <button type="button" @click="showPassword = !showPassword" class="px-3 py-2 text-gray-500 hover:text-gray-700 dark:text-slate-400 dark:hover:text-slate-200 text-sm border border-gray-300 dark:border-slate-600 rounded">
          {{ showPassword ? $t('upload.vault.hide') : $t('upload.vault.show') }}
        </button>
      </div>
      <p class="text-gray-500 dark:text-slate-500 text-sm mt-1">{{ $t('upload.vault.passwordHint') }}</p>
    </div>

    <TtlSelector ref="ttlRef" />

    <button class="btn-primary" :disabled="!canUpload || upload.uploading.value" @click="handleUpload">
      <svg class="inline w-4 h-4 mr-2 -mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
        <path d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
      </svg>
      {{ $t('upload.vault.encryptStore') }}
    </button>

    <ProgressBar :visible="upload.uploading.value" :percent="upload.progress.value" :text="upload.progressText.value" />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import DropZone from '../../components/DropZone.vue';
import FileList from '../../components/FileList.vue';
import TtlSelector from '../../components/TtlSelector.vue';
import ProgressBar from '../../components/ProgressBar.vue';
import { useUpload } from '../../composables/useUpload.js';
import { useRecaptcha } from '../../composables/useRecaptcha.js';
import { validateFiles, createBundle } from '../../lib/zip-bundle.js';
import { API_BASE } from '../../lib/api.js';
import * as CryptoModule from '../../lib/crypto.js';

const { t } = useI18n();
const router = useRouter();
const upload = useUpload();
const { getToken } = useRecaptcha();
const ttlRef = ref(null);

const contentType = ref('file');
const selectedFiles = ref([]);
const textInput = ref('');
const password = ref('');
const showPassword = ref(false);

const canUpload = computed(() => {
    const hasPassword = password.value.length >= 4;
    if (contentType.value === 'file') return selectedFiles.value.length > 0 && hasPassword;
    return textInput.value.trim().length > 0 && hasPassword;
});

function addFiles(files) {
    const combined = [...selectedFiles.value, ...Array.from(files)];
    const validation = validateFiles(combined);
    if (!validation.valid) { alert(validation.error); return; }
    selectedFiles.value = combined;
}

function removeFile(index) { selectedFiles.value.splice(index, 1); }
function clearFiles() { selectedFiles.value = []; }

async function handleUpload() {
    if (!canUpload.value) return;

    try {
        upload.uploading.value = true;

        let fileToUpload, uploadFileName;
        if (contentType.value === 'file') {
            if (selectedFiles.value.length === 1) {
                fileToUpload = selectedFiles.value[0];
                uploadFileName = selectedFiles.value[0].name;
                upload.updateProgress(0, t('upload.progress.preparingFile'));
            } else {
                upload.updateProgress(0, t('upload.progress.creatingZip'));
                const bundle = await createBundle(selectedFiles.value, (p, msg) => upload.updateProgress(p * 0.05, msg));
                fileToUpload = bundle.blob;
                uploadFileName = bundle.filename;
            }
        } else {
            fileToUpload = new TextEncoder().encode(textInput.value.trim());
            uploadFileName = null;
        }

        upload.updateProgress(5, t('upload.progress.generatingKey'));
        const dataKey = await CryptoModule.generateKey();

        upload.updateProgress(10, t('upload.progress.derivingPasswordKey'));
        const salt = CryptoModule.generateSalt();
        const passwordKey = await CryptoModule.deriveKeyFromPassword(password.value, salt);

        upload.updateProgress(15, t('upload.progress.securingKey'));
        const encryptedKeyData = await CryptoModule.encryptKey(dataKey, passwordKey);
        const encryptedKeyBase64 = CryptoModule.arrayToBase64(encryptedKeyData);
        const saltBase64 = CryptoModule.arrayToBase64(salt);

        upload.updateProgress(20, t('upload.progress.encryptingContent'));
        let encryptedData;
        if (contentType.value === 'file') {
            encryptedData = await CryptoModule.encryptFile(fileToUpload, dataKey, (p) => {
                upload.updateProgress(20 + p * 0.3, t('upload.progress.encrypting', { percent: Math.round(p) }));
            });
        } else {
            encryptedData = await CryptoModule.encrypt(fileToUpload, dataKey);
        }

        upload.updateProgress(50, t('upload.progress.initializingVault'));
        const ttl = ttlRef.value.getTtl();
        const recaptchaToken = await getToken('vault');

        const requestBody = {
            content_type: contentType.value,
            ttl,
            recaptcha_token: recaptchaToken,
            access_mode: 'multi',
            salt: saltBase64,
            encrypted_key: encryptedKeyBase64,
        };

        if (contentType.value === 'file') {
            requestBody.file_size = encryptedData.byteLength;
        } else {
            requestBody.encrypted_text = CryptoModule.arrayToBase64(encryptedData);
        }

        const response = await fetch(`${API_BASE}/upload/init`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(requestBody),
        });

        if (!response.ok) {
            const error = await response.json().catch(() => ({}));
            throw new Error(error.error || `Upload failed: ${response.status}`);
        }

        const data = await response.json();

        if (contentType.value === 'file') {
            upload.updateProgress(60, t('upload.progress.uploadingEncrypted'));
            await upload.uploadToS3(data.upload_url, encryptedData);
        }

        upload.updateProgress(100, t('upload.progress.completeRedirect'));

        const encodedFileName = uploadFileName ? encodeURIComponent(uploadFileName) : '';
        setTimeout(() => {
            router.push({ path: '/share', query: { id: data.file_id, salt: saltBase64, name: encodedFileName, vault: '1' } });
        }, 500);
    } catch (error) {
        console.error('Vault upload error:', error);
        alert(error.message || t('upload.progress.uploadFailed'));
        upload.uploading.value = false;
    }
}
</script>
