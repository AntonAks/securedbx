<template>
  <main class="flex-1 max-w-4xl mx-auto px-6 py-12 w-full">
    <h1 class="text-4xl font-semibold text-gray-900 dark:text-slate-100 mb-8 text-center">Frequently Asked Questions</h1>

    <!-- General Questions -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">General Questions</h2>
      <div class="divide-y divide-gray-200 dark:divide-slate-700">
        <FaqItem question="What is sdbx?">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">
            sdbx is a zero-knowledge file and text sharing service. It lets you securely share files (up to 500MB) and text secrets with end-to-end encryption. By default, content is accessed once and then permanently deleted. Vault and PIN modes also support multi-access sharing — content is automatically deleted when the expiration time (TTL) is reached.
          </p>
        </FaqItem>
        <FaqItem question="Do I need to create an account?">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">
            No! sdbx requires no registration, no email, and no personal information. Simply upload your file or text, get a secure link, and share it. That's it.
          </p>
        </FaqItem>
        <FaqItem question="Is sdbx really free?">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">
            Yes, sdbx is completely free to use. There are no premium tiers, no ads, and no hidden costs.
          </p>
        </FaqItem>
        <FaqItem question="What can I share with sdbx?">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">
            Two types of content: <strong class="text-gray-900 dark:text-slate-100">Files</strong> (any file type up to 500MB) and <strong class="text-gray-900 dark:text-slate-100">Text secrets</strong> (up to 1000 characters for passwords, API keys, notes, etc.).
          </p>
        </FaqItem>
      </div>
    </section>

    <!-- Security & Privacy -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">Security & Privacy</h2>
      <div class="divide-y divide-gray-200 dark:divide-slate-700">
        <FaqItem question='What does "zero-knowledge" mean?'>
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">
            Zero-knowledge means we cannot access your content, even if we wanted to. Your files are encrypted in your browser before upload, and the encryption key is stored only in the URL fragment (#), which is never sent to our servers.
          </p>
        </FaqItem>
        <FaqItem question="What encryption does sdbx use?">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">
            AES-256-GCM with 256-bit keys generated randomly in your browser. Each encryption uses a fresh initialization vector. All communications use TLS 1.2+.
          </p>
        </FaqItem>
        <FaqItem question="Can you see what I upload?">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">
            No. We cannot see your file contents or text secrets. Everything is encrypted before it leaves your browser. We only see encrypted data, file size, and expiration time. In PIN mode, we also store the filename to deliver it to the recipient, but only temporarily — for a maximum of TTL time (24 hours). All stored data is automatically deleted after the first download or when the expiration time (TTL) is reached.
          </p>
        </FaqItem>
        <FaqItem question="Do you log IP addresses or track users?">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">
            No tracking whatsoever. We store only a SHA-256 hash of your IP address for abuse prevention — the original IP cannot be recovered from the hash. We don't use cookies, don't run analytics, and don't track users in any way.
          </p>
        </FaqItem>
        <FaqItem question="Is sdbx open source?">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">
            Yes! Fully open source and available on <a href="https://github.com/antonaks/sdbx" target="_blank" rel="noopener noreferrer" class="text-blue-500 hover:text-blue-400 underline">GitHub</a>. Security researchers are welcome to audit our code.
          </p>
        </FaqItem>
      </div>
    </section>

    <!-- How It Works -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">How It Works</h2>
      <div class="divide-y divide-gray-200 dark:divide-slate-700">
        <FaqItem question="Why can files only be downloaded once by default?">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">
            One-time access is the default security model for URL link sharing. Once the recipient downloads the file, it's permanently deleted from our servers. However, Vault mode and PIN mode offer multi-access sharing — content stays available and is automatically deleted when the expiration time (TTL) is reached.
          </p>
        </FaqItem>
        <FaqItem question="What happens if my download fails due to network issues?">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">
            sdbx uses a two-phase download system: when you click download, the file is reserved for 10 minutes. If your connection drops, you can retry with the same link within that window. After successful download, the file is permanently deleted.
          </p>
        </FaqItem>
        <FaqItem question="How long do files stay available?">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">
            You choose when you upload: 1 hour, 12 hours, or 24 hours. After this time, files are automatically deleted, even if not downloaded.
          </p>
        </FaqItem>
        <FaqItem question="What's the difference between file sharing and text secrets?">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">
            Both are encrypted the same way. Files are for documents/images up to 500MB (downloaded as a file). Text secrets are for passwords/keys up to 1000 characters (displayed directly in browser with copy button).
          </p>
        </FaqItem>
        <FaqItem question="Where is the encryption key stored?">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">
            In the URL fragment (the part after #). URL fragments are never sent to servers by browsers. The key exists only in the shareable link.
          </p>
        </FaqItem>
      </div>
    </section>

    <!-- Limitations & Troubleshooting -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">Limitations & Troubleshooting</h2>
      <div class="divide-y divide-gray-200 dark:divide-slate-700">
        <FaqItem question="What's the maximum file size?">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">
            500 MB. For larger files, consider compressing or splitting them.
          </p>
        </FaqItem>
        <FaqItem question="What browsers are supported?">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">
            All modern browsers: Chrome 37+, Firefox 34+, Safari 11+, Edge 79+. Desktop browsers recommended for large files.
          </p>
        </FaqItem>
        <FaqItem question="I lost my link. Can you recover it?">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">
            Unfortunately, no. Since we never see your encryption key, there's no way to recover your link. Always save links before closing the page.
          </p>
        </FaqItem>
        <FaqItem question='The link says "already downloaded" but I didn&#39;t download anything!'>
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">
            Possible reasons: someone else with the link downloaded it first, the file expired, or link preview features in messaging apps triggered the download.
          </p>
        </FaqItem>
        <FaqItem question="Can I cancel or delete a shared file before it's downloaded?">
          <p class="text-gray-600 dark:text-slate-300 leading-relaxed">
            Currently no manual deletion. Files are automatically deleted after your chosen expiration time.
          </p>
        </FaqItem>
      </div>
    </section>

    <!-- Use Cases -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">Use Cases</h2>
      <div class="divide-y divide-gray-200 dark:divide-slate-700">
        <FaqItem question="What is sdbx perfect for?" title-class="text-green-600 dark:text-green-400">
          <ul class="text-gray-600 dark:text-slate-300 space-y-2 ml-6 list-disc">
            <li>Sharing passwords, API keys, or credentials securely</li>
            <li>Sending confidential documents that shouldn't linger</li>
            <li>Quick secure transfers without creating accounts</li>
            <li>Situations requiring proof that content can only be accessed once</li>
          </ul>
        </FaqItem>
        <FaqItem question="What is sdbx NOT suitable for?" title-class="text-red-500 dark:text-red-400">
          <ul class="text-gray-600 dark:text-slate-300 space-y-2 ml-6 list-disc">
            <li>Long-term file storage (max 7 days expiration)</li>
            <li>Large files</li>
            <li>Illegal content</li>
          </ul>
        </FaqItem>
      </div>
    </section>

    <!-- Still have questions -->
    <section class="bg-gray-50 dark:bg-slate-900/60 rounded-xl border border-gray-200 dark:border-slate-700 p-8 mb-6 text-center">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-slate-100 mb-4">Still have questions?</h2>
      <p class="text-gray-600 dark:text-slate-300 leading-relaxed mb-6">
        Have a question that's not answered here? Open an issue on GitHub.
      </p>
      <a href="https://github.com/antonaks/sdbx/issues" target="_blank" rel="noopener noreferrer" class="btn-outline btn-auto inline-block">
        Ask on GitHub
      </a>
    </section>

    <!-- Back to Home -->
    <div class="text-center mt-8">
      <router-link to="/" class="btn-primary btn-auto inline-block">&larr; Back to Home</router-link>
    </div>
  </main>
</template>

<script setup>
import FaqItem from '../components/FaqItem.vue';
</script>
