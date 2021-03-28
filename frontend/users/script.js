import TagsService from '/services/tags.js'
import SubscriptionsService from '/services/subscriptions.js'
import ItemsService from '/services/items.js'

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
    <span class="container" onclick="this.dispatchEvent(new CustomEvent('SubscriptionSelected', {
        detail: {
            id: '${subscription.id}',
        },
        bubbles: true,
    }))">
        <img class="icon" src="${!!subscription.icon_url ? subscription.icon_url : '/img/rss.svg'}"></img>
        <span class="title">${subscription.title}</span>
    </span>
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
        <span class="title" onclick="this.dispatchEvent(new CustomEvent('TagSelected', {
        detail: {
            id: '${tag.id}',
        },
        bubbles: true,
    }))">${tag.title}</span>
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
    const span = document.createElement('span')
    span.id = id
    span.classList.add('toast-message')
    span.innerText = message

    if (promise) promise.then((v) => {
        if (onSuccess) document.getElementById(id).innerText = onSuccess(v)
    }).catch(e => {
        document.getElementById(id).innerText = e
    }).finally(() => {
        setTimeout(() => document.getElementById(id).parentNode.remove(), 3000)
    })

    return span.outerHTML
}

const renderToast = (promise, message, onSuccess) => `
    <span class="toast-container show">
        ${renderToastMessage(promise, message, onSuccess)}
        <button class="toast-button" type="button" onclick="const parent = this.parentNode; parentNode.remove();">
            <svg viewBox="0 0 24 24" width="24" height="24" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
        </button>
    </span>
`

const addToastMessage = async (promise, message, onSuccess) => {
    const html = renderToast(promise, message, onSuccess)
    document.querySelector('#toasts-container').insertAdjacentHTML('afterbegin', html)
}

const renderItem = (item, subscription) => `
    <div class="container item-container" created="${item.created}" onclick="this.dispatchEvent(new CustomEvent('ItemSelected', {
        detail: {
            id: '${item.id}',
        },
        bubbles: true,
    }))">
        <span class="item-title">${item.title}</span>
        <span class="container-footer">
            ${subscription && !!subscription.icon_url ? '<img class="small-icon" src="' + subscription.icon_url + '"></img>' : ''}
            ${subscription && !!!subscription.icon_url ? '<img class="small-icon" src="/img/rss.svg"></img>' : ''}
            ${subscription ? '<span style="font-size:smaller;white-space: nowrap;overflow: hidden;text-overflow: ellipsis;">'
                    + subscription.title
                    + '</span>' : ''}
            <span title="${new Date(item.created).toLocaleString()}" class="item-date">
                ${Intl.DateTimeFormat(undefined, {
                    year: 'numeric',
                    month: 'numeric',
                    day: 'numeric',
                }).format(new Date(item.created))}
            </span>
        </span>
    </div>
`

const renderReader = (item) => `
    <h2><a href="${item.url}" target="_blank">${item.title}</a></h2>
    <div>${item.summary}</div>
`

document.querySelector('#tags-menu').addEventListener('SubscriptionSelected', async (e) => {
    const subscriptionId = e.detail.id

    deleteState('tag')
    storeState('subscription', subscriptionId)

    document.querySelector('#items-list').innerHTML = await listItemsBySubscription(subscriptionId).then((items) => {
        return items.map(item => renderItem(item, subscriptionById.get(item.subscription_id))).join('')
    })
})

document.querySelector('#tags-menu').addEventListener('TagSelected', async (e) => {
    const tagId = e.detail.id

    deleteState('subscription')
    storeState('tag', tagId)

    document.querySelector('#items-list').innerHTML = await listItemsByTag(tagId).then((items) => {
        return items.map(item => renderItem(item, subscriptionById.get(item.subscription_id))).join('')
    })
})

document.querySelector('#items-list').addEventListener('ItemSelected', (e) => {
    const itemId = e.detail.id

    storeState('item', itemId)
    ItemsService.get(itemId).then((item) => {
        document.querySelector('#reader').innerHTML = renderReader(item)
    })
})

document.querySelector('#tags-menu').addEventListener('SubscriptionCreate', (e) => {
    const params = e.detail.params
    const promise = e.detail.promise

    promise.then((subscription) => {
        const html = renderSubscription(subscription)
        if (subscription.tag_ids.length == 0) {
            document.getElementById('no-tags-list').insertAdjacentHTML('afterbegin', html)
        } else {
            subscription.tag_ids.forEach((tagId) => {
                document.getElementById(tagId).insertAdjacentHTML('afterbegin', html)
            })
        }
    })

    addToastMessage(promise,
        `Subscribing: ${params.url}`,
        (subscription) => `Subscribed: ${subscription.title}`,
    )
})

document.querySelector('#tags-menu').addEventListener('TagCreate', (e) => {
    const params = e.detail.params
    const promise = e.detail.promise

    promise.then((tag) => {
        const html = renderTag(tag, [])
        document.getElementById('tags-list').insertAdjacentHTML('afterbegin', html)
    })

    addToastMessage(promise,
        `Creating tag: ${params.title}`,
        (tag) => `New tag: ${tag.title}`,
    )
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
        document.querySelector('#items-list').insertAdjacentHTML('beforeend', items.map(renderItem).join(''))
    })
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
    document.querySelector('#reader').innerHTML = renderReader(item)
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

    if (itemsByTag.length > 0) document.querySelector('#items-list').innerHTML = itemsByTag.map((item) => {
        return renderItem(item, subscriptionById.get(item.subscription_id))
    }).join('')

    if (itemsBySubscription.length > 0) document.querySelector('#items-list').innerHTML = itemsBySubscription.map((item) => {
        return renderItem(item, subscriptionById.get(item.subscription_id))
    }).join('')

    document.querySelector("#tags-list").innerHTML = renderTags(tags, subscriptions)
    document.querySelector("#no-tags-list").innerHTML = subscriptions.filter(s => s.tag_ids.length === 0)
        .map(renderSubscription).join('')

    const addButton = document.querySelector('#add-button')
    tags.forEach(tag => addButton.addTag(tag))
})
