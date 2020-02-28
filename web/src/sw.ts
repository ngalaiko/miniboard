const version: string = '__VERSION__'

const filesToCache: string[] = [
    '/index.html',
    '/static/site.webmanifest',
    '/static/favicon-32x32.png',
]

const cacheName: string = `${version}-miniboard`

const initialize = (service: ServiceWorkerGlobalScope): void => {
    service.addEventListener('install', (event: ExtendableEvent) => {
        event.waitUntil(
            caches.open(cacheName).then(function addToCache(cache) {
                return cache.addAll(filesToCache)
            })
        )
    })

    service.addEventListener('activate', (event: ExtendableEvent) => {
        event.waitUntil(
            caches
                .keys()
                .then((keys) =>
                    Promise.all(keys.filter((key) => !key.startsWith(version)).map((key) => caches.delete(key)))
                )
        )
    })

    service.addEventListener('fetch', (event: FetchEvent) => {
        if (event.request.method !== 'GET') return

        const url = new URL(event.request.url)

        if (url.pathname.startsWith('/api/')) return

        event.respondWith(
            caches.match(event.request).then((cached) => {
                const fetchedFromNetwork = (response: Response) => {
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
}

initialize(self as any)
