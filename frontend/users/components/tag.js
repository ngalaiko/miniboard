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

        async connectedCallback() {
            if (this._feeds === undefined) throw 'feeds not set'

            await import('./feeds.js')

            const xFeeds = document.createElement('x-feeds')
            xFeeds.feeds = this._feeds
            this.shadowRoot.querySelector('details').appendChild(xFeeds)
        }

        static get observedAttributes() {
            return ['title']
        }

        attributeChangedCallback(attribute, oldValue, newValue) {
            switch (attribute) {
            case 'title':
                this.title = newValue
                break
            }
        }

        set title(value) {
            if (!this.shadowRoot) return

            this.shadowRoot.querySelector('#tag-title').innerText = value
        }

        set feeds(feeds) {
            this._feeds = feeds
        }
    }

    customElements.define('x-tag', Tag)
})()
