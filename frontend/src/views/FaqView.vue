<template>
  <main class="flex-1 max-w-4xl mx-auto px-6 py-12 w-full">
    <h1 class="text-4xl font-semibold text-gray-900 dark:text-slate-100 mb-8 text-center">{{ $t('faq.title') }}</h1>

    <!-- General Questions -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">{{ $t('faq.general.heading') }}</h2>
      <div class="divide-y divide-gray-200 dark:divide-slate-700">
        <FaqItem :question="$t('faq.general.q1')">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">{{ $t('faq.general.a1') }}</p>
        </FaqItem>
        <FaqItem :question="$t('faq.general.q2')">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">{{ $t('faq.general.a2') }}</p>
        </FaqItem>
        <FaqItem :question="$t('faq.general.q3')">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">{{ $t('faq.general.a3') }}</p>
        </FaqItem>
        <FaqItem :question="$t('faq.general.q4')">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed" v-html="$t('faq.general.a4html')"></p>
        </FaqItem>
      </div>
    </section>

    <!-- Security & Privacy -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">{{ $t('faq.security.heading') }}</h2>
      <div class="divide-y divide-gray-200 dark:divide-slate-700">
        <FaqItem :question="$t('faq.security.q1')">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">{{ $t('faq.security.a1') }}</p>
        </FaqItem>
        <FaqItem :question="$t('faq.security.q2')">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">{{ $t('faq.security.a2') }}</p>
        </FaqItem>
        <FaqItem :question="$t('faq.security.q3')">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">{{ $t('faq.security.a3') }}</p>
        </FaqItem>
        <FaqItem :question="$t('faq.security.q4')">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">{{ $t('faq.security.a4') }}</p>
        </FaqItem>
        <FaqItem :question="$t('faq.security.q5')">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed" v-html="$t('faq.security.a5html')"></p>
        </FaqItem>
      </div>
    </section>

    <!-- How It Works -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">{{ $t('faq.howItWorks.heading') }}</h2>
      <div class="divide-y divide-gray-200 dark:divide-slate-700">
        <FaqItem :question="$t('faq.howItWorks.q1')">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">{{ $t('faq.howItWorks.a1') }}</p>
        </FaqItem>
        <FaqItem :question="$t('faq.howItWorks.q2')">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">{{ $t('faq.howItWorks.a2') }}</p>
        </FaqItem>
        <FaqItem :question="$t('faq.howItWorks.q3')">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">{{ $t('faq.howItWorks.a3') }}</p>
        </FaqItem>
        <FaqItem :question="$t('faq.howItWorks.q4')">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">{{ $t('faq.howItWorks.a4') }}</p>
        </FaqItem>
        <FaqItem :question="$t('faq.howItWorks.q5')">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">{{ $t('faq.howItWorks.a5') }}</p>
        </FaqItem>
      </div>
    </section>

    <!-- Limitations & Troubleshooting -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">{{ $t('faq.limitations.heading') }}</h2>
      <div class="divide-y divide-gray-200 dark:divide-slate-700">
        <FaqItem :question="$t('faq.limitations.q1')">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">{{ $t('faq.limitations.a1') }}</p>
        </FaqItem>
        <FaqItem :question="$t('faq.limitations.q2')">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">{{ $t('faq.limitations.a2') }}</p>
        </FaqItem>
        <FaqItem :question="$t('faq.limitations.q3')">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">{{ $t('faq.limitations.a3') }}</p>
        </FaqItem>
        <FaqItem :question="$t('faq.limitations.q4')">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">{{ $t('faq.limitations.a4') }}</p>
        </FaqItem>
        <FaqItem :question="$t('faq.limitations.q5')">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">{{ $t('faq.limitations.a5') }}</p>
        </FaqItem>
      </div>
    </section>

    <!-- Use Cases -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">{{ $t('faq.useCases.heading') }}</h2>
      <div class="divide-y divide-gray-200 dark:divide-slate-700">
        <FaqItem :question="$t('faq.useCases.goodTitle')" title-class="text-green-600 dark:text-green-400">
          <ul class="text-gray-600 dark:text-slate-300 space-y-2 ml-6 list-disc">
            <li>{{ $t('faq.useCases.good1') }}</li>
            <li>{{ $t('faq.useCases.good2') }}</li>
            <li>{{ $t('faq.useCases.good3') }}</li>
            <li>{{ $t('faq.useCases.good4') }}</li>
          </ul>
        </FaqItem>
        <FaqItem :question="$t('faq.useCases.badTitle')" title-class="text-red-500 dark:text-red-400">
          <ul class="text-gray-600 dark:text-slate-300 space-y-2 ml-6 list-disc">
            <li>{{ $t('faq.useCases.bad1') }}</li>
            <li>{{ $t('faq.useCases.bad2') }}</li>
            <li>{{ $t('faq.useCases.bad3') }}</li>
          </ul>
        </FaqItem>
      </div>
    </section>

    <!-- Still have questions -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6 text-center">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">{{ $t('faq.stillHaveQuestions.heading') }}</h2>
      <p class="text-gray-600 dark:text-slate-300 leading-relaxed mb-6">{{ $t('faq.stillHaveQuestions.description') }}</p>
      <a href="https://github.com/antonaks/sdbx/issues" target="_blank" rel="noopener noreferrer" class="btn-outline btn-auto inline-block">
        {{ $t('faq.stillHaveQuestions.askOnGithub') }}
      </a>
    </section>

    <!-- Back to Home -->
    <div class="text-center mt-8">
      <router-link to="/" class="btn-primary btn-auto inline-block">{{ $t('faq.backToHome') }}</router-link>
    </div>
  </main>
</template>

<script setup>
import { useI18n } from 'vue-i18n';
import FaqItem from '../components/FaqItem.vue';

useI18n();
</script>
