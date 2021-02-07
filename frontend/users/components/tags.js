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

            this.feedsByTagId = new Map()
        }

        addFeeds(feeds) {
            feeds.forEach(feed => {
                feed.tag_ids.forEach(tagId => {
                    const feeds = this.feedsByTagId.get(tagId)
                    if (feeds) {
                        feeds.push(feed)
                    } else {
                        this.feedsByTagId.set(tagId, [feed])
                    }
                })
            })
        }

        addTags(tags) {
            tags.forEach(tag => {
                const feeds = this.feedsByTagId.has(tag.id)
                    ? this.feedsByTagId.get(tag.id)
                    : []

                _renderTag(this, tag, feeds)
            })
        }
    }

    const _renderTag = async (self, tag, feeds) => {
        await import('./tag.js')

        const list = self.shadowRoot.querySelector('#tags-list')

        const li = document.createElement('li')
        list.appendChild(li)

        const xTag = document.createElement('x-tag')
        xTag.addFeeds(feeds)
        xTag.setAttribute('title', tag.title)
        li.appendChild(xTag)
    }

    customElements.define('x-tags', Tags)
})()
