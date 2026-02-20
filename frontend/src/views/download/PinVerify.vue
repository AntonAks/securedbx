<template>
  <div>
    <!-- Code Entry -->
    <section v-if="step === 'code'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <PinCodeEntry :error="codeError" :loading="codeLoading" @submit="handleCodeSubmit" ref="codeEntryRef" />
    </section>

    <!-- PIN Entry with Timer -->
    <section v-if="step === 'pin'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <div class="flex items-center justify-between mb-4">
        <span class="text-gray-600 dark:text-slate-400 text-sm">{{ $t('download.pin.fileLabel') }} <strong class="text-gray-900 dark:text-slate-100 font-mono">{{ fileId }}</strong></span>
      </div>

      <div class="text-center py-3 px-4 rounded-lg mb-4 text-lg font-semibold" :class="timerClasses">
        {{ $t('download.pin.timeRemaining', { seconds: timerSeconds }) }}
      </div>

      <div class="text-center mb-4">
        <span class="text-sm text-gray-500 dark:text-slate-400">{{ $t('download.pin.attemptsLeft', { count: attemptsLeft }) }}</span>
      </div>

      <div class="mb-6">
        <label class="block text-gray-700 dark:text-slate-300 font-medium mb-2 text-center">{{ $t('download.pin.enterPin') }}</label>
        <div class="flex justify-center">
          <input
            ref="pinInputRef"
            v-model="pinValue"
            type="text"
            maxlength="4"
            autocomplete="off"
            autocorrect="off"
            autocapitalize="off"
            spellcheck="false"
            class="w-40 px-4 py-3 bg-white dark:bg-slate-800 border-2 border-gray-300 dark:border-slate-700 rounded-lg text-gray-900 dark:text-slate-200 text-center text-2xl font-mono tracking-[0.3em] focus:outline-none focus:border-blue-500"
            @keypress.enter="pinValid && handlePinVerify()"
          />
        </div>
      </div>

      <button class="btn-primary" :disabled="!pinValid || verifying" @click="handlePinVerify">
        {{ verifying ? $t('download.pin.verifying') : $t('download.pin.downloadFile') }}
      </button>

      <div v-if="pinError" class="mt-4 p-3 bg-red-500/10 border border-red-500 rounded-lg text-red-600 dark:text-red-400 text-sm text-center">
        {{ pinError }}
      </div>

      <p class="text-yellow-600 dark:text-yellow-400 text-xs text-center mt-4">{{ $t('download.pin.wrongPinWarning') }}</p>
    </section>

    <!-- Download Progress -->
    <section v-if="step === 'downloading'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <h2 class="text-xl font-semibold text-gray-900 dark:text-slate-100 mb-4 text-center">{{ $t('download.pin.downloading') }}</h2>
      <ProgressBar :visible="true" :percent="progress" :text="progressText" />
    </section>

    <!-- Text Result -->
    <section v-if="step === 'text-result'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <div class="bg-green-500/10 border border-green-500 rounded-lg p-4 mb-6">
        <h2 class="text-green-500 font-semibold text-xl mb-2">{{ $t('download.pin.secretText') }}</h2>
        <p class="text-gray-500 dark:text-slate-400 text-sm">{{ $t('download.pin.textDecrypted') }}</p>
      </div>
      <div class="mb-6">
        <textarea :value="decryptedText" readonly rows="10"
          class="w-full px-4 py-3 bg-gray-100 dark:bg-slate-800 border border-gray-300 dark:border-slate-700 rounded-lg text-gray-900 dark:text-slate-200 font-mono text-sm mb-4 resize-vertical focus:outline-none focus:border-blue-500"></textarea>
        <button @click="copyText" class="btn-secondary btn-auto">{{ textCopied ? $t('download.pin.textCopied') : $t('download.pin.copyText') }}</button>
      </div>
      <router-link to="/" class="btn-primary block text-center">{{ $t('download.pin.shareAnotherSecret') }}</router-link>
    </section>

    <!-- File Success -->
    <section v-if="step === 'file-success'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <div class="bg-green-500/10 border border-green-500 rounded-lg p-4 mb-6">
        <h2 class="text-green-500 font-semibold text-xl mb-2">{{ $t('download.pin.downloadComplete') }}</h2>
        <p class="text-gray-700 dark:text-slate-300">{{ $t('download.pin.fileDecrypted') }}</p>
        <p class="text-gray-500 dark:text-slate-400 text-sm mt-1">{{ $t('download.pin.codeInvalid') }}</p>
      </div>
      <router-link to="/" class="btn-primary block text-center">{{ $t('download.pin.shareAnotherFile') }}</router-link>
    </section>

    <!-- Session Timeout -->
    <section v-if="step === 'timeout'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <div class="bg-yellow-500/10 border border-yellow-500 rounded-lg p-4 mb-6">
        <h2 class="text-yellow-600 dark:text-yellow-400 font-semibold text-xl mb-2">{{ $t('download.pin.sessionExpired') }}</h2>
        <p class="text-gray-700 dark:text-slate-300">{{ $t('download.pin.sessionExpiredMsg') }}</p>
      </div>
      <p class="text-gray-600 dark:text-slate-400 text-sm mb-6 text-center">{{ $t('download.pin.enterCodeAgainMsg') }}</p>
      <button @click="handleTryAgain" class="btn-primary">{{ $t('download.pin.enterCodeAgain') }}</button>
    </section>

    <!-- Locked -->
    <section v-if="step === 'locked'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8">
      <div class="bg-red-500/10 border border-red-500 rounded-lg p-4 mb-6">
        <h2 class="text-red-500 font-semibold text-xl mb-2">{{ $t('download.pin.fileLocked') }}</h2>
        <p class="text-gray-700 dark:text-slate-300">{{ $t('download.pin.tooManyAttempts') }}</p>
        <p class="text-gray-600 dark:text-slate-400 text-sm mt-2">{{ $t('download.pin.lockedFor12h') }}</p>
        <p class="text-gray-600 dark:text-slate-400 text-sm">{{ $t('download.pin.tryAgainAt', { time: unlockTime }) }}</p>
      </div>
      <router-link to="/" class="btn-primary block text-center">{{ $t('download.pin.goToHomepage') }}</router-link>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onUnmounted, nextTick } from 'vue';
