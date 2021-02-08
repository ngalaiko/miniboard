import TagsService from './services/tags.js'
import SubscriptionsService from './services/subscriptions.js'

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

document.querySelector('#left').addEventListener('SubscriptionCreate', (e) => {
    const params = e.detail.params
    const promise = e.detail.promise

    promise.then(addSubscription)

    document.querySelector('#toasts').promise(`Subscribing: ${params.url}`, promise,
        (subscription) => `Subscribed: ${subscription.title}`,
    )
})

document.querySelector('#left').addEventListener('TagCreate', (e) => {
    const params = e.detail.params
    const promise = e.detail.promise

    promise.then(addTag)

    document.querySelector('#toasts').promise(`Creating tag: ${params.title}`, promise,
        (tag) => `New tag: ${tag.title}`,
    )
})

Promise.all([listAllSubscriptions(), listAllTags()]).then(async (values) => {
    const subscriptions = values[0]
    const tags = values[1]

    for (const tag of tags) {
        await addTag(tag)
    }

    for (const subscription of subscriptions) {
        await addSubscription(subscription)
    }
})

const addSubscription = async (subscription) => {
    if (subscription.tag_ids === null || subscription.tag_ids.length === 0) {
        document.querySelector('#subscriptions').addSubscription(subscription)
    } else {
        document.querySelector('#tags').addSubscription(subscription)
    }
}

const addTag = async (tag) => {
    const xAddButton = document.querySelector('#add-button')
    xAddButton.addTag(tag)

    const xTags = document.querySelector('#tags')
    await xTags.addTag(tag)
}
