import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';
import i18n from './i18n';
import { initTheme } from './stores/theme';
import './assets/main.css';

// Initialize theme immediately to prevent flash
initTheme();

const app = createApp(App);
app.use(createPinia());
app.use(i18n);
app.use(router);
app.mount('#app');

if (import.meta.env.PROD && 'serviceWorker' in navigator) {
    navigator.serviceWorker.register('/sw.js');
}

// If the app opened outside the share flow (e.g. via a link straight to
// /#/download), drop any stale share stash so plaintext never lingers on disk.
if ('caches' in window && !location.hash.includes('share=1')) {
    caches.delete('share-target-stash');
}
