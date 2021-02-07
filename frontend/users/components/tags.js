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

        async addFeed(feed) {
            for (const tagId of feed.tag_ids) {
                const tag = this.shadowRoot.querySelector(`#tag-${tagId}`)
                await tag.addFeed(feed)
            }
        }

        async addTag(tag) {
            await import('./tag.js')

            const list = this.shadowRoot.querySelector('#tags-list')

            const li = document.createElement('li')
            list.appendChild(li)

            const xTag = document.createElement('x-tag')
            xTag.setAttribute('id', `tag-${tag.id}`)
            xTag.setAttribute('title', tag.title)
            li.appendChild(xTag)
        }
    }

    customElements.define('x-tags', Tags)
})()
