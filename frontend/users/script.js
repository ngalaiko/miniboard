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

    console.log('create feed succeeded', feed)
})

document.querySelector('#left').addEventListener('FeedCreateFailed', (e) => {
    const params = e.detail.params
    const error = e.detail.error

    console.log('create feed failed', params, error)
})

Promise.all([listAllFeeds(), listAllTags()]).then((values) => {
    const feeds = values[0]
    const tags = values[1]

    const xAddButton = document.createElement('x-add-button')
    xAddButton.tags = tags
    document.querySelector('#left').appendChild(xAddButton)

    const xTags = document.createElement('x-tags')
    xTags.tags = tags
    xTags.feeds = feeds.filter(feed => feed.tag_ids.length !== 0)
    document.querySelector('#left').appendChild(xTags)

    const xFeeds = document.createElement('x-feeds')
    xFeeds.feeds = feeds.filter(feed => feed.tag_ids.length === 0)
    document.querySelector('#left').appendChild(xFeeds)
})
