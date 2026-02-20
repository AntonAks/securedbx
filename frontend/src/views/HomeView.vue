<template>
  <main class="flex-1 max-w-3xl mx-auto px-6 py-12 w-full">
    <!-- Hero -->
    <div class="text-center mb-10">
      <h1 class="text-4xl font-semibold text-gray-900 dark:text-slate-100 mb-3">
        {{ $t('home.hero.title') }}<br>{{ $t('home.hero.titleLine2') }}
      </h1>
      <p class="text-gray-600 dark:text-slate-400">
        {{ $t('home.hero.subtitle') }}
      </p>
      <div class="mt-4 flex justify-center">
        <div class="rounded-full px-4 py-1.5 text-sm font-medium bg-slate-100 text-slate-600 dark:bg-slate-800/80 dark:text-slate-300 backdrop-blur">
          {{ statsStore.text || $t('home.hero.loading') }}
        </div>
      </div>
    </div>

    <!-- Method Selection -->
    <MethodSelector v-if="!selectedMethod && !pinResult" @select="selectMethod" />

    <!-- PIN Upload Section -->
    <section v-if="selectedMethod === 'pin' && !pinResult" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-8">
      <button @click="goBack" class="text-sm text-gray-500 dark:text-slate-400 hover:text-gray-700 dark:hover:text-slate-300 mb-4">{{ $t('home.backToMethods') }}</button>
      <PinUpload @result="onPinResult" />
    </section>

    <!-- Link Upload Section (File / Text / Vault tabs) -->
    <section v-if="selectedMethod === 'link'" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-8">
      <button @click="goBack" class="text-sm text-gray-500 dark:text-slate-400 hover:text-gray-700 dark:hover:text-slate-300 mb-4">{{ $t('home.backToMethods') }}</button>

      <!-- Tab Navigation -->
      <div class="flex justify-center mb-8">
        <div class="inline-flex rounded-full bg-gray-200 dark:bg-slate-800 p-1">
          <button :class="['tab-btn', { active: linkTab === 'file' }]" @click="linkTab = 'file'">
            <svg class="inline w-4 h-4 mr-2 -mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
            </svg>
            {{ $t('home.tabs.file') }}
          </button>
          <button :class="['tab-btn', { active: linkTab === 'text' }]" @click="linkTab = 'text'">
            <svg class="inline w-4 h-4 mr-2 -mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
            </svg>
            {{ $t('home.tabs.text') }}
          </button>
          <button :class="['tab-btn', { active: linkTab === 'vault' }]" @click="linkTab = 'vault'">
            <svg class="inline w-4 h-4 mr-2 -mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
            </svg>
            {{ $t('home.tabs.vault') }}
          </button>
        </div>
      </div>

      <LinkUpload v-if="linkTab === 'file'" />
      <TextUpload v-if="linkTab === 'text'" />
      <VaultUpload v-if="linkTab === 'vault'" />
    </section>

    <!-- PIN Upload Result -->
    <section v-if="pinResult" class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-8">
      <div class="bg-green-500/10 border border-green-500 rounded-lg p-4 mb-6">
        <h2 class="text-green-500 font-semibold text-xl mb-2">{{ $t('home.pinResult.success') }}</h2>
      </div>

      <div class="text-center mb-6">
        <p class="text-gray-600 dark:text-slate-400 mb-2">{{ $t('home.pinResult.downloadCode') }}</p>
        <div class="text-5xl font-bold font-mono tracking-[0.3em] text-gray-900 dark:text-slate-100 py-4 px-6 border-3 border-blue-500 rounded-xl inline-block mb-3">
          {{ pinResult.fileId }}
        </div>
        <div>
          <button @click="copyCode" class="btn-secondary btn-auto text-sm mt-2">
            {{ codeCopied ? $t('home.pinResult.copied') : $t('home.pinResult.copyCode') }}
          </button>
        </div>
      </div>

      <div class="bg-gray-100 dark:bg-slate-800/50 rounded-lg p-4 mb-6 flex items-center justify-center gap-3">
        <span class="text-gray-700 dark:text-slate-300">{{ $t('home.pinResult.yourPin') }}</span>
        <span class="font-mono text-lg tracking-widest text-gray-900 dark:text-slate-100">
          {{ showPin ? pinResult.pin : '****' }}
        </span>
        <button @click="showPin = !showPin" class="text-sm text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300">
          {{ showPin ? $t('home.pinResult.hide') : $t('home.pinResult.show') }}
        </button>
      </div>

      <div class="bg-yellow-500/10 border border-yellow-500 rounded-lg p-4 mb-6">
        <p class="text-yellow-600 dark:text-yellow-400 font-semibold mb-2">{{ $t('home.pinResult.important') }}</p>
        <ul class="text-gray-700 dark:text-slate-300 text-sm space-y-1 ml-5 list-disc">
          <li>{{ $t('home.pinResult.rememberBoth') }}</li>
          <li>{{ $t('home.pinResult.codeWorksOnce') }}</li>
          <li>{{ $t('home.pinResult.expiresIn', { ttl: ttlLabel }) }}</li>
          <li>{{ $t('home.pinResult.wrongPinLockout') }}</li>
        </ul>
      </div>

      <div class="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4 mb-6">
        <p class="text-blue-700 dark:text-blue-300 font-medium mb-2">{{ $t('home.pinResult.toDownload') }}</p>
        <ol class="text-blue-700 dark:text-blue-300 text-sm space-y-1 ml-5 list-decimal">
          <li>{{ $t('home.pinResult.goToDomain', { domain }) }}</li>
          <li>{{ $t('home.pinResult.enterCode', { code: pinResult.fileId }) }}</li>
          <li>{{ $t('home.pinResult.enterYourPin') }}</li>
        </ol>
      </div>

      <button @click="uploadAnother" class="btn-primary">{{ $t('home.pinResult.uploadAnother') }}</button>
    </section>

    <!-- Features Section -->
    <section class="text-center">
      <h3 class="text-xl font-semibold text-gray-900 dark:text-slate-100 mb-8">{{ $t('home.features.heading') }}</h3>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <div class="flex flex-col items-center">
          <svg class="w-8 h-8 text-blue-600 dark:text-blue-500 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
          </svg>
          <h4 class="text-gray-900 dark:text-slate-100 font-medium mb-2">{{ $t('home.features.encryption') }}</h4>
          <p class="text-gray-600 dark:text-slate-400 text-sm">{{ $t('home.features.encryptionDesc') }}</p>
        </div>
        <div class="flex flex-col items-center">
          <svg class="w-8 h-8 text-blue-600 dark:text-blue-500 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/>
          </svg>
          <h4 class="text-gray-900 dark:text-slate-100 font-medium mb-2">{{ $t('home.features.zeroKnowledge') }}</h4>
          <p class="text-gray-600 dark:text-slate-400 text-sm">{{ $t('home.features.zeroKnowledgeDesc') }}</p>
        </div>
        <div class="flex flex-col items-center">
          <svg class="w-8 h-8 text-blue-600 dark:text-blue-500 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"/>
          </svg>
          <h4 class="text-gray-900 dark:text-slate-100 font-medium mb-2">{{ $t('home.features.noTracking') }}</h4>
          <p class="text-gray-600 dark:text-slate-400 text-sm">{{ $t('home.features.noTrackingDesc') }}</p>
        </div>
        <div class="flex flex-col items-center">
          <svg class="w-8 h-8 text-blue-600 dark:text-blue-500 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path d="M13 10V3L4 14h7v7l9-11h-7z"/>
          </svg>
          <h4 class="text-gray-900 dark:text-slate-100 font-medium mb-2">{{ $t('home.features.oneTime') }}</h4>
          <p class="text-gray-600 dark:text-slate-400 text-sm">{{ $t('home.features.oneTimeDesc') }}</p>
        </div>
      </div>
    </section>

    <!-- This Service Is For You If... -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mt-8">
      <h3 class="text-xl font-semibold text-gray-900 dark:text-slate-100 mb-6 text-center">{{ $t('home.forYou.heading') }}</h3>
      <ul class="text-gray-600 dark:text-slate-300 space-y-3">
        <li class="flex items-start">
          <span class="text-blue-600 dark:text-blue-500 mr-3">&#10003;</span>
          {{ $t('home.forYou.reason1') }}
        </li>
        <li class="flex items-start">
          <span class="text-blue-600 dark:text-blue-500 mr-3">&#10003;</span>
          {{ $t('home.forYou.reason2') }}
        </li>
        <li class="flex items-start">
          <span class="text-blue-600 dark:text-blue-500 mr-3">&#10003;</span>
          {{ $t('home.forYou.reason3') }}
        </li>
        <li class="flex items-start">
          <span class="text-blue-600 dark:text-blue-500 mr-3">&#10003;</span>
          {{ $t('home.forYou.reason4') }}
        </li>
        <li class="flex items-start">
          <span class="text-blue-600 dark:text-blue-500 mr-3">&#10003;</span>
          {{ $t('home.forYou.reason5') }}
        </li>
      </ul>
    </section>
  </main>
