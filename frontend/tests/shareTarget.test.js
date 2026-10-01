/**
 * shareTarget — round-trip of the share stash written by sw.js.
 *
 * Run with: npx vitest run tests/shareTarget.test.js
 */

import { describe, it, expect, beforeEach } from 'vitest';
import { takeShare } from '../src/shareTarget.js';

// Minimal in-memory Cache API stub (happy-dom has none)
function makeCaches() {
    const stores = new Map();
    return {
        async open(name) {
            if (!stores.has(name)) stores.set(name, new Map());
            const store = stores.get(name);
            return {
                async put(key, response) { store.set(key, response); },
                async match(key) { return store.get(key) ?? undefined; },
            };
        },
        async delete(name) { return stores.delete(name); },
        _stores: stores,
    };
}

async function stash(caches, meta, blobs = []) {
    const cache = await caches.open('share-target-stash');
    await cache.put('/share-meta', new Response(JSON.stringify(meta)));
    for (let i = 0; i < blobs.length; i++) {
        await cache.put(`/share-file-${i}`, new Response(blobs[i]));
    }
}

describe('takeShare', () => {
    beforeEach(() => {
        window.caches = makeCaches();
    });

    it('returns null when nothing was shared', async () => {
        expect(await takeShare()).toBeNull();
    });

    it('returns null when Cache API is unavailable', async () => {
        delete window.caches;
        expect(await takeShare()).toBeNull();
    });

    it('round-trips text and files, then clears the stash', async () => {
        await stash(
            window.caches,
            { text: 'hello', title: 't', url: 'https://x', files: [{ name: 'a.txt', type: 'text/plain' }] },
            [new Blob(['file-content'], { type: 'text/plain' })]
        );

        const shared = await takeShare();
        expect(shared.text).toBe('hello');
        expect(shared.files).toHaveLength(1);
        expect(shared.files[0].name).toBe('a.txt');
        expect(shared.files[0].type).toBe('text/plain');
        expect(await shared.files[0].text()).toBe('file-content');

        // stash cleared → second call finds nothing
        expect(await takeShare()).toBeNull();
    });
});
