(async () => {
    const HTMLTemplate = document.createElement('template')
    HTMLTemplate.innerHTML = `
    <style>
        #feed-container {
            display: flex;
            cursor: pointer;
            align-items: center;
            padding: 0.2em;
        }

        #feed-container:hover {
            background: #cccccc;
        }

        #feed-icon {
            margin-right: 0.2em;
            min-width: 20px;
            width: 20px;
            height: 20px;
        }

        #feed-title {
            text-overflow: ellipsis;
            white-space: nowrap;
            overflow: hidden;
        }
    </style>
    <span id="feed-container">
        <img id="feed-icon" src="/img/rss.svg"></img>
        <span id="feed-title"></span>
    </span>
    `

    class Feed extends HTMLElement {
        constructor() { 
            super()

            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)
        }

        static get observedAttributes() {
            return ['title', 'icon']
        }

        attributeChangedCallback(attribute, oldValue, newValue) {
            switch (attribute) {
            case 'title':
                this.title = newValue
                break
            case 'icon':
                this.icon = newValue
                break
            }
        }

        set title(value) {
            if (!this.shadowRoot) return

            this.shadowRoot.querySelector('#feed-title').innerText = value
        }

        set icon(value) {
            if (!this.shadowRoot) return

            if (value === '') return

            this.shadowRoot.querySelector('#feed-icon').src = value
        }
    }

    customElements.define('x-feed', Feed)
})()
