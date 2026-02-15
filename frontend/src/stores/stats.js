import { defineStore } from 'pinia';
import { ref } from 'vue';
import { API_BASE } from '../lib/api.js';
import { formatBytes } from '../lib/utils.js';

export const useStatsStore = defineStore('stats', () => {
    const text = ref('');
    const loaded = ref(false);

    async function load() {
        if (loaded.value) return;

        try {
            const response = await fetch(`${API_BASE}/stats`, {
                method: 'GET',
                headers: { 'Accept': 'application/json' },
            });

            if (!response.ok) throw new Error(`HTTP ${response.status}`);

            const data = await response.json();
            const downloads = (data.downloads || 0).toLocaleString();
            const size = formatBytes(data.total_bytes || 0);
            text.value = `${downloads} shares \u2022 ${size} transferred`;
            loaded.value = true;
        } catch (error) {
            console.error('Failed to load statistics:', error);
            text.value = '';
        }
    }

    return { text, loaded, load };
});
