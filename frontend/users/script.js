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

    return urlParams.get(key)
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
    <span class="subscription-container">
        <img class="subscription-icon" src="${!!subscription.icon_url ? subscription.icon_url : '/img/rss.svg'}"></img>
        <span class="subscription-title">${subscription.title}</span>
    </span>
`

const renderSubscriptions = (tagId, subscriptions) => `
    <div id="${tagId}">
        ${subscriptions.map(renderSubscription).join('')}
    </div>
`

const renderTag = (tag, subscriptions) => `
    <details class="tag-container">
        <summary class="tag-title">${tag.title}</summary>
        ${renderSubscriptions(tag.id, subscriptions)}
    </details>
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

const addToastMessage = async (promise, message, onSuccess) => {
    await import('./components/toast.js')
    
    const toast = document.createElement('x-toast')
    toast.setAttribute('message', message)

    if (promise) promise.then((v) => {
        if (onSuccess) toast.setAttribute('message', onSuccess(v))
    }).catch(e => {
        toast.setAttribute('message', e)
    }).finally(() => {
        toast.close()
    })

    document.querySelector('#toasts-container').appendChild(toast)
}

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

Promise.all([listAllSubscriptions(), listAllTags()]).then(async (values) => {
    const subscriptions = values[0]
    const tags = values[1]

    document.querySelector("#tags-list").innerHTML = renderTags(tags, subscriptions)
    document.querySelector("#no-tags-list").innerHTML = subscriptions.filter(s => s.tag_ids.length === 0)
        .map(renderSubscription).join('')

    const addButton = document.querySelector('#add-button')
    tags.forEach(tag => addButton.addTag(tag))
})
