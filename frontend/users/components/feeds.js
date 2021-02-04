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

        async connectedCallback() {
            const feeds = await _loadAllFeeds(this._tagId)
            if (feeds.length == 0) return

            import('./feed.js')

            feeds.forEach(feed => _renderFeed(this, feed))
        }

        static get observedAttributes() {
            return ['tag_id']
        }

        attributeChangedCallback(attribute, oldValue, newValue) {
            switch (attribute) {
            case 'tag_id':
                this.tagId = newValue
                break
            }
        }

        set tagId(value) {
            this._tagId  = value
        }
    }

    const _renderFeed = async (self, feed) => {
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

    const _loadAllFeeds = async (tagId, pageSize, createdLt) => {
        const params = {}

        if (pageSize === undefined) pageSize = 100
        if (tagId !== undefined) params.tagId = tagId
        if (pageSize !== undefined) params.pageSize = pageSize
        if (createdLt !== undefined) params.createdLt = createdLt

        const tags = await FeedsService.list(params)
        
        if (tags.length < pageSize) {
            return tags
        }

        params.createdLt = tags[tags.length - 1].created

        return tags.concat(await FeedsService.list(params))
    }

    customElements.define('x-feeds', Feeds)
})()
