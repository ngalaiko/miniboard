const apiUrl = window.location.hostname !== 'localhost'
    ? 'https://api.miniboard.app'
    : 'http://localhost:80';

class API {
    async get(url) {
        return await this.fetch(url)
    }

    async post(url, request) {
        return await this.fetch(url, {
            method: 'POST',
            body: JSON.stringify(request),
            headers: new Headers({
                "Content-Type": "application/json",
            }),
        })
    }

    async fetch(url, params) {
        if (params == undefined) params = {}
        params.credentials = 'include'

        const response = await fetch(apiUrl + url, params)
        const body = await response.json()
        if (response.status !== 200) {
            throw new Error(body.message)
        }

        return body
    }
}

const Api = new API()

class Subscriptions {
    async list(params) {
        if (params === undefined) params = {}

        const pageSizeQuery = params.pageSize !== undefined
            ? `&page_size=${params.pageSize}`
            : ''

        const createdLtQuery = params.createdLt !== undefined
            ? `&created_lt=${encodeURIComponent(params.createdLt)}`
            : ''

        const url = '/v1/subscriptions/?' + pageSizeQuery + createdLtQuery

        const body = await Api.get(url)

        return body.subscriptions
    }

    async create(params) {
        if (params === undefined) params = {}

        const request = {
            url: params.url,
        }

        if (params.tagIds !== undefined) {
            request.tag_ids = params.tagIds
        }

        return await Api.post('/v1/subscriptions/', request)
    }
}

class Items {
    async get(id) {
        return await Api.get(`/v1/items/${id}/`)
    }

    async list(params) {
        if (params === undefined) params = {}

        const pageSizeQuery = params.pageSize !== undefined
            ? `&page_size=${params.pageSize}`
            : ''

        const createdLtQuery = params.createdLt !== undefined
            ? `&created_lt=${encodeURIComponent(params.createdLt)}`
            : ''

        let url = ''
        switch (true) {
        case !!params.subscriptionId && !!params.tagId:
            throw new Error('not implemented')
        case !!params.subscriptionId:
            url = `/v1/subscriptions/${params.subscriptionId}/items/?` + pageSizeQuery + createdLtQuery
            break
        case !!params.tagId:
            url = `/v1/tags/${params.tagId}/items/?` + pageSizeQuery + createdLtQuery
            break
        default:
            url = `/v1/items/?` + pageSizeQuery + createdLtQuery
            break
        }

        const body = await Api.get(url)
        return body.items
    }
}

class Operations {
    async get(id) {
        return await Api.get(`/v1/operations/${id}/`)
    }

    async wait(id) {
        const operation = await this.get(id)

        switch (true) {
        case !operation.done:
            await new Promise(r => setTimeout(r, 1000))
            return await this.wait(id)
        case operation.result.error !== undefined:
            throw new Error(operation.result.error.message)
        case operation.result.response !== undefined:
            return operation.result.response
        default:
            throw new Error(`invalid state for operation ${id}`)
        }
    }
}

class Imports {
    async create(raw) {
        return await Api.fetch('/v1/imports/', {
            method: 'POST',
            body: raw,
            headers: new Headers({
                "Content-Type": "application/xml",
            }),
        })
    }
}

class Tags {
    async list(params) {
        if (params === undefined) params = {}

        const pageSizeQuery = params.pageSize !== undefined
            ? `&page_size=${params.pageSize}`
            : ''

        const createdLtQuery = params.createdLt !== undefined
            ? `&created_lt=${encodeURIComponent(params.createdLt)}`
            : ''

        const url = '/v1/tags/?' + pageSizeQuery + createdLtQuery

        const body = await Api.get(url)

        return body.tags
    }

    async create(params) {
        if (params === undefined) params = {}

        const request = {
            title: params.title,
        }

        return await Api.post('/v1/tags/', request)
    }
}

const TagsService = new Tags()
const OperationsService = new Operations()
const ItemsService = new Items()
const SubscriptionsService = new Subscriptions()
const ImportsService = new Imports()

const storeState = (key, value) => {
    const urlParams = new URLSearchParams(window.location.search.slice(1))

    if (urlParams.get(key) === value) return

    urlParams.set(key, value)

    let refresh = window.location.protocol +
        "//" + window.location.host + window.location.pathname +
        `?${urlParams.toString()}`
    window.history.pushState({ path: refresh }, '', refresh)
}

const getState = (key) => {
    const urlParams = new URLSearchParams(window.location.search.slice(1))
    const value = urlParams.get(key)
    return value ? value : undefined
}

const deleteState = (key) => {
    const urlParams = new URLSearchParams(window.location.search.slice(1))
    urlParams.delete(key)

    let refresh = window.location.protocol +
        "//" + window.location.host + window.location.pathname +
        `?${urlParams.toString()}`
    window.history.pushState({ path: refresh }, '', refresh)
}

