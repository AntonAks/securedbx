// Reads the payload stashed by sw.js when Android shares files/text into the PWA.

const STASH = 'share-target-stash';

// Returns { text, title, url, files: File[] } or null.
// ALWAYS clears the stash, even on partial/corrupt payloads — a failed read
// must never leave plaintext secrets sitting in Cache storage.
export async function takeShare() {
    if (!('caches' in window)) return null;
    const cache = await caches.open(STASH);
    try {
        const metaRes = await cache.match('/share-meta');
        if (!metaRes) return null;
        const meta = await metaRes.json();
        const files = (await Promise.all(
            meta.files.map(async (f, i) => {
                const res = await cache.match(`/share-file-${i}`);
                if (!res) return null; // entry evicted/corrupt — skip, don't throw
                return new File([await res.blob()], f.name || `shared-${i}`, { type: f.type });
            })
        )).filter(Boolean);
        return { text: meta.text, title: meta.title, url: meta.url, files };
    } catch {
        return null;
    } finally {
        await caches.delete(STASH);
    }
}
