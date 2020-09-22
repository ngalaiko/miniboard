import './UserAddButton/UserAddButton.js'

(async () => {
    const res = await fetch('/users/components/UserAddController/UserAddController.html')
    const textTemplate = await res.text()

    const HTMLTemplate = new DOMParser().parseFromString(textTemplate, 'text/html')
                            .querySelector('template')

    class UserAddController extends HTMLElement {
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
        const addButton = self.shadowRoot.querySelector('#user-add-button')

        addButton.addEventListener('AddClicked', (e) => {
            import('./UserAddModal/UserAddModal.js')

            let addModal = document.createElement('user-add-modal')
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

    customElements.define('user-add-controller', UserAddController)
})()
