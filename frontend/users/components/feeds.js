import FeedsService from '../services/feeds.js'

(async () => {
    const HTMLTemplate = document.createElement('template')
    HTMLTemplate.innerHTML = `
    <style>
        #feeds-list{
            list-style: none;
            margin: 0;
            padding: 0;
        }

        .feed {
            cursor: pointer;
            padding: 0.1em;
        }
    </style>

    <ul id="feeds-list"></ul>
    `

    class Feeds extends HTMLElement {
        constructor() { 
            super()

            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)
        }

        static get observedAttributes() {
            return ['tag_id']
        }

        addFeeds(feeds) {
            feeds.forEach(feed => _renderFeed(this, feed))
        }
    }

    const _renderFeed = async (self, feed) => {
        await import('./feed.js')
        const list = self.shadowRoot.querySelector('#feeds-list')

        const li = document.createElement('li')
        list.appendChild(li)

        const xFeed = document.createElement('x-feed')
        xFeed.setAttribute('title', feed.title)
        if (feed.icon_url !== undefined) {
            xFeed.setAttribute('icon', feed.icon_url)
        }
        li.appendChild(xFeed)
    }

    customElements.define('x-feeds', Feeds)
})()
