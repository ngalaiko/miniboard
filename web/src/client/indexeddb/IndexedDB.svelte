<script context='module'>
    const disabled = () => {
        let $ = {}
        $.add = () => Promise.resolve(undefined)
        $.update = () => Promise.resolve(undefined)
        $.get = () => Promise.resolve(undefined)
        $.delete = () => Promise.resolve(undefined)
        $.forEach = () => Promise.resolve(undefined)
        return $
    }

    export const IndexedDB = () => {
        let $ = {}

        const indexedDB = window.indexedDB || window.mozIndexedDB
            || window.webkitIndexedDB || window.msIndexedDB

        let db

        if (!indexedDB) {
            return disabled()
        }

        const request = indexedDB.open('miniboard', 3)
        request.onerror = (event) => {
            console.error(event.target.errorCode)
        }
        request.onsuccess = (event) => db = event.target.result

        request.onupgradeneeded = (event) => {
            let db = event.target.result

            let collectionNames = [
                'articles',
            ]

            collectionNames.forEach(collectionName => {
                db.createObjectStore(collectionName, {
                    keyPath: 'name'
                }).createIndex('name', 'name', { unique: true });
            })
        }

        const collectionName = (name) => {
            let parts = name.split('/')
            return parts[parts.length - 2]
        }

        $.add = (value) => new Promise((resolve, reject) => {
            let name = collectionName(value.name)
            let tx = db.transaction(name, 'readwrite')
            tx.onerror = (event) => reject(event.target.errorCode)

            let request = tx.objectStore(name).add(value)
            request.onerror = (event) => reject(event.target.errorCode)
            request.onsuccess = (event) => resolve(event.target.result)
        })

        $.update = (value) => new Promise((resolve, reject) => {
            let name = collectionName(value.name)
            let tx = db.transaction(name, 'readwrite')
            tx.onerror = (event) => reject(event.target.errorCode)

            let request = tx.objectStore(name).put(value)
            request.onsuccess = (event) => resolve(event.target.result)
            request.onerror = (event) => reject(event.target.errorCode)
        })

        $.get = (key) => new Promise((resolve, reject) => {
            let name = collectionName(key)
            let tx = db.transaction(name, 'readonly')
            tx.onerror = (event) => { reject(event.target.errorCode) }

            let request = tx.objectStore(name).get(key)
            request.onerror = (event) => reject(event.target.errorCode)
            request.onsuccess = (event) => resolve(event.target.result)
        })

        $.delete = (key) => new Promise((resolve, reject) => {
            let name = collectionName(key)
            let tx = db.transaction(name, 'readwrite')
            tx.onerror = (event) => reject(event.target.errorCode)

            let request = tx.objectStore(name).delete(key)
            request.onerror = (event) => reject(event.target.errorCode)
            request.onsuccess = (event) => resolve(event.target.result)
        })

        $.forEach = (name, eachFunc) => new Promise((resolve, reject) => {
            let tx = db.transaction(name, 'readonly')
            tx.onerror = (event) => reject(event.target.errorCode)

            let request = tx.objectStore(name).index('name').openCursor()
            request.onerror = (event) => reject(event.target.errorCode)
            request.onsuccess = (event) => {
                let cursor = event.target.result

                if (cursor) {
                    eachFunc(cursor.value)
                    cursor.continue()
                }

                resolve(undefined)
            }
        })

        return $
    }
</script>
