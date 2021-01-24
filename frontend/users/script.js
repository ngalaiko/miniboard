const apiUrl = window.location.hostname == 'localhost'
    ? 'http://localhost:80'
    : 'https://api.miniboard.app'

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
        if (item.icon_url !== undefined) {
            item.setAttribute('icon', feed.icon_url)
        }
    })
})()

document.querySelector("#add-button").addEventListener('click', (e) => {
    import('./components/modal.js')

    const addModal = document.createElement('add-modal')
    document.body.appendChild(addModal)

    addModal.addEventListener('Closed', (e) => {
        document.body.removeChild(addModal)
    })

    addModal.addEventListener('FeedAdded', (e) => addFeed(e.detail.url))
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
