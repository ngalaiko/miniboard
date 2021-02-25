(async () => {
    const HTMLTemplate = document.createElement('template')
    HTMLTemplate.innerHTML = `
    <style>
        #subscription-container {
            display: flex;
            cursor: pointer;
            align-items: center;
            padding: 0.2em;
        }

        #subscription-container:hover {
            background: #cccccc;
        }

        #subscription-icon {
            margin-right: 0.2em;
            min-width: 20px;
            width: 20px;
            height: 20px;
        }

        #subscription-title {
            text-overflow: ellipsis;
            white-space: nowrap;
            overflow: hidden;
        }
    </style>
    <span id="subscription-container">
        <img id="subscription-icon" src="/img/rss.svg"></img>
        <span id="subscription-title"></span>
    </span>
    `

    class Subscription extends HTMLElement {
        constructor() { 
            super()

            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)

            shadowRoot.querySelector('#subscription-container').addEventListener('click', () => {
                const event = new CustomEvent('SubscriptionSelected', {
                    detail: {
                        id: this._id,
                    },
                    bubbles: true,
                    composed: true
                })
                this.shadowRoot.dispatchEvent(event)
            })
        }

        static get observedAttributes() {
            return ['id', 'title', 'icon']
        }

        attributeChangedCallback(attribute, oldValue, newValue) {
            switch (attribute) {
            case 'id':
                this.id = newValue
                break
            case 'title':
                this.title = newValue
                break
            case 'icon':
                this.icon = newValue
                break
            }
        }

        set id(value) {
            this._id = value
        }

        set title(value) {
            if (!this.shadowRoot) return

            this.shadowRoot.querySelector('#subscription-title').innerText = value
        }

        set icon(value) {
            if (!this.shadowRoot) return

            if (value === '') return

            this.shadowRoot.querySelector('#subscription-icon').src = value
        }
    }

    customElements.define('x-subscription', Subscription)
})()
