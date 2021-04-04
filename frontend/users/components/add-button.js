import ImportsService from '/services/imports.js'
import SubscriptionsService from '/services/subscriptions.js'
import OperationsService from '/services/operations.js'

(async () => {
    const HTMLTemplate = document.createElement('template')
    HTMLTemplate.innerHTML = `
    <style>
        #add-button {
            width: 100%;
            font-size: larger;
        }
    </style>
    <button id="add-button" type="button">Add</button>
    `

    class AddButton extends HTMLElement {
        constructor() {
            super()

            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)
        }

        connectedCallback() {
            this.shadowRoot.querySelector('#add-button')
                .addEventListener('click', _clickHandler(this))
        }
    }

    const _clickHandler = (self) => {
        return () => {
            import('./add-modal.js')

            const addModal = document.createElement('x-add-modal')
            document.body.appendChild(addModal)

            addModal.addEventListener('Closed', (e) => {
                document.body.removeChild(addModal)
            })

            addModal.addEventListener('FileAdded', (e) => {
                new Promise((resolve, reject) => {
                    const reader = new FileReader()
                    reader.readAsBinaryString(e.detail.file)
                    reader.onload = e => {
                        resolve(e.target.result)
                    }
                    reader.onerror = () => reject(new Error('failed to read file'))
                })
                .then(async (raw) => {
                    self.dispatchEvent(new CustomEvent('ImportCreate', {
                        detail: {
                            promise: ImportsService.create(raw)
                                .then(operation => OperationsService.wait(operation.id)),
                        },
                        bubbles: true,
                    }))
                })
            })

            addModal.addEventListener('UrlAdded', (e) => {
                self.dispatchEvent(new CustomEvent('SubscriptionCreate', {
                    detail: {
                        params: e.detail,
                        promise: SubscriptionsService.create(e.detail)
                            .then(operation => OperationsService.wait(operation.id)),
                    },
                    bubbles: true,
                }))
            })
        }
    }

    customElements.define('x-add-button', AddButton)
})()
