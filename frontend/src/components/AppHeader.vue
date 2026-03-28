<template>
  <header class="bg-gray-100 dark:bg-slate-900 border-b border-gray-200 dark:border-slate-800">
    <div class="max-w-4xl mx-auto px-6 py-4 flex items-center justify-between">
      <a @click.prevent="goHome" href="/" class="flex items-center gap-3 cursor-pointer">
        <img src="/img/logo.png" alt="SecureDBX logo" class="w-14 h-14">
        <span class="text-xl font-semibold text-gray-900 dark:text-slate-100">SecureDBX</span>
      </a>

      <!-- Desktop nav -->
      <nav class="hidden sm:flex items-center gap-6">
        <router-link
          v-if="route.path !== '/'"
          to="/"
          class="text-sm text-gray-600 dark:text-slate-400 hover:text-gray-900 dark:hover:text-slate-200 transition"
        >{{ $t('nav.home') }}</router-link>
        <router-link
          to="/about"
          :class="route.path === '/about'
            ? 'text-sm text-blue-600 dark:text-blue-500'
            : 'text-sm text-gray-600 dark:text-slate-400 hover:text-gray-900 dark:hover:text-slate-200 transition'"
        >{{ $t('nav.about') }}</router-link>
        <router-link
          to="/faq"
          :class="route.path === '/faq'
            ? 'text-sm text-blue-600 dark:text-blue-500'
            : 'text-sm text-gray-600 dark:text-slate-400 hover:text-gray-900 dark:hover:text-slate-200 transition'"
        >{{ $t('nav.faq') }}</router-link>

        <div class="flex items-center gap-1">
          <button
            @click="switchLocale('en')"
            :class="['text-xl leading-none transition-opacity', locale === 'en' ? 'opacity-100' : 'opacity-30 hover:opacity-60']"
            :aria-label="$t('language.en')"
            title="English"
          >🇬🇧</button>
          <button
            @click="switchLocale('uk')"
            :class="['text-xl leading-none transition-opacity', locale === 'uk' ? 'opacity-100' : 'opacity-30 hover:opacity-60']"
            :aria-label="$t('language.uk')"
            title="Солов'їна"
          >🇺🇦</button>
        </div>

        <button
          @click="themeStore.toggle()"
          class="p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-slate-800 transition-colors"
          :aria-label="themeStore.isDark ? $t('nav.switchToLight') : $t('nav.switchToDark')"
        >
          <svg v-if="themeStore.isDark" class="w-5 h-5 text-slate-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="4"/>
            <path d="M12 2v2m0 16v2M4.93 4.93l1.41 1.41m11.32 11.32l1.41 1.41M2 12h2m16 0h2M4.93 19.07l1.41-1.41m11.32-11.32l1.41-1.41"/>
          </svg>
          <svg v-else class="w-5 h-5 text-slate-700" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/>
          </svg>
        </button>
      </nav>

      <!-- Mobile: theme toggle always visible + hamburger -->
      <div class="flex items-center gap-2 sm:hidden">
        <button
          @click="themeStore.toggle()"
          class="p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-slate-800 transition-colors"
          :aria-label="themeStore.isDark ? $t('nav.switchToLight') : $t('nav.switchToDark')"
        >
          <svg v-if="themeStore.isDark" class="w-5 h-5 text-slate-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="4"/>
            <path d="M12 2v2m0 16v2M4.93 4.93l1.41 1.41m11.32 11.32l1.41 1.41M2 12h2m16 0h2M4.93 19.07l1.41-1.41m11.32-11.32l1.41-1.41"/>
          </svg>
          <svg v-else class="w-5 h-5 text-slate-700" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/>
          </svg>
        </button>

        <button
          @click="menuOpen = !menuOpen"
          class="p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-slate-800 transition-colors"
          aria-label="Toggle menu"
        >
          <svg v-if="!menuOpen" class="w-5 h-5 text-gray-700 dark:text-slate-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" d="M4 6h16M4 12h16M4 18h16"/>
          </svg>
          <svg v-else class="w-5 h-5 text-gray-700 dark:text-slate-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- Mobile dropdown -->
    <div v-if="menuOpen" class="sm:hidden border-t border-gray-200 dark:border-slate-800 bg-gray-100 dark:bg-slate-900 px-6 py-4 flex flex-col gap-4">
      <router-link
        v-if="route.path !== '/'"
        to="/"
        @click="menuOpen = false"
        class="text-sm text-gray-600 dark:text-slate-400 hover:text-gray-900 dark:hover:text-slate-200 transition"
      >{{ $t('nav.home') }}</router-link>
      <router-link
        to="/about"
        @click="menuOpen = false"
        :class="route.path === '/about'
          ? 'text-sm text-blue-600 dark:text-blue-500'
          : 'text-sm text-gray-600 dark:text-slate-400 hover:text-gray-900 dark:hover:text-slate-200 transition'"
      >{{ $t('nav.about') }}</router-link>
      <router-link
        to="/faq"
        @click="menuOpen = false"
        :class="route.path === '/faq'
          ? 'text-sm text-blue-600 dark:text-blue-500'
          : 'text-sm text-gray-600 dark:text-slate-400 hover:text-gray-900 dark:hover:text-slate-200 transition'"
      >{{ $t('nav.faq') }}</router-link>

      <div class="flex items-center gap-2">
        <button
          @click="switchLocale('en')"
          :class="['text-xl leading-none transition-opacity', locale === 'en' ? 'opacity-100' : 'opacity-30 hover:opacity-60']"
          :aria-label="$t('language.en')"
        >🇬🇧</button>
        <button
          @click="switchLocale('uk')"
          :class="['text-xl leading-none transition-opacity', locale === 'uk' ? 'opacity-100' : 'opacity-30 hover:opacity-60']"
          :aria-label="$t('language.uk')"
        >🇺🇦</button>
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref, watch } from 'vue';
import { inject } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useThemeStore } from '../stores/theme.js';
import { setLocale } from '../i18n';

const route = useRoute();
const router = useRouter();
const themeStore = useThemeStore();
const triggerHomeReset = inject('triggerHomeReset');
const { locale } = useI18n();
const menuOpen = ref(false);

// Close menu on navigation
watch(() => route.path, () => { menuOpen.value = false; });

function goHome() {
  triggerHomeReset();
  router.push('/');
}

function switchLocale(loc) {
  setLocale(loc);
}
</script>