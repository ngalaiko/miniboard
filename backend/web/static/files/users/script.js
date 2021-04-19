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

        const response = await fetch(url, params)
        const body = await response.json()
        if (response.status !== 200) {
            throw new Error(body.message)
        }

        return body
    }
}

const Api = new API()

class Subscriptions {
    async create(params) {
        if (params === undefined) params = {}

        const request = {
            url: params.url,
        }

        if (params.tagIds !== undefined) {
            request.tag_ids = params.tagIds
        }

        return await Api.post('/api/v1/subscriptions/', request)
    }
}

class Operations {
    async get(id) {
        return await Api.get(`/api/v1/operations/${id}/`)
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
        return await Api.fetch('/api/v1/imports/', {
            method: 'POST',
            body: raw,
            headers: new Headers({
                "Content-Type": "application/xml",
            }),
        })
    }
}

const OperationsService = new Operations()
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

const onSubscriptionSelected = (subscriptionId) => {
    deleteState('tag')
    storeState('subscription', subscriptionId)

    sendMessage("items:load", {
        subscriptionId: subscriptionId,
    })
}

const onItemSelected = (itemId) => {
    sendMessage("item:selected", {id: itemId})
    storeState('item', itemId)
}

const onTagSelected = (tagId) => {
    deleteState('subscription')
    storeState('tag', tagId)

    sendMessage("items:load", {
        tagId: tagId,
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

    sendMessage("items:loadmore", {
        tagId: tagId,
        subscriptionId: subscriptionId,
        createdLt: createdLt,
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
