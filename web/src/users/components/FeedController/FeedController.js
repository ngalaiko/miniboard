import './FeedList/FeedList.js'

(async () => {
    const res = await fetch('/users/components/FeedController/FeedController.html')
    const textTemplate = await res.text()

    const HTMLTemplate = new DOMParser().parseFromString(textTemplate, 'text/html')
                            .querySelector('template')

    class FeedController extends HTMLElement {
        constructor() { 
             super()
        }

        connectedCallback() {
            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)

            this.render()
        }

        async render() { 
            let feedListElement = this.shadowRoot.querySelector('#feed-list')

            let feeds = await _loadFeeds()

            feedListElement.list = feeds
        }
    }

    const _loadFeeds = async (pageToken) => {
        if (pageToken === '') return []
        if (pageToken === undefined) pageToken = ''

        const response = await fetch(`/api/v1/feeds?page_size=10&page_token=${pageToken}`)
        if (response.status !== 200) {
            throw `failed to fetch feeds: ${(await response.json()).message}`
        }

        const body = await response.json()

        return body.feeds.concat(await _loadFeeds(body.nextPageToken))
    }

    customElements.define('feed-controller', FeedController)
})()
