<template>
  <main class="flex-1 max-w-3xl mx-auto px-6 py-12 w-full">
    <LinkDownload v-if="mode === 'link'" />
    <VaultDownload v-else-if="mode === 'vault'" />
    <PinVerify v-else />

    <!-- Info Section -->
    <section class="mt-12 text-center">
      <h3 class="text-xl font-semibold text-gray-900 dark:text-slate-100 mb-3">{{ $t('download.info.title') }}</h3>
      <p class="text-gray-500 dark:text-slate-400 text-sm leading-relaxed max-w-2xl mx-auto">
        {{ $t('download.info.description') }}
      </p>
    </section>
  </main>
</template>

<script setup>
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import LinkDownload from './download/LinkDownload.vue';
import VaultDownload from './download/VaultDownload.vue';
import PinVerify from './download/PinVerify.vue';

const route = useRoute();

const mode = computed(() => {
    const { id, key, vault } = route.query;
    if (id && vault === '1') return 'vault';
    if (id && key) return 'link';
    return 'pin';
});
</script>
