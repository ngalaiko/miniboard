import SubscriptionsService from '../services/subscriptions.js'

(async () => {
    const HTMLTemplate = document.createElement('template')
    HTMLTemplate.innerHTML = `
    <style>
        #subscriptions-list{
            list-style: none;
            margin: 0;
            padding: 0;
        }

        .subscription {
            cursor: pointer;
            padding: 0.1em;
        }
    </style>

    <ul id="subscriptions-list"></ul>
    `

    class Subscriptions extends HTMLElement {
        constructor() { 
            super()

            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)
        }

        async addSubscription(subscription) {
            await import('./subscription.js')
            const list = this.shadowRoot.querySelector('#subscriptions-list')

            const li = document.createElement('li')
            list.appendChild(li)

            const xSubscription = document.createElement('x-subscription')
            xSubscription.setAttribute('id', subscription.id)
            xSubscription.setAttribute('title', subscription.title)
            if (subscription.icon_url !== undefined) {
                xSubscription.setAttribute('icon', subscription.icon_url)
            }
            li.appendChild(xSubscription)
        }
    }

    customElements.define('x-subscriptions', Subscriptions)
})()
