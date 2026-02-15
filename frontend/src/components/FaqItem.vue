<template>
  <div class="faq-item">
    <button
      class="w-full flex items-center justify-between py-4 text-left cursor-pointer"
      :aria-expanded="open"
      @click="open = !open"
    >
      <span :class="[
        'text-lg font-semibold pr-4',
        titleClass || 'text-gray-900 dark:text-slate-100'
      ]">{{ question }}</span>
      <svg
        :class="['w-5 h-5 shrink-0 text-gray-500 dark:text-slate-400 transition-transform duration-300', { 'rotate-180': open }]"
        fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"
      >
        <path d="M19 9l-7 7-7-7"/>
      </svg>
    </button>
    <div
      class="overflow-hidden transition-all duration-300"
      :style="{ maxHeight: open ? '500px' : '0', paddingTop: open ? '0.5rem' : '0', paddingBottom: open ? '0.5rem' : '0' }"
    >
      <slot />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';

defineProps({
  question: { type: String, required: true },
  titleClass: { type: String, default: '' },
});

const open = ref(false);
</script>
