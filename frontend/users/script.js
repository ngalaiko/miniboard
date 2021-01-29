const apiUrl = window.location.hostname !== 'localhost'
    ? 'https://api.miniboard.app'
    : 'http://localhost:80';

(async () => {
    const response = await fetch(apiUrl + '/v1/feeds', {
        credentials: 'include',
    })

    const body = await response.json()

    if (response.status !== 200) {
        console.error(`failed to fetch feeds: ${body.message}`)
        return
    }

    const list = document.querySelector('#feeds-list')
    body.feeds.map(feed => {
        const item = list.appendChild(
            document.createElement('li').appendChild(
                document.createElement('x-feed')
            )
        )

        item.setAttribute('title', feed.title)
        if (feed.icon_url !== undefined) {
            item.setAttribute('icon', feed.icon_url)
        }
    })
})()

document.querySelector("#feeds-add-button").addEventListener('click', (e) => {
    import('./components/modal.js')

    const addModal = document.createElement('add-modal')
    document.body.appendChild(addModal)

    addModal.addEventListener('Closed', (e) => {
        document.body.removeChild(addModal)
    })

    addModal.addEventListener('FileAdded', (e) => {
        new Promise((resolve, reject) => {
            const reader = new FileReader()
            reader.readAsBinaryString(e.detail.file)
            reader.onload = e => {
                resolve(e.target.result)
            }
            reader.onerror = () => reject(new Error('failed to read file'))
        })
        .then(async (raw) => {
            const parser = new DOMParser()
            const dom = parser.parseFromString(raw, "text/xml")

            if (dom.documentElement.nodeName !== 'opml') {
                throw new Error('only opml files supported')
            }

            const urls = Array.from(dom.getElementsByTagName('outline'))
                .map(item => item.getAttribute('xmlUrl'))
                .filter(item => item !== null)

            for (const url of urls) {
                try {
                    await addFeed(url)
                } catch {}
            }
        })
    })

    addModal.addEventListener('UrlAdded', (e) => addFeed(e.detail.url))
})

const addFeed = async (url) => {
    const response = await fetch(apiUrl + '/v1/feeds', {
        credentials: 'include',
        method: 'POST',
        body: JSON.stringify({
            url: url,
        }),
    })

    const body = await response.json()
    if (response.status !== 200) {
        alert(`failed to create feed: ${body.message}`)
        return
    }

    watchOperationStatus(url, body)
}

const watchOperationStatus = async (url, operation) => {
    switch (true) {
    case !operation.done:
        const response = await fetch(apiUrl + `/v1/operations/${operation.id}`, {
            credentials: 'include',
        })
        const body = await response.json()
        if (response.status !== 200) {
            console.error(`failed to fetch operation status: ${body.message}`)
            return
        }
        window.setTimeout(() => watchOperationStatus(url, body), 1000)
        break
    case !(operation.result.error === undefined):
        console.error(`${url}: ${operation.result.error.message}`)
        break
    case !(operation.result.response === undefined):
        console.debug(`${url}: added`)
        break
    }
}
