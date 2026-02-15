<template>
  <footer class="border-t border-gray-200 dark:border-slate-800 py-8 text-center">
    <p class="text-gray-500 dark:text-slate-400 text-sm mb-2">
      &copy; {{ year }} sdbx - Built with privacy in mind
    </p>
    <p class="text-sm">
      <router-link to="/about" class="text-blue-500 hover:text-blue-400">About</router-link>
      <span class="text-gray-400 dark:text-slate-600 mx-2">|</span>
      <router-link to="/faq" class="text-blue-500 hover:text-blue-400">FAQ</router-link>
      <span class="text-gray-400 dark:text-slate-600 mx-2">|</span>
      <a href="https://github.com/antonaks/sdbx" target="_blank" rel="noopener noreferrer" class="text-blue-500 hover:text-blue-400">Source Code</a>
    </p>
  </footer>
</template>

<script setup>
const year = new Date().getFullYear();
</script>
