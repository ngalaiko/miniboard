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
            const tags = await _loadAllTags()
            if (tags.length == 0) return

            import('./tag.js')

            tags.forEach(tag => _renderTag(this, tag))
        }
    }

    const _renderTag = (self, tag) => {
        const list = self.shadowRoot.querySelector('#tags-list')

        const li = document.createElement('li')
        list.appendChild(li)

        const xTag = document.createElement('x-tag')
        xTag.setAttribute('id', tag.id)
        xTag.setAttribute('title', tag.title)
        li.appendChild(xTag)
    }

    const _loadAllTags = async (pageSize, createdLt) => {
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

    customElements.define('x-tags', Tags)
})()
