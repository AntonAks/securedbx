<template>
  <div class="mb-6">
    <label class="block text-gray-700 dark:text-slate-300 font-medium mb-3">Delete after:</label>
    <div class="flex flex-wrap gap-4">
      <label class="flex items-center cursor-pointer">
        <input type="radio" value="1h" v-model="selectedPreset" class="mr-2 text-blue-600 focus:ring-blue-500">
        <span class="text-gray-700 dark:text-slate-300">1 hour</span>
      </label>
      <label class="flex items-center cursor-pointer">
        <input type="radio" value="12h" v-model="selectedPreset" class="mr-2 text-blue-600 focus:ring-blue-500">
        <span class="text-gray-700 dark:text-slate-300">12 hours</span>
      </label>
      <label class="flex items-center cursor-pointer">
        <input type="radio" value="24h" v-model="selectedPreset" class="mr-2 text-blue-600 focus:ring-blue-500">
        <span class="text-gray-700 dark:text-slate-300">24 hours</span>
      </label>
      <label v-if="showCustom" class="flex items-center cursor-pointer">
        <input type="radio" value="custom" v-model="selectedPreset" class="mr-2 text-blue-600 focus:ring-blue-500">
        <span class="text-gray-700 dark:text-slate-300">Custom</span>
      </label>
    </div>
    <div v-if="isCustom" class="mt-3 p-3 bg-gray-100 dark:bg-slate-800 rounded-lg border border-gray-200 dark:border-slate-700">
      <div class="flex items-center gap-2">
        <select v-model="customValue" class="px-3 py-2 bg-white dark:bg-slate-900 border border-gray-300 dark:border-slate-600 rounded text-gray-900 dark:text-slate-200 focus:outline-none focus:border-blue-500">
          <option v-for="v in unitOptions" :key="v" :value="v">{{ v }}</option>
        </select>
        <select v-model="customUnit" class="px-3 py-2 bg-white dark:bg-slate-900 border border-gray-300 dark:border-slate-600 rounded text-gray-900 dark:text-slate-200 focus:outline-none focus:border-blue-500">
          <option value="minutes">minutes</option>
          <option value="hours">hours</option>
          <option value="days">days</option>
        </select>
      </div>
    </div>
    <p class="text-gray-500 dark:text-slate-400 text-sm mt-2">Expires: {{ expirationPreview }}</p>
  </div>
</template>

<script setup>
import { useTtl } from '../composables/useTtl.js';

defineProps({
  showCustom: { type: Boolean, default: true },
});

const {
  selectedPreset,
  customUnit,
  customValue,
  isCustom,
  unitOptions,
  expirationPreview,
  getTtl,
} = useTtl();

defineExpose({ getTtl });
</script>
