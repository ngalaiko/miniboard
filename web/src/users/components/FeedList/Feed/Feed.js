(async () => {
    const res = await fetch('/users/components/FeedList/Feed/Feed.html')
    const textTemplate = await res.text()

    const HTMLTemplate = new DOMParser().parseFromString(textTemplate, 'text/html')
                            .querySelector('template')

    class Feed extends HTMLElement {
        constructor() { 
            super()

            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)
        }

        static get observedAttributes() {
            return ['title', 'iconurl']
        }

        attributeChangedCallback(attribute, oldValue, newValue) {
            switch (attribute) {
            case 'title':
                this.title = newValue !== '' ? newValue : "Not Provided!"
                break
            case 'iconurl':
                this.iconUrl = newValue
                break
            }
        }

        set title(value) {
            if (this.shadowRoot)
                this.shadowRoot.querySelector('.feed__feed-title').innerText = value
        }

        set iconUrl(value) {
            if (!this.shadowRoot) return

            if (value !== 'null') {
                this.shadowRoot.querySelector('.feed__feed-icon').src = value
            }
        }
    }

    customElements.define('x-feed', Feed)
})()
