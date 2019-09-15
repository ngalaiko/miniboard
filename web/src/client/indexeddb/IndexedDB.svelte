<script context='module'>
    const disabled = () => {
        let $ = {}
        $.add = () => Promise.reject('not supported')
        $.get = () => Promise.reject('not supported')
        $.delete = () => Promise.reject('not supported')
        $.forEach = () => Promise.reject('not supported')
        return $
    }

    let database = () => new Promise((resolve, reject) => {
        const indexedDB = window.indexedDB || window.mozIndexedDB
            || window.webkitIndexedDB || window.msIndexedDB

        if (!indexedDB) {
            reject('indexed db not supported')
        }

        const request = indexedDB.open('miniboard', 3)
        request.onerror = (event) => reject(`can't open IndexedDB: ${event.target.errorCode}`)
        request.onsuccess = (event) => resolve(event.target.result)
        request.onupgradeneeded = (event) => {
            let db = event.target.result

            let collectionNames = [
                'articles',
            ]

            collectionNames.forEach(collectionName => {
                db.createObjectStore(collectionName, {
                    keyPath: 'name'
                }).createIndex('name', 'name', { unique: true })
            })
        }
    })

    export const IndexedDB = async () => {
        let $ = {}

        let db
        try {
            db = await database()
        } catch (e) {
            console.log(`can't use IndexedDB because: ${e}`)
            return disabled()
        }

        const collectionName = (name) => {
            let parts = name.split('/')
            return parts[parts.length - 2]
        }

        $.add = (value) => new Promise((resolve, reject) => {
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

                if (cursor && eachFunc(cursor.value)) {
                    cursor.continue()
                } else {
                    resolve(undefined)
                }
            }
        })

        return $
    }
</script>