</template>

<script setup>
import { ref, computed, onMounted, inject, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useStatsStore } from '../stores/stats.js';
import { useClipboard } from '../composables/useClipboard.js';
import MethodSelector from '../components/MethodSelector.vue';
import PinUpload from './home/PinUpload.vue';
import LinkUpload from './home/LinkUpload.vue';
import TextUpload from './home/TextUpload.vue';
import VaultUpload from './home/VaultUpload.vue';

const { t } = useI18n();

const statsStore = useStatsStore();
const { copied: codeCopied, copy } = useClipboard();

const selectedMethod = ref(null);
const linkTab = ref('file');
const pinResult = ref(null);
const showPin = ref(false);

const domain = computed(() => window.location.host);
const ttlLabel = computed(() => {
  if (!pinResult.value) return '';
  const ttlMap = { '1h': t('upload.pin.ttl1h'), '12h': t('upload.pin.ttl12h'), '24h': t('upload.pin.ttl24h') };
  return ttlMap[pinResult.value.ttl] || pinResult.value.ttl;
});

const homeResetKey = inject('homeResetKey');

onMounted(() => statsStore.load());

watch(homeResetKey, () => {
    selectedMethod.value = null;
    pinResult.value = null;
    showPin.value = false;
    linkTab.value = 'file';
});

function selectMethod(method) {
    selectedMethod.value = method;
}

function goBack() {
    selectedMethod.value = null;
}

function onPinResult(result) {
    pinResult.value = result;
}

function copyCode() {
    if (pinResult.value) copy(pinResult.value.fileId);
}

function uploadAnother() {
    pinResult.value = null;
    selectedMethod.value = null;
    showPin.value = false;
}
</script>
