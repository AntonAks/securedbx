/**
 * Cryptography module for sdbx
 * Handles AES-256-GCM encryption and decryption
 */

'use strict';

const CryptoModule = (function() {
    const ALGORITHM = 'AES-GCM';
    const KEY_LENGTH = 256;
    const IV_LENGTH = 12; // 96 bits

    /**
     * Generate a random encryption key
     * @returns {Promise<CryptoKey>} Generated key
     */
    async function generateKey() {
        return crypto.subtle.generateKey(
            {
                name: ALGORITHM,
                length: KEY_LENGTH,
            },
            true, // extractable
            ['encrypt', 'decrypt']
        );
    }

    /**
     * Generate random IV (Initialization Vector)
     * @returns {Uint8Array} Random IV
     */
    function generateIV() {
        return crypto.getRandomValues(new Uint8Array(IV_LENGTH));
    }

    /**
     * Encrypt data with AES-256-GCM
     * @param {ArrayBuffer} data - Data to encrypt
     * @param {CryptoKey} key - Encryption key
     * @returns {Promise<Uint8Array>} IV + ciphertext
     */
    async function encrypt(data, key) {
        const iv = generateIV();

        const ciphertext = await crypto.subtle.encrypt(
            {
                name: ALGORITHM,
                iv: iv,
            },
            key,
            data
        );

        // Prepend IV to ciphertext
        const result = new Uint8Array(iv.length + ciphertext.byteLength);
        result.set(iv, 0);
        result.set(new Uint8Array(ciphertext), iv.length);

        return result;
    }

    /**
     * Decrypt data with AES-256-GCM
     * @param {Uint8Array} data - IV + ciphertext
     * @param {CryptoKey} key - Decryption key
     * @returns {Promise<ArrayBuffer>} Decrypted data
     */
    async function decrypt(data, key) {
        // Extract IV and ciphertext
        const iv = data.slice(0, IV_LENGTH);
        const ciphertext = data.slice(IV_LENGTH);

        return crypto.subtle.decrypt(
            {
                name: ALGORITHM,
                iv: iv,
            },
            key,
            ciphertext
        );
    }

    /**
     * Export key to raw format
     * @param {CryptoKey} key - Key to export
     * @returns {Promise<ArrayBuffer>} Raw key bytes
     */
    async function exportKey(key) {
        return crypto.subtle.exportKey('raw', key);
    }

    /**
     * Import key from raw format
     * @param {ArrayBuffer} keyData - Raw key bytes
     * @returns {Promise<CryptoKey>} Imported key
     */
    async function importKey(keyData) {
        return crypto.subtle.importKey(
            'raw',
            keyData,
            {
                name: ALGORITHM,
            },
            true,
            ['encrypt', 'decrypt']
        );
    }

    /**
     * Convert key to base64 for URL
     * @param {CryptoKey} key - Key to convert
     * @returns {Promise<string>} Base64-encoded key
     */
    async function keyToBase64(key) {
        const rawKey = await exportKey(key);
        const keyArray = new Uint8Array(rawKey);
        const keyString = String.fromCharCode.apply(null, keyArray);
        return btoa(keyString);
    }

    /**
     * Convert base64 string to key
     * @param {string} base64 - Base64-encoded key
     * @returns {Promise<CryptoKey>} Imported key
     */
    async function base64ToKey(base64) {
        const keyString = atob(base64);
        const keyArray = new Uint8Array(keyString.length);
        for (let i = 0; i < keyString.length; i++) {
            keyArray[i] = keyString.charCodeAt(i);
        }
        return importKey(keyArray.buffer);
    }

    /**
     * Encrypt file using Web Worker (non-blocking)
     * @param {File} file - File to encrypt
     * @param {CryptoKey} key - Encryption key
     * @param {function(number): void} onProgress - Progress callback (0-100)
     * @returns {Promise<Uint8Array>} Encrypted data
     */
    async function encryptFile(file, key, onProgress) {
        return new Promise(async (resolve, reject) => {
            try {
                // Create worker
                const worker = new Worker('/js/crypto-worker.js');

                // Handle worker messages
                worker.onmessage = function(e) {
                    const { type, percent, result, error } = e.data;

                    if (type === 'progress') {
                        if (onProgress) {
                            // Map worker progress (0-100) to overall progress (10-100)
                            // Reserve 0-10% for file reading
                            const overallProgress = 10 + (percent * 0.9);
                            onProgress(overallProgress);
                        }
                    } else if (type === 'complete') {
                        worker.terminate();
                        resolve(new Uint8Array(result));
                    } else if (type === 'error') {
                        worker.terminate();
                        reject(new Error(error));
                    }
                };

                worker.onerror = function(error) {
                    worker.terminate();
                    reject(error);
                };

                // Read file
                if (onProgress) onProgress(0);
                const arrayBuffer = await file.arrayBuffer();
                if (onProgress) onProgress(10);

                // Export key to raw bytes
                const keyData = await exportKey(key);

                // Send to worker
                worker.postMessage({
                    action: 'encrypt',
                    data: arrayBuffer,
                    keyData: keyData
                }, [arrayBuffer]); // Transfer ownership for efficiency

            } catch (error) {
                reject(error);
            }
        });
    }

    /**
     * Decrypt file using Web Worker (non-blocking)
     * @param {Uint8Array} encryptedData - Encrypted data
     * @param {CryptoKey} key - Decryption key
     * @param {function(number): void} onProgress - Progress callback (0-100)
     * @returns {Promise<ArrayBuffer>} Decrypted data
     */
    async function decryptFile(encryptedData, key, onProgress) {
        return new Promise(async (resolve, reject) => {
            try {
                // Create worker
                const worker = new Worker('/js/crypto-worker.js');

                // Handle worker messages
                worker.onmessage = function(e) {
                    const { type, percent, result, error } = e.data;

                    if (type === 'progress') {
                        if (onProgress) {
                            onProgress(percent);
                        }
                    } else if (type === 'complete') {
                        worker.terminate();
                        resolve(result);
                    } else if (type === 'error') {
                        worker.terminate();
                        reject(new Error(error));
                    }
                };

                worker.onerror = function(error) {
                    worker.terminate();
                    reject(error);
                };

                // Export key to raw bytes
                if (onProgress) onProgress(0);
                const keyData = await exportKey(key);

                // Send to worker
                worker.postMessage({
                    action: 'decrypt',
                    data: encryptedData.buffer,
                    keyData: keyData
                }, [encryptedData.buffer]); // Transfer ownership for efficiency

            } catch (error) {
                reject(error);
            }
        });
    }

    return {
        generateKey,
        encrypt,
        decrypt,
        encryptFile,
        decryptFile,
        keyToBase64,
        base64ToKey,
    };
})();
