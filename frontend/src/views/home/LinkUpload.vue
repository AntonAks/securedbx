<template>
  <div>
    <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-6">{{ $t('upload.link.shareFiles') }}</h2>

    <DropZone v-if="selectedFiles.length === 0" @files="addFiles" />
    <FileList :files="selectedFiles" @remove="removeFile" @clear="clearFiles" />

    <TtlSelector ref="ttlRef" />

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
import { validateFiles, createBundle } from '../../lib/zip-bundle.js';
import * as CryptoModule from '../../lib/crypto.js';

const { t } = useI18n();
const router = useRouter();
const upload = useUpload();
const ttlRef = ref(null);
const selectedFiles = ref([]);

function addFiles(files) {
    const combined = [...selectedFiles.value, ...Array.from(files)];
    const validation = validateFiles(combined);
    if (!validation.valid) {
        alert(validation.error);
        return;
    }
    selectedFiles.value = combined;
}

function removeFile(index) {
    selectedFiles.value.splice(index, 1);
}

function clearFiles() {
    selectedFiles.value = [];
}

async function handleUpload() {
    if (selectedFiles.value.length === 0) return;

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

        upload.updateProgress(15, t('upload.progress.generatingKey'));
        const key = await CryptoModule.generateKey();

        upload.updateProgress(20, t('upload.progress.encrypting', { percent: 0 }));
        const encryptedData = await CryptoModule.encryptFile(fileToUpload, key, (progress) => {
            upload.updateProgress(20 + (progress * 0.3), t('upload.progress.encrypting', { percent: Math.round(progress) }));
        });

        upload.updateProgress(50, t('upload.progress.initializing'));
        const ttl = ttlRef.value.getTtl();
        const fileSize = fileToUpload.size || fileToUpload.byteLength;
        const data = await upload.initializeUpload({ file_size: fileSize, ttl });

        upload.updateProgress(60, t('upload.progress.uploadingEncrypted'));
        await upload.uploadToS3(data.upload_url, encryptedData);

        upload.updateProgress(100, t('upload.progress.completeRedirect'));
        const keyBase64 = await CryptoModule.keyToBase64(key);
        const encodedFileName = encodeURIComponent(uploadFileName);

        setTimeout(() => {
            router.push({ path: '/share', query: { id: data.file_id, key: keyBase64, name: encodedFileName } });
        }, 500);
    } catch (error) {
        console.error('Upload error:', error);
        alert(error.message || t('upload.progress.uploadFailed'));
        upload.uploading.value = false;
    }
}
</script>
