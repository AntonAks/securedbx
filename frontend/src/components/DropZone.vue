<template>
  <div
    :class="['border-2 border-dashed border-gray-300 dark:border-slate-700 rounded-xl p-12 text-center cursor-pointer transition-all hover:border-blue-500 hover:bg-gray-100 dark:hover:bg-slate-800/50 mb-6', { 'border-blue-500 bg-gray-100 dark:bg-slate-800/80': isDragOver }]"
    @click="$refs.fileInput.click()"
    @dragover="fileDrop.handleDragOver"
    @dragleave="fileDrop.handleDragLeave"
    @drop="fileDrop.handleDrop"
    @keydown.enter.prevent="$refs.fileInput.click()"
    @keydown.space.prevent="$refs.fileInput.click()"
    tabindex="0"
  >
    <svg class="mx-auto mb-4 w-16 h-16 text-gray-400 dark:text-slate-500" fill="currentColor" viewBox="0 0 24 24">
      <path d="M19.35 10.04A7.49 7.49 0 0012 4a7.49 7.49 0 00-7.35 6.04A5.994 5.994 0 000 15c0 3.31 2.69 6 6 6h13c2.76 0 5-2.24 5-5 0-2.64-2.05-4.78-4.65-4.96zM13 13v4h-2v-4H8l4-5 4 5h-3z"/>
    </svg>
    <p class="text-gray-700 dark:text-slate-300 mb-2">{{ $t('components.dropZone.dragDrop') }}</p>
    <p class="text-sm text-gray-500 dark:text-slate-500">{{ $t('components.dropZone.orBrowse') }}</p>
    <input ref="fileInput" type="file" accept="*/*" multiple hidden @change="handleFileChange">
  </div>
</template>

<script setup>
import { useFileDrop } from '../composables/useFileDrop.js';

const emit = defineEmits(['files']);

const fileDrop = useFileDrop((files) => emit('files', files));
const { isDragOver } = fileDrop;

function handleFileChange(e) {
    const files = e.target.files;
    if (files.length > 0) {
        emit('files', files);
    }
    e.target.value = '';
}
</script>
