import TagsService from '/services/tags.js'
import SubscriptionsService from '/services/subscriptions.js'
import OperationsService from '/services/operations.js'
import ItemsService from '/services/items.js'
import ImportsService from '/services/imports.js'

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
    <div class="container" onclick="this.dispatchEvent(new CustomEvent('SubscriptionSelected', {
        detail: {
            id: '${subscription.id}',
        },
        bubbles: true,
    }))">
        <img class="icon" src="${!!subscription.icon_url ? subscription.icon_url : '/img/rss.svg'}"></img>
        <div class="title">${subscription.title}</div>
    </div>
`

const renderTag = (tag, subscriptions) => `
    <div style="display:flex;align-items:center;cursor:pointer;">
        <button type="button" style="background:none;border:none;padding:0;" onclick="
            const el = document.getElementById('${tag.id}');
            el.hidden = !el.hidden;
            document.getElementById('${tag.id}-arrow').classList.toggle('rotate');
            ">
            <svg id="${tag.id}-arrow" viewBox="0 0 24 24" width="20" height="20" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="9 18 15 12 9 6"></polyline>
            </svg>
        </button>
        <div class="title" onclick="this.dispatchEvent(new CustomEvent('TagSelected', {
        detail: {
            id: '${tag.id}',
        },
        bubbles: true,
    }))">${tag.title}</div>
    </div>
    <div id="${tag.id}" hidden>
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
    <div class="toast-container show">
        ${renderToastMessage(promise, message, onSuccess)}
        <button class="toast-button" type="button" onclick="const parent = this.parentNode; parentNode.remove();">
            <svg viewBox="0 0 24 24" width="24" height="24" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round">
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

const renderItemCreated = (created) => {
    if (!created) return  '<div class="item-date">N/A</div>'
    const date = new Date(created)
    const formatter = Intl.DateTimeFormat(undefined, {
        year: 'numeric',
        month: 'numeric',
        day: 'numeric',
    })
    return `<div title="${date.toLocaleString()}" class="item-date">${formatter.format(date)}</div>`
}

const renderItemSubscriptionTitle = (subscription) => `
    <div style="font-size:smaller;white-space: nowrap;overflow: hidden;text-overflow: ellipsis;">
        ${subscription.title}
    </div>
`

const renderItemSubscriptionIcon = (subscription) => !!subscription.icon_url
    ? `<img class="small-icon" src="${subscription.icon_url}"></img>`
    : `<img class="small-icon" src="/img/rss.svg"></img>`

const renderItemSubscription = (subscription) => `
    ${renderItemSubscriptionIcon(subscription)}
    ${renderItemSubscriptionTitle(subscription)}
`

const renderItem = (item, subscription) => `
    <div id="${item.id}" class="container item-container" created="${item.created}">
        <div class="item-title">${item.title}</div>
        <div class="container-footer">
            ${renderItemSubscription(subscription)}
            ${renderItemCreated(item.created)}
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
    items.map((item) => {
        const template = document.createElement('template')
        template.innerHTML = renderItem(item, subscriptionById.get(item.subscription_id))
        target.appendChild(template.content)
        document.getElementById(item.id).addEventListener('click', () => {
            storeState('item', item.id)
            displayReader(document.getElementById('reader'), item)
        })
    })
}

document.querySelector('#tags-menu').addEventListener('SubscriptionSelected', async (e) => {
    const subscriptionId = e.detail.id

    deleteState('tag')
    storeState('subscription', subscriptionId)

    listItemsBySubscription(subscriptionId).then((items) => {
        const list = document.querySelector('#items-list')
        list.innerHTML = ''
        displayItems(list, items)
    })
})

document.querySelector('#tags-menu').addEventListener('TagSelected', async (e) => {
    const tagId = e.detail.id

    deleteState('subscription')
    storeState('tag', tagId)

    listItemsByTag(tagId).then((items) => {
        const list = document.querySelector('#items-list')
        list.innerHTML = ''
        displayItems(list, items)
    })
})

document.querySelector('#items-list').addEventListener('scroll', (e) => {
    const { scrollTop, scrollHeight, clientHeight } = e.target
    const needMore = scrollTop + clientHeight >= scrollHeight - 50
    if (!needMore) return

    const subscriptionId = getState('subscription')
    const tagId = getState('tag')
    const createdLt = document.querySelector('#items-list').lastElementChild.getAttribute('created')

    if (!createdLt) return

    document.querySelector('#items-list').insertAdjacentHTML('beforeend', '<div class="page-separator"></div')

    const pageSize = 100
    ItemsService.list({
        pageSize: pageSize,
        subscriptionId: subscriptionId,
        tagId: tagId,
        createdLt: createdLt,
    }).then((items) => {
        displayItems(document.querySelector('#items-list'), items)
    })
})

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

        promise.then((subscription) => {
            const html = renderSubscription(subscription)
            subscriptionById.set(subscription.id, subscription)
            if (subscription.tag_ids.length == 0) {
                document.getElementById('no-tags-list').insertAdjacentHTML('afterbegin', html)
            } else {
                subscription.tag_ids.forEach((tagId) => {
                    document.getElementById(tagId).insertAdjacentHTML('afterbegin', html)
                })
            }
        })

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

        if (imported.subscriptions) imported.subscriptions.forEach((subscription) => {
            subscriptionById.set(subscription.id, subscription)
            const html = renderSubscription(subscription)
            if (subscription.tag_ids.length == 0) {
                document.getElementById('no-tags-list').insertAdjacentHTML('afterbegin', html)
            } else {
                subscription.tag_ids.forEach((tagId) => {
                    document.getElementById(tagId).insertAdjacentHTML('afterbegin', html)
                })
            }
        })
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
        subscriptionId: subscriptionId,
    })
}

const listItemsByTag = async (tagId) => {
    if (!tagId) return []

    return await ItemsService.list({
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
    displayItems(list, itemsBySubscription)
    displayItems(list, itemsByTag)

    document.querySelector("#tags-list").innerHTML = renderTags(tags, subscriptions)
    document.querySelector("#no-tags-list").innerHTML = subscriptions.filter(s => s.tag_ids.length === 0)
        .map(renderSubscription).join('')
})
