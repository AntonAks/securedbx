<template>
  <div v-if="visible" class="mt-6">
    <div class="h-2 bg-gray-200 dark:bg-slate-800 rounded-full overflow-hidden mb-2">
      <div class="h-full bg-blue-600 transition-all duration-300" :style="{ width: percent + '%' }"></div>
    </div>
    <p class="text-center text-sm text-gray-600 dark:text-slate-400">{{ text }}</p>
  </div>
</template>

<script setup>
defineProps({
  visible: { type: Boolean, default: false },
  percent: { type: Number, default: 0 },
  text: { type: String, default: '' },
});
</script>
