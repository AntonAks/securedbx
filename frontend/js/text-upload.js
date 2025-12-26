/**
 * Text Secret Upload Module
 */

(function() {
    'use strict';

    // Configuration
    const CONFIG = {
        API_BASE_URL: '/prod',
        RECAPTCHA_SITE_KEY: '6LdulTIsAAAAAJdhvyMU6B1og7GE7d5DySrQUQiv',
    };

    // DOM Elements
    const textInput = document.getElementById('text-input');
    const charCount = document.getElementById('char-count');
    const textUploadBtn = document.getElementById('text-upload-btn');
    const textProgressSection = document.getElementById('text-progress-section');
    const textProgressFill = document.getElementById('text-progress-fill');
    const textProgressText = document.getElementById('text-progress-text');
    const textResultSection = document.getElementById('text-result-section');
    const textShareUrl = document.getElementById('text-share-url');
    const textCopyBtn = document.getElementById('text-copy-btn');
    const newTextBtn = document.getElementById('new-text-btn');

    // Tab switching
    const tabButtons = document.querySelectorAll('.tab-btn');
    const tabPanels = document.querySelectorAll('.tab-panel');

    /**
     * Initialize text upload module
     * Sets up event listeners for tabs, text input, and buttons
     */
    function init() {
        // Tab switching
        tabButtons.forEach(btn => {
            btn.addEventListener('click', handleTabSwitch);
        });

        // Text input
        textInput.addEventListener('input', handleTextInput);

        // Upload button
        textUploadBtn.addEventListener('click', handleTextUpload);

        // Copy button
        textCopyBtn.addEventListener('click', handleCopyUrl);

        // New text button
        newTextBtn.addEventListener('click', resetTextUpload);
    }

    /**
     * Handle tab switching between file and text tabs
     * @param {Event} e - Click event
     */
    function handleTabSwitch(e) {
        const targetTab = e.target.dataset.tab;

        // Update tab buttons
        tabButtons.forEach(btn => btn.classList.remove('active'));
        e.target.classList.add('active');

        // Update tab panels
        tabPanels.forEach(panel => {
            if (panel.id === targetTab) {
                panel.classList.add('active');
            } else {
                panel.classList.remove('active');
            }
        });
    }

    /**
     * Handle text input and update character counter
     * Enables/disables upload button based on text length
     */
    function handleTextInput() {
        const length = textInput.value.length;
        charCount.textContent = length;

        // Enable/disable upload button
        textUploadBtn.disabled = length === 0;
    }

    /**
     * Get selected TTL value for text secret
     * @returns {string} TTL value: "1h", "12h", or "24h"
     */
    function getTextTTL() {
        const selected = document.querySelector('input[name="text-ttl"]:checked');
        return selected ? selected.value : '1h';
    }

    /**
     * Handle text upload: encrypt and share text secret
     * @returns {Promise<void>}
     * @throws {Error} If upload fails
     */
    async function handleTextUpload() {
        const text = textInput.value.trim();

        if (!text) {
            alert('Please enter some text to share');
            return;
        }

        try {
            // Disable button and show progress
            textUploadBtn.disabled = true;
            textProgressSection.style.display = 'block';
            updateProgress(10, 'Generating encryption key...');

            // Generate encryption key
            const key = await CryptoModule.generateKey();
            updateProgress(30, 'Encrypting text...');

            // Encrypt text
            const encoder = new TextEncoder();
            const textBuffer = encoder.encode(text);
            const encryptedData = await CryptoModule.encrypt(textBuffer, key);
            updateProgress(50, 'Preparing upload...');

            // Convert encrypted data to base64
            const base64Text = btoa(String.fromCharCode(...encryptedData));

            // Export key for URL (use base64 like file uploads)
            const keyBase64 = await CryptoModule.keyToBase64(key);

            // Get reCAPTCHA token
            updateProgress(60, 'Verifying...');
            const recaptchaToken = await grecaptcha.execute(CONFIG.RECAPTCHA_SITE_KEY, { action: 'upload' });

            // Initialize upload
            updateProgress(70, 'Initializing...');
            const ttl = getTextTTL();
            const response = await fetch(`${CONFIG.API_BASE_URL}/upload/init`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    content_type: 'text',
                    encrypted_text: base64Text,
                    ttl: ttl,
                    recaptcha_token: recaptchaToken,
                }),
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Upload failed');
            }

            const data = await response.json();
            updateProgress(100, 'Complete!');

            // Generate share URL with key (same format as files)
            const shareUrl = `${window.location.origin}/download.html#${data.file_id}#${keyBase64}`;

            // Show result
            setTimeout(() => {
                textProgressSection.style.display = 'none';
                textShareUrl.value = shareUrl;
                textResultSection.style.display = 'block';
            }, 500);

        } catch (error) {
            console.error('Upload error:', error);
            alert(`Upload failed: ${error.message}`);
            resetTextUpload();
        }
    }

    /**
     * Update progress indicator during upload
     * @param {number} percent - Progress percentage (0-100)
     * @param {string} message - Status message to display
     */
    function updateProgress(percent, message) {
        textProgressFill.style.width = `${percent}%`;
        textProgressText.textContent = message;
    }

    /**
     * Handle copy URL to clipboard
     * @returns {Promise<void>}
     */
    async function handleCopyUrl() {
        try {
            await navigator.clipboard.writeText(textShareUrl.value);

            const originalText = textCopyBtn.textContent;
            textCopyBtn.textContent = '✓ Copied!';
            textCopyBtn.classList.add('btn-success');

            setTimeout(() => {
                textCopyBtn.textContent = originalText;
                textCopyBtn.classList.remove('btn-success');
            }, 2000);
        } catch (error) {
            // Fallback: select text
            textShareUrl.select();
            alert('Press Ctrl+C (or Cmd+C on Mac) to copy');
        }
    }

    /**
     * Reset text upload form to initial state
     * Clears input, hides progress and result sections
     */
    function resetTextUpload() {
        textInput.value = '';
        charCount.textContent = '0';
        textUploadBtn.disabled = true;
        textProgressSection.style.display = 'none';
        textResultSection.style.display = 'none';
        textProgressFill.style.width = '0%';
        textProgressText.textContent = '';
    }

    // Initialize when DOM is ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }

})();
