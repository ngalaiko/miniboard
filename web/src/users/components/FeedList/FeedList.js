import './Feed/Feed.js'
import FeedService from '../../services/FeedService.js'

(async () => {
    const res = await fetch('/users/components/FeedList/FeedList.html')
    const textTemplate = await res.text()

    const HTMLTemplate = new DOMParser().parseFromString(textTemplate, 'text/html')
                            .querySelector('template')

    class FeedList extends HTMLElement {
        constructor() { 
            super()

            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)
        }

        connectedCallback() {
            this.render()
        }

        async render() { 
            const ulElement = this.shadowRoot.querySelector('.feed-list__list')
            const feeds = await _loadFeeds()

            feeds.forEach(feed => {
                const li = _createFeedElement(this, feed)
                ulElement.appendChild(li)
            })
        }
    }

    const _loadFeeds = async (pageToken) => {
        const {feeds, nextPageToken} = await FeedService.list(25, pageToken)

        if (nextPageToken === '') return feeds

        return feeds.concat(await _loadFeeds(nextPageToken))
    }

    const _createFeedElement = (self, feed) => {
        let userFeed = document.createElement('x-feed')
        Object.keys(feed).forEach((key) => userFeed.setAttribute(key, feed[key]))

        let li = document.createElement('li')
        li.appendChild(userFeed)

        li.onclick = () => {
            let event = new CustomEvent('FeedSelected', {
                detail: {
                    id: feed.id
                },
                bubbles: true,
                composed: true
            })
            self.dispatchEvent(event)
        }

        return li
    }

    customElements.define('feed-list', FeedList)
})()
