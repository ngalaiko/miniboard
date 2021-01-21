const apiUrl = window.location.hostname == 'localhost'
    ? 'http://localhost:80'
    : 'https://api.miniboard.app'

const createFeedElement = (feed) => {
    const list = document.querySelector('#feeds-list')

    const listItem = document.createElement('li')
    list.appendChild(listItem)

    const item = document.createElement('x-feed')
    item.setAttribute('title', feed.title)
    if (item.iconUrl !== undefined) {
        item.setAttribute('icon', feed.icon_url)
    }
    listItem.appendChild(item)
}

const loadFeeds = async () => {
    const response = await fetch(apiUrl + '/v1/feeds', {
        credentials: 'include',
    })

    const body = await response.json()

    if (response.status !== 200) {
        consle.error(`failed to fetch feeds: ${body.message}`)
        return
    }

    body.feeds.map(createFeedElement)
}

loadFeeds()
