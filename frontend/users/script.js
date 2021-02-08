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

document.querySelector('#left').addEventListener('SubscriptionCreateSucceded', (e) => {
    const subscription = e.detail.subscription

    addSubscription(subscription)
})

document.querySelector('#left').addEventListener('SubscriptionCreateFailed', (e) => {
    const params = e.detail.params
    const error = e.detail.error

    console.log('create subscription failed', params, error)
})

document.querySelector('#left').addEventListener('TagCreateSucceded', (e) => {
    const tag = e.detail.tag

    addTag(tag)
})

document.querySelector('#left').addEventListener('TagCreateFailed', (e) => {
    const params = e.detail.params
    const error = e.detail.error

    console.log('create tag failed', params, error)
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
