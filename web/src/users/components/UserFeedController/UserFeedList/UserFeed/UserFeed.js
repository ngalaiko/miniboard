(async () => {
    const res = await fetch('/users/components/UserFeedController/UserFeedList/UserFeed/UserFeed.html')
    const textTemplate = await res.text()

    const HTMLTemplate = new DOMParser().parseFromString(textTemplate, 'text/html')
                            .querySelector('template')

    class UserFeed extends HTMLElement {
        constructor() { 
             super()
        }

        connectedCallback() {
            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)

            this.title = this.getAttribute('title')
        }

        static get observedAttributes() {
            return ['title']
        }

        attributeChangedCallback(attribute, oldValue, newValue) {
            if (attribute === 'title') {
                this.title = newValue !== '' ? newValue : "Not Provided!"
            }
        }

        set title(value) {
            this._title = value
            if (this.shadowRoot)
                this.shadowRoot.querySelector('.feed__user-feed-container').innerHTML = value
        }

        get title() {
            return this._title
        }
    }

    customElements.define('user-feed', UserFeed)
})()
