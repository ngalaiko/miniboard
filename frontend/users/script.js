const apiUrl = window.location.hostname !== 'localhost'
    ? 'https://api.miniboard.app'
    : 'http://localhost:80';

const loadFeeds = async (pageSize, createdLt) => {
    const pageSizeQuery = pageSize !== undefined
        ? `&page_size=${pageSize}`
        : ''

    const createdLtQuery = createdLt !== undefined
        ? `&created_lt=${encodeURIComponent(createdLt)}`
        : ''

    const url = apiUrl + '/v1/feeds?' + pageSizeQuery + createdLtQuery
    const response = await fetch(url, {
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

    if (body.feeds.length < pageSize) return

    loadFeeds(pageSize, body.feeds.pop().created)
}

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

            const tagIdsByUrl = {}
            Array.from(dom.getElementsByTagName('outline'))
                .filter(item => item.getAttribute('xmlUrl') === null)
                .map(item => {
                    const title = item.getAttribute('title')
                    tag = tagsByTitle[title] !== undefined
                        ? tagsByTitle[title]
                        : createTag(title)

                    return {
                        tagId: tag.id,
                        urls: Array.from(item.getElementsByTagName('outline'))
                            .filter(item => item.getAttribute('xmlUrl') !== null)
                            .map(item => item.getAttribute('xmlUrl')),
                    }
                })
                .forEach(item => item.urls.forEach(
                    url => {
                        if (!tagIdsByUrl[url]) tagIdsByUrl[url] = []
                        tagIdsByUrl[url].push(item.tagId)
                    })
                )

            for (const url of Object.keys(tagIdsByUrl)) {
                try {
                    await addFeed(url, tagIdsByUrl[url])
                } catch {}
            }
        })
    })

    addModal.addEventListener('UrlAdded', (e) => addFeed(e.detail.url))
})

const tagsByTitle = {}

const loadTags = async (pageSize, createdLt) => {
    const pageSizeQuery = pageSize !== undefined
        ? `&page_size=${pageSize}`
        : ''

    const createdLtQuery = createdLt !== undefined
        ? `&created_lt=${encodeURIComponent(createdLt)}`
        : ''

    const url = apiUrl + '/v1/tags?' + pageSizeQuery + createdLtQuery
    const response = await fetch(url, {
        credentials: 'include',
    })

    const body = await response.json()

    if (response.status !== 200) {
        throw new Exception(`failed to fetch feeds: ${body.message}`)
        return
    }

    body.tags.forEach(tag => tagsByTitle[tag.title] = tag)

    if (body.tags.length < pageSize) return

    loadFeeds(pageSize, body.feeds.pop().created)
}

const createTag = async (title) => {
    const response = await fetch(apiUrl + '/v1/tags', {
        credentials: 'include',
        method: 'POST',
        body: JSON.stringify({
            title: title,
        }),
    })

    const body = await response.json()
    if (response.status !== 200) {
        throw new Exception(`failed to create feed: ${body.message}`)
        return
    }

    return body
}

const addFeed = async (url, tagIds) => {
    const request = {
        url: url,
    }

    if (tagIds !== undefined) {
        request.tag_ids = tagIds
    }

    const response = await fetch(apiUrl + '/v1/feeds', {
        credentials: 'include',
        method: 'POST',
        body: JSON.stringify(request),
    })

    const body = await response.json()
    if (response.status !== 200) {
        throw new Exception(`failed to create feed: ${body.message}`)
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

loadFeeds(100)
loadTags(100)
