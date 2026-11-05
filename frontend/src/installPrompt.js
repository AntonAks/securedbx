// Captures Chrome's beforeinstallprompt so the UI can offer PWA install.
// The event only fires when installable (Android/desktop Chrome, not yet
// installed) — so `installPrompt` doubles as the visibility flag.
import { ref } from 'vue';

export const installPrompt = ref(null);

if (typeof window !== 'undefined') {
    window.addEventListener('beforeinstallprompt', (e) => {
        e.preventDefault();
        installPrompt.value = e;
    });
    window.addEventListener('appinstalled', () => {
        installPrompt.value = null;
    });
}

export async function promptInstall() {
    const e = installPrompt.value;
    if (!e) return;
    e.prompt();
    await e.userChoice;
    installPrompt.value = null;
}
