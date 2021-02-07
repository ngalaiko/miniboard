import TagsService from './services/tags.js'
import FeedsService from './services/feeds.js'

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

const listAllFeeds = async (pageSize, createdLt) => {
    const params = {}

    if (pageSize === undefined) pageSize = 100
    if (pageSize !== undefined) params.pageSize = pageSize
    if (createdLt !== undefined) params.createdLt = createdLt

    const tags = await FeedsService.list(params)
    
    if (tags.length < pageSize) {
        return tags
    }

    params.createdLt = tags[tags.length - 1].created

    return tags.concat(await FeedsService.list(params))
}

document.querySelector('#left').addEventListener('FeedCreateSucceded', (e) => {
    const feed = e.detail.feed

    addFeed(feed)
})

document.querySelector('#left').addEventListener('FeedCreateFailed', (e) => {
    const params = e.detail.params
    const error = e.detail.error

    console.log('create feed failed', params, error)
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

Promise.all([listAllFeeds(), listAllTags()]).then(async (values) => {
    const feeds = values[0]
    const tags = values[1]

    for (const tag of tags) {
        await addTag(tag)
    }

    for (const feed of feeds) {
        await addFeed(feed)
    }
})

const addFeed = async (feed) => {
    if (feed.tag_ids === null || feed.tag_ids.length === 0) {
        document.querySelector('#feeds').addFeed(feed)
    } else {
        document.querySelector('#tags').addFeed(feed)
    }
}

const addTag = async (tag) => {
    const xAddButton = document.querySelector('#add-button')
    xAddButton.addTag(tag)

    const xTags = document.querySelector('#tags')
    await xTags.addTag(tag)
}
