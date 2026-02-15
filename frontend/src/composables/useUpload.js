import { ref } from 'vue';
import { API_BASE } from '../lib/api.js';
import { useRecaptcha } from './useRecaptcha.js';

/**
 * Composable for upload operations (init + S3 upload)
 */
export function useUpload() {
    const progress = ref(0);
    const progressText = ref('');
    const uploading = ref(false);

    const { getToken } = useRecaptcha();

    function updateProgress(percent, text) {
        progress.value = percent;
        progressText.value = text;
    }

    /**
     * Initialize upload via API
     * @param {Object} params - Upload parameters
     * @returns {Promise<Object>} API response with upload_url, file_id, etc.
     */
    async function initializeUpload(params) {
        const recaptchaToken = await getToken('upload');

        const response = await fetch(`${API_BASE}/upload/init`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                ...params,
                recaptcha_token: recaptchaToken,
            }),
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.error || `Upload initialization failed: ${response.status}`);
        }

        return response.json();
    }

    /**
     * Upload encrypted data to S3 via presigned URL
     * @param {string} presignedUrl - S3 presigned URL
     * @param {Uint8Array} data - Encrypted data
     * @param {number} progressStart - Start of progress range (0-100)
     * @param {number} progressEnd - End of progress range (0-100)
     */
    async function uploadToS3(presignedUrl, data, progressStart = 60, progressEnd = 90) {
        return new Promise((resolve, reject) => {
            const xhr = new XMLHttpRequest();

            xhr.upload.onprogress = (event) => {
                if (event.lengthComputable) {
                    const pct = (event.loaded / event.total) * 100;
                    const overall = progressStart + (pct / 100) * (progressEnd - progressStart);
                    updateProgress(overall, `Uploading... ${pct.toFixed(1)}%`);
                }
            };

            xhr.onload = () => {
                if (xhr.status >= 200 && xhr.status < 300) {
                    resolve();
                } else {
                    reject(new Error(`S3 upload failed: ${xhr.status}`));
                }
            };

            xhr.onerror = () => reject(new Error('Network error during upload'));
            xhr.ontimeout = () => reject(new Error('Upload timeout'));

            xhr.open('PUT', presignedUrl);
            xhr.setRequestHeader('Content-Type', 'application/octet-stream');
            xhr.send(data);
        });
    }

    return {
        progress,
        progressText,
        uploading,
        updateProgress,
        initializeUpload,
        uploadToS3,
    };
}
