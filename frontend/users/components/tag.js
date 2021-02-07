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
            <x-feeds id="tag-feeds"></x-feeds>
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

        async addFeeds(feeds) {
            await import('./feeds.js')
            const xFeeds = this.shadowRoot.querySelector('#tag-feeds')
            xFeeds.addFeeds(feeds)
        }
    }

    customElements.define('x-tag', Tag)
})()
