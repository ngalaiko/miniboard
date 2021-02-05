import TagsService from '../services/tags.js'

(async () => {
    const HTMLTemplate = document.createElement('template')
    HTMLTemplate.innerHTML = `
    <style>
        #tags-list{
            display: flex;
            flex-direction: column;

            list-style: none;
            margin: 0;
            padding: 0;
        }
    </style>

    <ul id="tags-list"></ul>
    `

    class Tags extends HTMLElement {
        constructor() {
            super()

            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)
        }

        async connectedCallback() {
            if (this._tags === undefined) throw 'tags not set'
            if (this._feeds === undefined) throw 'feeds not set'

            await import('./tag.js')

            const feedsByTagId = new Map()
            this._feeds.forEach(feed => {
                feed.tag_ids.forEach(tagId => {
                    const feeds = feedsByTagId.get(tagId)
                    if (feeds) {
                        feeds.push(feed)
                    } else {
                        feedsByTagId.set(tagId, [feed])
                    }
                })
            })

            this._tags.forEach(tag => {
                const feeds = feedsByTagId.has(tag.id)
                    ? feedsByTagId.get(tag.id)
                    : []

                _renderTag(this, tag, feeds)
            })
        }

        set feeds(feeds) {
            this._feeds = feeds
        }

        set tags(tags) {
            this._tags = tags
        }
    }

    const _renderTag = (self, tag, feeds) => {
        const list = self.shadowRoot.querySelector('#tags-list')

        const li = document.createElement('li')
        list.appendChild(li)

        const xTag = document.createElement('x-tag')
        xTag.feeds = feeds
        xTag.setAttribute('title', tag.title)
        li.appendChild(xTag)
    }

    customElements.define('x-tags', Tags)
})()
