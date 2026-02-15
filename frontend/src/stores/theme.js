import { defineStore } from 'pinia';
import { ref } from 'vue';

const THEME_KEY = 'sdbx-theme';
const DARK_CLASS = 'dark';

/**
 * Initialize theme immediately (before app mount) to prevent flash.
 * Called from main.js before createApp.
 */
export function initTheme() {
    const savedTheme = localStorage.getItem(THEME_KEY);
    const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    const shouldBeDark = savedTheme === 'dark' || (!savedTheme && systemPrefersDark) || !savedTheme;

    if (shouldBeDark) {
        document.documentElement.classList.add(DARK_CLASS);
    } else {
        document.documentElement.classList.remove(DARK_CLASS);
    }
}

export const useThemeStore = defineStore('theme', () => {
    const isDark = ref(document.documentElement.classList.contains(DARK_CLASS));

    function toggle() {
        isDark.value = document.documentElement.classList.toggle(DARK_CLASS);
        localStorage.setItem(THEME_KEY, isDark.value ? 'dark' : 'light');
    }

    return { isDark, toggle };
});
