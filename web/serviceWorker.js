const filesToCache = [
    '.',
    'index.html', 
    'site.webmanifest', 
    './app.es5.min.js'
];

const cacheName = 'miniboard'

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(cacheName).then(function addToCache(cache) {
      return cache.addAll(filesToCache);
    })
  );
});

self.addEventListener('fetch', (event) => {
  if (event.request.method !== 'GET') return;

  event.respondWith(
    caches.match(event.request).then((cached) => {
      function fetchedFromNetwork(response) {
        const cacheResponseCopy = response.clone();

        caches.open(cacheName).then((cache) => {
          cache.put(event.request, cacheResponseCopy);
        });
        return response;
      }

      function unableToResolve() {
        return new Response('Service Unavailable', {
          status: 503,
        });
      }

      const networked = fetch(event.request)
        .then(fetchedFromNetwork, unableToResolve)
        .catch(unableToResolve);

      return cached || networked;
    })
  );
});

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches
      .keys()
      .then((keys) =>
        Promise.all(keys.filter((key) => !key.startsWith(version)).map((key) => caches.delete(key)))
      )
  );
});