const listAllTags = async (pageSize, createdLt) => {
    const params = {}

    if (pageSize === undefined) pageSize = 100
    if (pageSize !== undefined) params.pageSize = pageSize
    if (createdLt !== undefined) params.createdLt = createdLt

    const tags = await TagsService.list(params)

    if (tags.length < pageSize) {
        return tags
    }

    params.createdLt = tags[tags.length - 1].created

    return tags.concat(await TagsService.list(params))
}

const listAllSubscriptions = async (pageSize, createdLt) => {
    const params = {}

    if (pageSize === undefined) pageSize = 100
    if (pageSize !== undefined) params.pageSize = pageSize
    if (createdLt !== undefined) params.createdLt = createdLt

    const tags = await SubscriptionsService.list(params)

    if (tags.length < pageSize) {
        return tags
    }

    params.createdLt = tags[tags.length - 1].created

    return tags.concat(await SubscriptionsService.list(params))
}

const renderSubscription = (subscription) => `
    <div class="subscription pure-g" style="text-align: left;" onclick="onSubscriptionSelected('${subscription.id}')">
        <img class="pure-u subscription-icon" width="20px" height="20px" src="${!!subscription.icon_url ? subscription.icon_url : '/img/rss.svg'}"></img>
        <span class="pure-u-3-4 subscription-title">${subscription.title}</span>
    </div>
`

const toggleTag = (tagId) => {
    const children = document.getElementById(`${tagId}-children`)
    children.hidden = !children.hidden
    document.getElementById(`${tagId}-arrow`).classList.toggle('rotate-90')
}

const renderTag = (tag, subscriptions) => `
    <div class="tag pure-g">
        <button class="pure-u pure-button button-hidden" style="text-align: left;" onclick="toggleTag('${tag.id}')" >
            <svg id="${tag.id}-arrow" class="tag-arrow" viewBox="0 0 20 20" width="20" height="20" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="9 18 15 12 9 6"></polyline>
            </svg>
        </button>
        <button class="pure-u-3-4 pure-button button-hidden tag-title" style="text-align: left;" onclick="onTagSelected('${tag.id}')">
            ${tag.title}
        </button>
    </div>
    <div id="${tag.id}-children" class="tag-subscriptions" hidden>
        ${subscriptions.map(renderSubscription).join('')}
    </div>
`

const renderTags = (tags, subscriptions) => {
    const subscriptionsByTagId = subscriptions.reduce((map, subscription) => {
        subscription.tag_ids.forEach((tagId) => {
            const list = map.get(tagId)
            if (!list) {
                map.set(tagId, [subscription])
            } else {
                list.push(subscription)
            }
        })
        return map
    }, new Map())

    return tags.map((tag) => renderTag(tag, subscriptionsByTagId.get(tag.id) || [])).join('')
}

const renderToastMessage = (promise, message, onSuccess) => {
    const id = `toast-${performance.now()}`
    const div = document.createElement('div')
    div.id = id
    div.classList.add('toast-message')
    div.innerText = message

    if (promise) promise.then((v) => {
        if (onSuccess) document.getElementById(id).innerText = onSuccess(v)
    }).catch(e => {
        document.getElementById(id).innerText = e
    }).finally(() => {
        setTimeout(() => document.getElementById(id).parentNode.remove(), 3000)
    })

    return div.outerHTML
}

const renderToast = (promise, message, onSuccess) => `
    <div class="toast-container">
        ${renderToastMessage(promise, message, onSuccess)}
        <button class="toast-button" type="button" onclick="const parent = this.parentNode; parentNode.remove();">
            <svg viewBox="0 0 24 24" width="20" height="20" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
        </button>
    </div>
`

const addToastMessage = async (promise, message, onSuccess) => {
    const html = renderToast(promise, message, onSuccess)
    document.querySelector('#toasts-container').insertAdjacentHTML('afterbegin', html)
}

const renderItemSubscriptionIcon = (subscription) => !!subscription && !!subscription.icon_url
    ? `<img class="subscription-icon" width="20" height="20" src="${subscription.icon_url}"></img>`
    : `<img class="subscription-icon" width="20" height="20" src="/img/rss.svg"></img>`

const renderItem = (item, subscription) => `
    <div class="item pure-g" created="${item.created}" onclick="onItemSelected('${item.id}')">
        <div class="pure-u">${renderItemSubscriptionIcon(subscription)}</div>
        <div class="pure-u-3-4">
            <h5 class="item-subscription-title">${subscription.title}</h5>
            <h4 class="item-title">${item.title}</h4>
        </div>
    </div>
`

const renderReader = (item) => `
    <h2><a href="${item.url}" target="_blank">${item.title}</a></h2>
    <div>${item.summary}</div>
`

const displayReader = (target, item) => {
    target.innerHTML = renderReader(item)
}

const displayItems = (target, items) => {
    const html = items.map((item) => {
        return renderItem(item, subscriptionById.get(item.subscription_id))
    }).join('')

    target.insertAdjacentHTML('beforeend', html)
}

const onSubscriptionSelected = (subscriptionId) => {
    deleteState('tag')
    storeState('subscription', subscriptionId)

    listItemsBySubscription(subscriptionId).then((items) => {
        const list = document.querySelector('#items-list')
        list.innerHTML = ''
        displayItems(list, items)
    })
}

