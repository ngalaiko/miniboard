(async () => {
    const HTMLTemplate = document.createElement('template')
    HTMLTemplate.innerHTML = `
    <style>
        #tag-container {
            display: flex;

            cursor: pointer;
            align-items: center;
            padding: 0.2em;
        }
        
        details{
            width: 100%;
        }

        #tag-title {
            text-overflow: ellipsis;
            white-space: nowrap;
            overflow: hidden;
        }
    </style>
    <span id="tag-container">
        <details>
            <summary id="tag-title"></summary>
        </details>
    </span>
    `

    class Tag extends HTMLElement {
        constructor() { 
            super()

            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)
        }

        static get observedAttributes() {
            return ['title', 'id']
        }

        attributeChangedCallback(attribute, oldValue, newValue) {
            switch (attribute) {
            case 'title':
                this.title = newValue
                break
            case 'id':
                this.id = newValue
                break
            }
        }

        set title(value) {
            if (!this.shadowRoot) return

            this.shadowRoot.querySelector('#tag-title').innerText = value
        }

        set id(value) {
            if (!this.shadowRoot) return

            import('./feeds.js')

            const xFeeds = document.createElement('x-feeds')
            this.shadowRoot.querySelector('details').appendChild(xFeeds)
            xFeeds.setAttribute('tag_ids', [value])
        }
    }

    customElements.define('x-tag', Tag)
})()
