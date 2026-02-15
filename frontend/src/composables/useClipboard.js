import { ref } from 'vue';

/**
 * Composable for clipboard operations with success feedback
 */
export function useClipboard() {
    const copied = ref(false);

    async function copy(text) {
        try {
            await navigator.clipboard.writeText(text);
            copied.value = true;
            setTimeout(() => { copied.value = false; }, 2000);
            return true;
        } catch (err) {
            console.error('Failed to copy:', err);
            return false;
        }
    }

    return { copied, copy };
}
