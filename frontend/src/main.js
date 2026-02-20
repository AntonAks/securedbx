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
