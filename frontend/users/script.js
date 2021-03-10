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

const renderSubscription = async (subscription) => {
    if (subscription.tag_ids === null || subscription.tag_ids.length === 0) {
        document.querySelector('#subscriptions').addSubscription(subscription)
    } else {
        subscription.tag_ids
            .map((tagId) => document.getElementById(tagId))
            .forEach((xTag) => xTag.addSubscription(subscription))
    }
}

const renderTag = async (tag) => {
    document.querySelector('#add-button')
        .addTag(tag)

    await import('./components/tag.js')

    const xtag = document.querySelector('#tags-list')
        .appendChild(document.createElement('li'))
        .appendChild(document.createElement('x-tag'))

    xtag.setAttribute('id', tag.id)
    xtag.setAttribute('title', tag.title)
}

const renderItems = async (subscriptionId) => {
    const items = await ItemsService.list({
        subscriptionIdEq: subscriptionId,
    })

    items.forEach(renderItem)
}

const renderItem = (item) => {
    console.log(item)
}

const renderToastMessage = async (promise, message,onSucess) => {
    await import('./components/toast.js')
    
    const toast = document.createElement('x-toast')
    toast.setAttribute('message', message)

    if (promise) promise.then((v) => {
        if (messageOnSuccess) toast.setAttribute('message', messageOnSuccess(v))
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

    promise.then(renderSubscription)

    renderToastMessage(promise,
        `Subscribing: ${params.url}`,
        (subscription) => `Subscribed: ${subscription.title}`,
    )
})

document.querySelector('#tags-menu').addEventListener('TagCreate', (e) => {
    const params = e.detail.params
    const promise = e.detail.promise

    promise.then(addTag)

    renderToastMessage(promise,
        `Creating tag: ${params.title}`,
        (tag) => `New tag: ${tag.title}`,
    )
})

document.querySelector('#tags-menu').addEventListener('SubscriptionSelected', async (e) => {
    const subscriptionId = e.detail.id

    storeState('subscription', subscriptionId)
    renderItems(subscriptionId)
})

const subscriptionId = getState('subscription')
if (subscriptionId !== "") renderItems(subscriptionId)

Promise.all([listAllSubscriptions(), listAllTags()]).then(async (values) => {
    const subscriptions = values[0]
    const tags = values[1]

    for (const tag of tags) {
        await renderTag(tag)
    }

    for (const subscription of subscriptions) {
        await renderSubscription(subscription)
    }
})
