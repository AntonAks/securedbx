// Reads the payload stashed by sw.js when Android shares files/text into the PWA.

const STASH = 'share-target-stash';

// Returns { text, title, url, files: File[] } or null. Clears the stash.
export async function takeShare() {
    if (!('caches' in window)) return null;
    const cache = await caches.open(STASH);
    const metaRes = await cache.match('/share-meta');
    if (!metaRes) return null;
    const meta = await metaRes.json();
    const files = await Promise.all(
        meta.files.map(async (f, i) => {
            const res = await cache.match(`/share-file-${i}`);
            return new File([await res.blob()], f.name || `shared-${i}`, { type: f.type });
        })
    );
    await caches.delete(STASH);
    return { text: meta.text, title: meta.title, url: meta.url, files };
}
