// Service worker for PWA installability and Android share-target intake.
// Deliberately network-only: this service deals in one-time links,
// caching an API response could burn a download.

const STASH = 'share-target-stash';

self.addEventListener('install', () => self.skipWaiting());
self.addEventListener('activate', (e) => e.waitUntil(self.clients.claim()));

self.addEventListener('fetch', (e) => {
    const url = new URL(e.request.url);
    if (e.request.method === 'POST' && url.pathname === '/share-target') {
        e.respondWith(handleShare(e.request));
    }
    // Everything else falls through to the network untouched.
});

async function handleShare(request) {
    const form = await request.formData();
    const files = form.getAll('files').filter((f) => f instanceof File);
    const cache = await caches.open(STASH);
    const meta = {
        text: form.get('text') || '',
        title: form.get('title') || '',
        url: form.get('url') || '',
        files: files.map((f) => ({ name: f.name, type: f.type })),
    };
    await cache.put('/share-meta', new Response(JSON.stringify(meta)));
    await Promise.all(
        files.map((f, i) => cache.put(`/share-file-${i}`, new Response(f)))
    );
    return Response.redirect('/#/?share=1', 303);
}
