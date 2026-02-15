<template>
  <div ref="containerRef" class="bg-white p-4 rounded-lg inline-block">
    <canvas v-if="!error" ref="canvasRef"></canvas>
    <div v-else class="flex items-center justify-center text-center p-4" style="width: 256px; height: 256px;">
      <p class="text-gray-500 dark:text-slate-400 text-sm">
        URL too long for QR code.<br>Use "Copy Link" instead.
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import QRCode from 'qrcode';

const props = defineProps({
    value: { type: String, required: true },
    size: { type: Number, default: 256 },
});

const canvasRef = ref(null);
const error = ref(false);

defineExpose({ canvasRef });

onMounted(() => render());

watch(() => props.value, () => render());

async function render() {
    if (!canvasRef.value || !props.value) return;

    try {
        await QRCode.toCanvas(canvasRef.value, props.value, {
            width: props.size,
            margin: 0,
            errorCorrectionLevel: 'L',
            color: { dark: '#000000', light: '#ffffff' },
        });
        error.value = false;
    } catch (err) {
        console.warn('QR code generation failed:', err.message);
        error.value = true;
    }
}
</script>
