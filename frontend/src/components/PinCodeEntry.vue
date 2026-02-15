<template>
  <div>
    <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-6 text-center">Download Your File</h2>
    <p class="text-gray-600 dark:text-slate-400 text-center mb-6">Enter your 6-digit code:</p>

    <div class="flex gap-2 sm:gap-3 justify-center mb-6" ref="boxContainer">
      <input
        v-for="(_, i) in 6" :key="i"
        ref="digitInputs"
        type="text"
        class="code-digit w-12 h-14 text-center text-2xl font-semibold bg-white dark:bg-slate-800 border-2 border-gray-300 dark:border-slate-700 rounded-lg focus:border-blue-500 focus:outline-none text-gray-900 dark:text-slate-100"
        maxlength="1"
        inputmode="numeric"
        pattern="[0-9]"
        :aria-label="`Digit ${i + 1}`"
        @input="onInput($event, i)"
        @keydown="onKeydown($event, i)"
        @paste="onPaste($event)"
      />
    </div>

    <button class="btn-primary" :disabled="!isComplete || loading" @click="$emit('submit', code)">
      {{ loading ? 'Verifying...' : 'Continue' }}
    </button>

    <div v-if="error" class="mt-4 p-3 bg-red-500/10 border border-red-500 rounded-lg text-red-600 dark:text-red-400 text-sm text-center">
      {{ error }}
    </div>

    <p class="text-gray-400 dark:text-slate-500 text-xs text-center mt-4">The code looks like: 482973</p>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue';

defineProps({
    error: { type: String, default: '' },
    loading: { type: Boolean, default: false },
});

defineEmits(['submit']);

const digitInputs = ref([]);
const digits = ref(['', '', '', '', '', '']);

const code = computed(() => digits.value.join(''));
const isComplete = computed(() => code.value.length === 6 && /^\d{6}$/.test(code.value));

defineExpose({ code, clear });

onMounted(() => {
    nextTick(() => {
        if (digitInputs.value[0]) digitInputs.value[0].focus();
    });
});

function onInput(event, index) {
    const value = event.target.value;

    if (value && !/^[0-9]$/.test(value)) {
        event.target.value = '';
        digits.value[index] = '';
        return;
    }

    digits.value[index] = value;

    if (value && index < 5) {
        digitInputs.value[index + 1].focus();
    }
}

function onKeydown(event, index) {
    if (event.key === 'Backspace' && !event.target.value && index > 0) {
        digitInputs.value[index - 1].focus();
        digitInputs.value[index - 1].value = '';
        digits.value[index - 1] = '';
    }

    if (event.key === 'Enter' && isComplete.value) {
        // Let the parent handle submit via the button
    }
}

function onPaste(event) {
    event.preventDefault();
    const pastedData = (event.clipboardData || window.clipboardData).getData('text').trim();
    const pastedDigits = pastedData.replace(/[^0-9]/g, '');

    if (pastedDigits.length === 0) return;

    for (let i = 0; i < 6 && i < pastedDigits.length; i++) {
        digits.value[i] = pastedDigits[i];
        if (digitInputs.value[i]) digitInputs.value[i].value = pastedDigits[i];
    }

    const nextEmpty = digits.value.findIndex(d => !d);
    if (nextEmpty >= 0) {
        digitInputs.value[nextEmpty].focus();
    } else {
        digitInputs.value[5].focus();
    }
}

function clear() {
    digits.value = ['', '', '', '', '', ''];
    digitInputs.value.forEach(inp => { if (inp) inp.value = ''; });
    nextTick(() => {
        if (digitInputs.value[0]) digitInputs.value[0].focus();
    });
}
</script>
