<template>
  <main class="flex-1 max-w-4xl mx-auto px-6 py-12 w-full">
    <h1 class="text-4xl font-semibold text-gray-900 dark:text-slate-100 mb-8 text-center">{{ $t('about.title') }}</h1>

    <!-- What is sdbx -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">{{ $t('about.whatIs.heading') }}</h2>
      <p class="text-gray-600 dark:text-slate-300 leading-relaxed mb-4" v-html="$t('about.whatIs.p1html')"></p>
      <p class="text-gray-600 dark:text-slate-300 leading-relaxed" v-html="$t('about.whatIs.p2html')"></p>
    </section>

    <!-- Our Mission -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">{{ $t('about.mission.heading') }}</h2>
      <p class="text-gray-600 dark:text-slate-300 leading-relaxed mb-4">{{ $t('about.mission.p1') }}</p>
      <p class="text-gray-600 dark:text-slate-300 leading-relaxed mb-4">{{ $t('about.mission.p2') }}</p>
      <ul class="text-gray-600 dark:text-slate-300 space-y-2 ml-6 list-disc">
        <li>{{ $t('about.mission.item1') }}</li>
        <li>{{ $t('about.mission.item2') }}</li>
        <li>{{ $t('about.mission.item3') }}</li>
        <li>{{ $t('about.mission.item4') }}</li>
      </ul>
    </section>

    <!-- How It Works -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-6">{{ $t('about.howItWorks.heading') }}</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div class="flex flex-col items-start">
          <div class="flex items-center gap-3 mb-3">
            <svg class="w-6 h-6 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
            </svg>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-slate-100">{{ $t('about.howItWorks.encryption') }}</h3>
          </div>
          <p class="text-gray-500 dark:text-slate-400 text-sm leading-relaxed">{{ $t('about.howItWorks.encryptionDesc') }}</p>
        </div>

        <div class="flex flex-col items-start">
          <div class="flex items-center gap-3 mb-3">
            <svg class="w-6 h-6 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/>
            </svg>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-slate-100">{{ $t('about.howItWorks.zeroKnowledge') }}</h3>
          </div>
          <p class="text-gray-500 dark:text-slate-400 text-sm leading-relaxed">{{ $t('about.howItWorks.zeroKnowledgeDesc') }}</p>
        </div>

        <div class="flex flex-col items-start">
          <div class="flex items-center gap-3 mb-3">
            <svg class="w-6 h-6 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"/>
            </svg>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-slate-100">{{ $t('about.howItWorks.noTracking') }}</h3>
          </div>
          <p class="text-gray-500 dark:text-slate-400 text-sm leading-relaxed">{{ $t('about.howItWorks.noTrackingDesc') }}</p>
        </div>

        <div class="flex flex-col items-start">
          <div class="flex items-center gap-3 mb-3">
            <svg class="w-6 h-6 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path d="M13 10V3L4 14h7v7l9-11h-7z"/>
            </svg>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-slate-100">{{ $t('about.howItWorks.selfDestruct') }}</h3>
          </div>
          <p class="text-gray-500 dark:text-slate-400 text-sm leading-relaxed">{{ $t('about.howItWorks.selfDestructDesc') }}</p>
        </div>

        <div class="flex flex-col items-start">
          <div class="flex items-center gap-3 mb-3">
            <svg class="w-6 h-6 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path d="M12 11c0-1.1.9-2 2-2s2 .9 2 2-.9 2-2 2-2-.9-2-2zm-6 0c0-1.1.9-2 2-2s2 .9 2 2-.9 2-2 2-2-.9-2-2zm12 0c0-1.1.9-2 2-2s2 .9 2 2-.9 2-2 2-2-.9-2-2z"/>
            </svg>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-slate-100">{{ $t('about.howItWorks.pinCode') }}</h3>
          </div>
          <p class="text-gray-500 dark:text-slate-400 text-sm leading-relaxed">{{ $t('about.howItWorks.pinCodeDesc') }}</p>
        </div>

        <div class="flex flex-col items-start">
          <div class="flex items-center gap-3 mb-3">
            <svg class="w-6 h-6 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path d="M8 11V7a4 4 0 118 0m-4 8v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2z"/>
            </svg>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-slate-100">{{ $t('about.howItWorks.vaultMode') }}</h3>
          </div>
          <p class="text-gray-500 dark:text-slate-400 text-sm leading-relaxed">{{ $t('about.howItWorks.vaultModeDesc') }}</p>
        </div>
      </div>
    </section>

    <!-- Security Features -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">{{ $t('about.security.heading') }}</h2>
      <ul class="text-gray-600 dark:text-slate-300 space-y-2 ml-6 list-disc">
        <li><strong class="text-gray-900 dark:text-slate-100">{{ $t('about.security.e2e') }}</strong> {{ $t('about.security.e2eDesc') }}</li>
        <li><strong class="text-gray-900 dark:text-slate-100">{{ $t('about.security.https') }}</strong> {{ $t('about.security.httpsDesc') }}</li>
        <li><strong class="text-gray-900 dark:text-slate-100">{{ $t('about.security.bot') }}</strong> {{ $t('about.security.botDesc') }}</li>
        <li><strong class="text-gray-900 dark:text-slate-100">{{ $t('about.security.twoPhase') }}</strong> {{ $t('about.security.twoPhaseDesc') }}</li>
        <li><strong class="text-gray-900 dark:text-slate-100">{{ $t('about.security.origin') }}</strong> {{ $t('about.security.originDesc') }}</li>
        <li><strong class="text-gray-900 dark:text-slate-100">{{ $t('about.security.pbkdf2') }}</strong> {{ $t('about.security.pbkdf2Desc') }}</li>
        <li><strong class="text-gray-900 dark:text-slate-100">{{ $t('about.security.cleanup') }}</strong> {{ $t('about.security.cleanupDesc') }}</li>
      </ul>
    </section>

    <!-- Technology Stack -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">{{ $t('about.techStack.heading') }}</h2>
      <p class="text-gray-600 dark:text-slate-300 leading-relaxed mb-4">{{ $t('about.techStack.intro') }}</p>
      <ul class="text-gray-600 dark:text-slate-300 space-y-2 ml-6 list-disc">
        <li><strong class="text-gray-900 dark:text-slate-100">{{ $t('about.techStack.frontend') }}</strong> {{ $t('about.techStack.frontendDesc') }}</li>
        <li><strong class="text-gray-900 dark:text-slate-100">{{ $t('about.techStack.backend') }}</strong> {{ $t('about.techStack.backendDesc') }}</li>
        <li><strong class="text-gray-900 dark:text-slate-100">{{ $t('about.techStack.storage') }}</strong> {{ $t('about.techStack.storageDesc') }}</li>
        <li><strong class="text-gray-900 dark:text-slate-100">{{ $t('about.techStack.cdn') }}</strong> {{ $t('about.techStack.cdnDesc') }}</li>
        <li><strong class="text-gray-900 dark:text-slate-100">{{ $t('about.techStack.infra') }}</strong> {{ $t('about.techStack.infraDesc') }}</li>
      </ul>
    </section>

    <!-- Limitations -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">{{ $t('about.limitations.heading') }}</h2>
      <p class="text-gray-600 dark:text-slate-300 leading-relaxed mb-4">{{ $t('about.limitations.intro') }}</p>
      <ul class="text-gray-600 dark:text-slate-300 space-y-2 ml-6 list-disc">
        <li>{{ $t('about.limitations.maxSize') }}</li>
        <li>{{ $t('about.limitations.expiry') }}</li>
        <li>{{ $t('about.limitations.accessModes') }}</li>
        <li>{{ $t('about.limitations.maxText') }}</li>
        <li>{{ $t('about.limitations.maxFiles') }}</li>
        <li>{{ $t('about.limitations.desktop') }}</li>
      </ul>
    </section>

    <!-- Privacy Policy -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">{{ $t('about.privacy.heading') }}</h2>
      <p class="text-gray-600 dark:text-slate-300 leading-relaxed mb-4" v-html="$t('about.privacy.simplehtml')"></p>
      <ul class="text-gray-600 dark:text-slate-300 space-y-2 ml-6 list-disc mb-4">
        <li>{{ $t('about.privacy.item1') }}</li>
        <li>{{ $t('about.privacy.item2') }}</li>
        <li>{{ $t('about.privacy.item3') }}</li>
        <li>{{ $t('about.privacy.item4') }}</li>
        <li>{{ $t('about.privacy.item5') }}</li>
        <li>{{ $t('about.privacy.item6') }}</li>
      </ul>
      <p class="text-gray-600 dark:text-slate-300 leading-relaxed">{{ $t('about.privacy.detail') }}</p>
    </section>

    <!-- Open Source -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">{{ $t('about.openSource.heading') }}</h2>
      <p class="text-gray-600 dark:text-slate-300 leading-relaxed mb-6">{{ $t('about.openSource.description') }}</p>
      <div class="text-center">
        <a href="https://github.com/antonaks/sdbx" target="_blank" rel="noopener noreferrer" class="btn-outline btn-auto inline-block">
          {{ $t('about.openSource.viewSource') }}
        </a>
      </div>
    </section>

    <!-- Disclaimer -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">{{ $t('about.disclaimer.heading') }}</h2>
      <p class="text-gray-600 dark:text-slate-300 leading-relaxed mb-4">{{ $t('about.disclaimer.p1') }}</p>
      <p class="text-gray-600 dark:text-slate-300 leading-relaxed">{{ $t('about.disclaimer.p2') }}</p>
    </section>

    <!-- Questions -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6 text-center">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">{{ $t('about.questions.heading') }}</h2>
      <p class="text-gray-600 dark:text-slate-300 leading-relaxed mb-6">{{ $t('about.questions.description') }}</p>
      <a href="https://github.com/antonaks/sdbx/issues" target="_blank" rel="noopener noreferrer" class="btn-outline btn-auto inline-block">
        {{ $t('about.questions.reportIssue') }}
      </a>
    </section>

    <!-- Back to Home -->
    <div class="text-center mt-8">
      <router-link to="/" class="btn-primary btn-auto inline-block">{{ $t('about.backToHome') }}</router-link>
    </div>
  </main>
</template>

<script setup>
import { useI18n } from 'vue-i18n';

useI18n();
</script>
