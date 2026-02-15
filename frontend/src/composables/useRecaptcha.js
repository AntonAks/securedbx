import { RECAPTCHA_SITE_KEY } from '../lib/api.js';

/**
 * Composable for reCAPTCHA v3 token generation
 */
export function useRecaptcha() {
    async function getToken(action) {
        if (!window.grecaptcha) {
            throw new Error('reCAPTCHA not loaded');
        }

        return new Promise((resolve, reject) => {
            grecaptcha.ready(async () => {
                try {
                    const token = await grecaptcha.execute(RECAPTCHA_SITE_KEY, { action });
                    resolve(token);
                } catch (error) {
                    console.error('reCAPTCHA error:', error);
                    reject(new Error('Failed to verify reCAPTCHA. Please refresh and try again.'));
                }
            });
        });
    }

    return { getToken };
}