import { useI18n } from 'vue-i18n';
import PinCodeEntry from '../../components/PinCodeEntry.vue';
import ProgressBar from '../../components/ProgressBar.vue';
import { useRecaptcha } from '../../composables/useRecaptcha.js';
import { useClipboard } from '../../composables/useClipboard.js';
import { API_BASE } from '../../lib/api.js';
import * as CryptoModule from '../../lib/crypto.js';

const { t } = useI18n();
const PIN_REGEX = /^[a-zA-Z0-9]{4}$/;

const { getToken } = useRecaptcha();
const { copied: textCopied, copy } = useClipboard();

const step = ref('code');
const fileId = ref('');
const sessionExpires = ref(0);
const attemptsLeft = ref(3);
const timerSeconds = ref(60);
let timerInterval = null;

const codeEntryRef = ref(null);
const pinInputRef = ref(null);
const codeError = ref('');
const codeLoading = ref(false);
const pinValue = ref('');
const pinError = ref('');
const verifying = ref(false);

const progress = ref(0);
const progressText = ref('');
const decryptedText = ref('');
const unlockTime = ref('');

const pinValid = computed(() => PIN_REGEX.test(pinValue.value));

const timerClasses = computed(() => {
    const base = 'text-center py-3 px-4 rounded-lg mb-4 text-lg font-semibold';
    if (timerSeconds.value > 20) return `${base} bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300`;
    if (timerSeconds.value > 10) return `${base} bg-yellow-100 dark:bg-yellow-900/30 text-yellow-700 dark:text-yellow-300`;
    return `${base} bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-300 animate-pulse`;
});

onUnmounted(() => {
    if (timerInterval) clearInterval(timerInterval);
});

async function handleCodeSubmit(code) {
    codeError.value = '';
    codeLoading.value = true;

    try {
        const recaptchaToken = await getToken('pin_initiate');

        const response = await fetch(`${API_BASE}/pin/initiate`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ file_id: code, recaptcha_token: recaptchaToken }),
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            const errorMsg = errorData.error || `Request failed: ${response.status}`;

            if (errorMsg.toLowerCase().includes('locked')) {
                showLocked(errorData.locked_until);
                return;
            }

            codeError.value = errorMsg;
            return;
        }

        const data = await response.json();
        fileId.value = code;
        sessionExpires.value = data.session_expires;
        attemptsLeft.value = data.attempts_left || 3;

        step.value = 'pin';
        startTimer();
        nextTick(() => { if (pinInputRef.value) pinInputRef.value.focus(); });

    } catch (error) {
        console.error('Code submit error:', error);
        codeError.value = t('download.pin.failedVerify');
    } finally {
        codeLoading.value = false;
    }
}

function startTimer() {
    if (timerInterval) clearInterval(timerInterval);

    const updateTimer = () => {
        const now = Math.floor(Date.now() / 1000);
        const remaining = sessionExpires.value - now;

        if (remaining <= 0) {
            clearInterval(timerInterval);
            timerInterval = null;
            step.value = 'timeout';
            return;
        }

        timerSeconds.value = remaining;
    };

    updateTimer();
    timerInterval = setInterval(updateTimer, 1000);
}

