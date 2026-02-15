<template>
  <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-8">
    <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-6 text-center">How would you like to share?</h2>
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 justify-center mb-6">
      <button @click="$emit('select', 'pin')" class="p-6 bg-white dark:bg-slate-800 rounded-xl border-2 border-gray-200 dark:border-slate-700 hover:border-blue-500 dark:hover:border-blue-500 transition-all text-center cursor-pointer group">
        <div class="text-3xl mb-3">&#128273;</div>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-slate-100 mb-1">Use PIN Code</h3>
        <p class="text-gray-500 dark:text-slate-400 text-sm mb-2">Share with 6 digits + PIN</p>
        <span class="text-xs text-gray-400 dark:text-slate-500">Easy to type</span>
      </button>
      <button @click="$emit('select', 'link')" class="p-6 bg-white dark:bg-slate-800 rounded-xl border-2 border-gray-200 dark:border-slate-700 hover:border-blue-500 dark:hover:border-blue-500 transition-all text-center cursor-pointer group">
        <div class="text-3xl mb-3">&#128279;</div>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-slate-100 mb-1">Use URL Link</h3>
        <p class="text-gray-500 dark:text-slate-400 text-sm mb-2">Copy & paste a secure URL</p>
        <span class="text-xs text-gray-400 dark:text-slate-500">Works everywhere</span>
      </button>
      <router-link to="/download" class="p-6 bg-white dark:bg-slate-800 rounded-xl border-2 border-gray-200 dark:border-slate-700 hover:border-green-500 dark:hover:border-green-500 transition-all text-center cursor-pointer group">
        <div class="text-3xl mb-3">&#128229;</div>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-slate-100 mb-1">Download</h3>
        <p class="text-gray-500 dark:text-slate-400 text-sm mb-2">Have a code? Get your file</p>
        <span class="text-xs text-gray-400 dark:text-slate-500">Enter 6-digit code</span>
      </router-link>
    </div>
    <div class="text-center">
      <button @click="showHelp = true" class="text-sm text-gray-500 dark:text-slate-400 hover:text-gray-700 dark:hover:text-slate-300 underline">Which method should I use?</button>
    </div>

    <!-- Help Modal -->
    <div v-if="showHelp" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showHelp = false">
      <div class="bg-white dark:bg-slate-900 rounded-xl p-6 max-w-md w-11/12 relative">
        <button @click="showHelp = false" class="absolute top-3 right-3 text-gray-400 hover:text-gray-600 dark:hover:text-slate-300 text-xl">&times;</button>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-slate-100 mb-4">Which method should I use?</h3>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <h4 class="font-medium text-gray-900 dark:text-slate-100 mb-2">Use Link if:</h4>
            <ul class="text-gray-600 dark:text-slate-400 text-sm space-y-1 list-disc ml-4">
              <li>Sharing via email or chat</li>
              <li>Sending to multiple people</li>
              <li>Recipient has a browser open</li>
            </ul>
          </div>
          <div>
            <h4 class="font-medium text-gray-900 dark:text-slate-100 mb-2">Use PIN if:</h4>
            <ul class="text-gray-600 dark:text-slate-400 text-sm space-y-1 list-disc ml-4">
              <li>Sending from phone to computer</li>
              <li>Don't want to copy long URLs</li>
              <li>Want a simple code to type</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref } from 'vue';

defineEmits(['select']);

const showHelp = ref(false);
</script>