const onItemSelected = (itemId) => {
    storeState('item', itemId)
    ItemsService.get(itemId).then((item) => {
        displayReader(document.querySelector('#reader'), item)
    })
}

const onTagSelected = (tagId) => {
    deleteState('subscription')
    storeState('tag', tagId)

    listItemsByTag(tagId).then((items) => {
        const list = document.querySelector('#items-list')
        list.innerHTML = ''
        displayItems(list, items)
    })
}

document.querySelector('#items-list').addEventListener('scroll', (e) => {
    const { scrollTop, scrollHeight, clientHeight } = e.target
    const needMore = scrollTop + clientHeight >= scrollHeight - 50
    if (!needMore) return

    const subscriptionId = getState('subscription')
    const tagId = getState('tag')
    const createdLt = document.querySelector('#items-list').lastElementChild.getAttribute('created')

    if (!createdLt) return

    document.querySelector('#items-list').insertAdjacentHTML('beforeend', '<div class="page-separator"></div')

    ItemsService.list({
        pageSize: 50,
        subscriptionId: subscriptionId,
        tagId: tagId,
        createdLt: createdLt,
    }).then((items) => {
        displayItems(document.querySelector('#items-list'), items)
    })
})

const displaySubscription = (subscription) => {
    const html = renderSubscription(subscription)
    subscriptionById.set(subscription.id, subscription)
    if (subscription.tag_ids.length == 0) {
        document.getElementById('no-tags-list').insertAdjacentHTML('afterbegin', html)
    } else {
        subscription.tag_ids.forEach((tagId) => {
            document.getElementById(tagId).insertAdjacentHTML('afterbegin', html)
        })
    }
}

window.addEventListener('keydown', (e) => {
    if (isModalClosed()) return

    switch (e.key) {
    case 'Enter':
        closeModal()

        const url = document.querySelector("#input-url").value
        if (url === '') return

        document.querySelector("#input-url").value = ''

        const promise = SubscriptionsService.create({
            url: url,
        }).then((operation) => OperationsService.wait(operation.id))

        promise.then(displaySubscription)

        addToastMessage(promise,
            `Subscribing: ${url}`,
            (subscription) => `Subscribed: ${subscription.title}`,
        )
        break
    case 'Escape':
        closeModal()
        break
    }
})

document.querySelector("#input-file").addEventListener('change', async (e) => {
    const files = e.target.files

    closeModal()

    if (files.length === 0) return

    const promise = new Promise((resolve, reject) => {
        const reader = new FileReader()
        reader.readAsBinaryString(files[0])
        reader.onload = e => {
            resolve(e.target.result)
        }
        reader.onerror = () => reject(new Error('failed to read file'))
    })
    .then((raw) => ImportsService.create(raw))
    .then((operation) => OperationsService.wait(operation.id))

    promise.then((imported) => {
        if (imported.tags) imported.tags.forEach((tag) => {
            const html = renderTag(tag, [])
            document.getElementById('tags-list').insertAdjacentHTML('afterbegin', html)
        })

        if (imported.subscriptions) imported.subscriptions.forEach(displaySubscription)
    })

    addToastMessage(promise, `Importing file...`, () => `Imported`)
})

const isModalClosed = () => {
    return document.querySelector('#modal').hidden
}

const closeModal = () => {
    document.querySelector('#modal').hidden = true
}

const showModal = () => {
    document.querySelector('#modal').hidden = false
}

document.querySelector("#background").addEventListener('click', () => {
    if (!isModalClosed()) closeModal()
})

document.querySelector('#add-button').addEventListener('click', () => {
    showModal()
})

const subscriptionById = new Map()

const listItemsBySubscription = async (subscriptionId) => {
    if (!subscriptionId) return []

    return await ItemsService.list({
        pageSize: 50,
        subscriptionId: subscriptionId,
    })
}

const listItemsByTag = async (tagId) => {
    if (!tagId) return []

    return await ItemsService.list({
        pageSize: 50,
        tagId: tagId,
    })
}

const itemId = getState('item')
if (itemId) ItemsService.get(itemId).then((item) => {
    displayReader(document.querySelector('#reader'), item)
})

Promise.all([
    listAllSubscriptions(),
    listAllTags(),
    listItemsByTag(getState('tag')),
    listItemsBySubscription(getState('subscription')),
]).then(async (values) => {
    const subscriptions = values[0]
    const tags = values[1]
    const itemsByTag = values[2]
    const itemsBySubscription = values[3]

    subscriptions.forEach(s => subscriptionById.set(s.id, s))

    const list = document.querySelector('#items-list')
    if (itemsBySubscription.length > 0) displayItems(list, itemsBySubscription)
    if (itemsByTag.length > 0) displayItems(list, itemsByTag)

    document.querySelector("#tags-list").innerHTML = renderTags(tags, subscriptions)
    document.querySelector("#no-tags-list").innerHTML = subscriptions.filter(s => s.tag_ids.length === 0)
        .map(renderSubscription).join('')
})
