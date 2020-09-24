import './AddButton/AddButton.js'

(async () => {
    const res = await fetch('/users/components/AddController/AddController.html')
    const textTemplate = await res.text()

    const HTMLTemplate = new DOMParser().parseFromString(textTemplate, 'text/html')
                            .querySelector('template')

    class AddController extends HTMLElement {
        constructor() { 
             super()
        }

        connectedCallback() {
            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)

            _registerEventListeners(this)
        }
    }

    const _registerEventListeners = (self) => {
        const addButton = self.shadowRoot.querySelector('#add-button')

        addButton.addEventListener('AddClicked', (e) => {
            import('./AddModal/AddModal.js')

            let addModal = document.createElement('add-modal')
            self.shadowRoot.appendChild(addModal)

            addModal.addEventListener('Closed', (e) => {
                self.shadowRoot.removeChild(addModal)
            })

            addModal.addEventListener('FileAdded', (e) => {
                console.log(e)
            })

            addModal.addEventListener('UrlAdded', (e) => {
                console.log(e)
            })
        })
    }

    customElements.define('add-controller', AddController)
})()
