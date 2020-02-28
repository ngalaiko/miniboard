const version = '__VERSION__'

const filesToCache = [
    '/index.html',
    '/static/site.webmanifest',
    '/static/favicon-32x32.png',
]

const cacheName = `${version}-miniboard`

self.addEventListener('install', (event) => {
    event.waitUntil(
        caches.open(cacheName).then(function addToCache(cache) {
            return cache.addAll(filesToCache)
        })
    )
})

self.addEventListener('activate', (event) => {
    event.waitUntil(
        caches
            .keys()
            .then((keys) =>
                Promise.all(keys.filter((key) => !key.startsWith(version)).map((key) => caches.delete(key)))
            )
    )
})

self.addEventListener('fetch', (event) => {
    if (event.request.method !== 'GET') return

    const url = new URL(event.request.url)

    if (url.pathname.startsWith('/api/')) return

    event.respondWith(
        caches.match(event.request).then((cached) => {
            const fetchedFromNetwork = (response) => {
                const cacheResponseCopy = response.clone()

                caches.open(cacheName).then((cache) => {
                    cache.put(event.request, cacheResponseCopy)
                })
                return response
            }

            const unableToResolve = () => {
                return new Response('Service Unavailable', {
                    status: 503,
                })
            }

            const networked = fetch(event.request)
                .then(fetchedFromNetwork, unableToResolve)
                .catch(unableToResolve)

            return cached || networked
        })
    )
})