async function handlePinVerify() {
    if (!pinValid.value || verifying.value) return;

    verifying.value = true;
    pinError.value = '';

    try {
        const recaptchaToken = await getToken('pin_verify');

        const response = await fetch(`${API_BASE}/pin/verify`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                file_id: fileId.value,
                pin: pinValue.value,
                recaptcha_token: recaptchaToken,
            }),
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            const errorMsg = errorData.error || `Verification failed: ${response.status}`;

            if (errorMsg.toLowerCase().includes('locked')) {
                clearInterval(timerInterval);
                timerInterval = null;
                showLocked(errorData.locked_until);
                return;
            }

            if (errorMsg.toLowerCase().includes('expired') || errorMsg.toLowerCase().includes('session')) {
                clearInterval(timerInterval);
                timerInterval = null;
                step.value = 'timeout';
                return;
            }

            if (errorData.attempts_left !== undefined) {
                attemptsLeft.value = errorData.attempts_left;
            } else {
                attemptsLeft.value = Math.max(0, attemptsLeft.value - 1);
            }

            pinError.value = errorMsg;
            pinValue.value = '';
            nextTick(() => { if (pinInputRef.value) pinInputRef.value.focus(); });
            return;
        }

        clearInterval(timerInterval);
        timerInterval = null;

        const data = await response.json();
        await handleDecryptAndDownload(pinValue.value, data);

    } catch (error) {
        console.error('PIN verify error:', error);
        pinError.value = t('download.pin.failedVerify');
    } finally {
        verifying.value = false;
    }
}

function hexToUint8Array(hex) {
    const bytes = new Uint8Array(hex.length / 2);
    for (let i = 0; i < hex.length; i += 2) bytes[i / 2] = parseInt(hex.substr(i, 2), 16);
    return bytes;
}

async function handleDecryptAndDownload(pin, data) {
    step.value = 'downloading';

    try {
        progress.value = 5;
        progressText.value = t('download.progress.derivingKey');
        const saltBytes = hexToUint8Array(data.salt);
        const key = await CryptoModule.deriveKeyFromPassword(pin, saltBytes, true);

        if (data.content_type === 'text') {
            progress.value = 30;
            progressText.value = t('download.progress.decryptingText');
            const encryptedBytes = Uint8Array.from(atob(data.encrypted_text), c => c.charCodeAt(0));

            const decryptedData = await CryptoModule.decryptFile(encryptedBytes, key, (p) => {
                progress.value = 30 + p * 0.6;
                progressText.value = t('download.progress.decrypting', { percent: Math.round(p) });
            });

            progress.value = 100;
            progressText.value = t('download.progress.complete');
            decryptedText.value = new TextDecoder().decode(decryptedData);
            step.value = 'text-result';

        } else {
            progress.value = 10;
            progressText.value = t('download.progress.downloadingEncrypted');
            const encryptedData = await downloadEncryptedFile(data.download_url);

            progress.value = 60;
            progressText.value = t('download.progress.decryptingFile', { percent: 0 });
            const decryptedData = await CryptoModule.decryptFile(encryptedData, key, (p) => {
                progress.value = 60 + p * 0.35;
                progressText.value = t('download.progress.decrypting', { percent: Math.round(p) });
            });

            const fileName = data.file_name || 'download';

            if (fileName === 'secret.txt') {
                progress.value = 100;
                progressText.value = t('download.progress.complete');
                decryptedText.value = new TextDecoder().decode(decryptedData);
                step.value = 'text-result';
            } else {
                progress.value = 95;
                progressText.value = t('download.progress.savingFile');
                saveFile(decryptedData, fileName);
                progress.value = 100;
                progressText.value = t('download.progress.downloadComplete');
                step.value = 'file-success';
            }
        }

        confirmDownload(fileId.value);

    } catch (error) {
        console.error('Decrypt/download error:', error);
        step.value = 'code';
        codeError.value = t('download.pin.decryptionFailed');
        if (codeEntryRef.value) codeEntryRef.value.clear();
    }
}

function downloadEncryptedFile(url) {
    return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        xhr.responseType = 'arraybuffer';

        xhr.onprogress = (event) => {
            if (event.lengthComputable) {
                const pct = (event.loaded / event.total) * 100;
                progress.value = 10 + pct * 0.5;
                progressText.value = t('download.progress.downloading', { percent: pct.toFixed(1) });
            }
        };

        xhr.onload = () => {
            if (xhr.status >= 200 && xhr.status < 300) {
                resolve(new Uint8Array(xhr.response));
            } else {
                reject(new Error(`Download failed: ${xhr.status}`));
            }
        };

        xhr.onerror = () => reject(new Error('Network error during download'));
        xhr.ontimeout = () => reject(new Error('Download timeout'));

        xhr.open('GET', url);
        xhr.send();
    });
}

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

function confirmDownload(id) {
    fetch(`${API_BASE}/files/${id}/confirm`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
    }).catch(err => console.warn('Confirm download error:', err));
}

function copyText() {
    copy(decryptedText.value);
}

function showLocked(lockedUntil) {
    step.value = 'locked';
    if (lockedUntil) {
        const d = typeof lockedUntil === 'number' ? new Date(lockedUntil * 1000) : new Date(lockedUntil);
        unlockTime.value = d.toLocaleString();
    } else {
        unlockTime.value = new Date(Date.now() + 12 * 60 * 60 * 1000).toLocaleString();
    }
}

function handleTryAgain() {
    step.value = 'code';
    pinValue.value = '';
    pinError.value = '';
    codeError.value = '';
    if (codeEntryRef.value) codeEntryRef.value.clear();
}
</script>
