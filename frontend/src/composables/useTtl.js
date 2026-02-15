import { ref, computed, watch } from 'vue';

const TTL_UNIT_OPTIONS = {
    minutes: [5, 10, 20, 30, 40, 50],
    hours: [2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 24],
    days: [1, 2, 3, 4, 5, 6, 7],
};

/**
 * Composable for TTL (time-to-live) selection
 */
export function useTtl() {
    const selectedPreset = ref('1h');
    const customUnit = ref('hours');
    const customValue = ref(2);

    const isCustom = computed(() => selectedPreset.value === 'custom');

    const unitOptions = computed(() => TTL_UNIT_OPTIONS[customUnit.value] || TTL_UNIT_OPTIONS.hours);

    /**
     * Get TTL in minutes
     * @returns {number|string} Minutes for custom, preset string otherwise
     */
    function getTtl() {
        if (selectedPreset.value === 'custom') {
            const val = parseInt(customValue.value, 10);
            switch (customUnit.value) {
                case 'hours': return val * 60;
                case 'days': return val * 60 * 24;
                default: return val; // minutes
            }
        }
        return selectedPreset.value;
    }

    /**
     * Get TTL in minutes (always numeric)
     * @returns {number}
     */
    function getTtlMinutes() {
        const ttl = getTtl();
        if (typeof ttl === 'number') return ttl;
        const presetMinutes = { '1h': 60, '12h': 720, '24h': 1440 };
        return presetMinutes[ttl] || 60;
    }

    const expirationPreview = computed(() => {
        const minutes = getTtlMinutes();
        const expirationDate = new Date(Date.now() + minutes * 60 * 1000);
        return expirationDate.toLocaleString(undefined, {
            weekday: 'short',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    });

    // Reset custom value when unit changes to a valid option
    watch(customUnit, () => {
        const opts = TTL_UNIT_OPTIONS[customUnit.value];
        if (opts && !opts.includes(parseInt(customValue.value, 10))) {
            customValue.value = opts[0];
        }
    });

    return {
        selectedPreset,
        customUnit,
        customValue,
        isCustom,
        unitOptions,
        expirationPreview,
        getTtl,
        getTtlMinutes,
    };
}
