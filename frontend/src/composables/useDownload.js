import { ref } from 'vue';
import { API_BASE } from '../lib/api.js';
import { useRecaptcha } from './useRecaptcha.js';

/**
 * Composable for download operations
 */
export function useDownload() {
    const progress = ref(0);
    const progressText = ref('');
    const downloading = ref(false);

    const { getToken } = useRecaptcha();

    function updateProgress(percent, text) {
        progress.value = percent;
        progressText.value = text;
    }

    /**
     * Download file from API
     * @param {string} fileId - File identifier
     * @returns {Promise<{data: ArrayBuffer, headers: Headers}>}
     */
    async function downloadFile(fileId) {
        const recaptchaToken = await getToken('download');

        const response = await fetch(`${API_BASE}/download/${fileId}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ recaptcha_token: recaptchaToken }),
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.error || `Download failed: ${response.status}`);
        }

        return response;
    }

    /**
     * Save decrypted data as a file download
     * @param {ArrayBuffer} data - Decrypted file data
     * @param {string} filename - Filename for the download
     */
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

    return {
        progress,
        progressText,
        downloading,
        updateProgress,
        downloadFile,
        saveFile,
    };
}
