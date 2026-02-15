<template>
  <div v-if="files.length > 0" class="bg-gray-100 dark:bg-slate-800/50 rounded-lg p-4 mb-6 max-h-48 overflow-y-auto">
    <div class="flex justify-between items-center mb-3">
      <span class="text-gray-900 dark:text-slate-100 font-medium">
        Selected: {{ files.length }}/10 files
      </span>
      <span class="text-gray-600 dark:text-slate-400 text-sm">
        Total: {{ formattedTotalSize }}
      </span>
    </div>
    <ul class="space-y-2">
      <li v-for="(file, index) in files" :key="index" class="flex justify-between items-center text-sm py-1 border-b border-gray-200 dark:border-slate-700 last:border-0">
        <span class="text-gray-700 dark:text-slate-300 truncate mr-2" :title="file.name">{{ file.name }}</span>
        <div class="flex items-center gap-3 flex-shrink-0">
          <span class="text-gray-500 dark:text-slate-500">{{ formatFileSize(file.size) }}</span>
          <button type="button" class="text-red-500 hover:text-red-600 dark:text-red-400 dark:hover:text-red-300 p-1" @click="$emit('remove', index)" title="Remove file">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>
      </li>
    </ul>
    <button type="button" class="mt-3 text-sm text-red-500 hover:text-red-600 dark:text-red-400 dark:hover:text-red-300" @click="$emit('clear')">
      Clear All
    </button>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { formatFileSize } from '../lib/utils.js';

const props = defineProps({
  files: { type: Array, required: true },
});

defineEmits(['remove', 'clear']);

const formattedTotalSize = computed(() => {
    const total = props.files.reduce((sum, f) => sum + f.size, 0);
    return formatFileSize(total);
});
</script>
