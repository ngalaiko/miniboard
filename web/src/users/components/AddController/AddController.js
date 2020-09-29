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
                new Promise((resolve, reject) => {
                    const reader = new FileReader()
                    reader.readAsBinaryString(e.detail.file)
                    reader.onload = e => {
                        resolve(e.target.result)
                    }
                    reader.onerror = () => reject(new Error('failed to read file'))
                })
                .then(btoa)
                .then(_addFile)
            })

            addModal.addEventListener('UrlAdded', (e) => _addUrl(e.detail.url))
        })
    }

    const _addUrl = async (url) => {
        const response = await fetch('/api/v1/sources', {
            method: 'POST',
            body: JSON.stringify({
                url: url,
            }),
        })

        const body = await response.json()

        _watchOperation(body.name)
    }

    const _addFile = async (raw) => {
        const response = await fetch('/api/v1/sources', {
            method: 'POST',
            body: JSON.stringify({
                raw: raw,
            }),
        })

        const body = await response.json()

        _watchOperation(body.name)
    }

    const _watchOperation = async (name) => {
        const response = await fetch(`/api/v1/${name}`)
        const body = await response.json()

        switch (true) {
        case !body.done:
            console.log(name, 'not done')
            window.setTimeout(() => _watchOperation(name), 1000)
            break
        case !(body.error === undefined):
            _handleError(body.error)
            break
        case !(body.response === undefined):
            _handleResponse(body.response)
            break
        }
    }

    const _handleResponse = (response) => {
        switch (response['@type']) {
        case 'type.googleapis.com/miniboard.users.feeds.v1.Feed':
            console.log('feed added', response.title)
            break
        case 'type.googleapis.com/miniboard.users.articles.v1.Article':
            console.log('article added', response.title)
            break
        default:
            console.log(response)
            break
        }
    }

    const _handleError = (error) => {
        console.error(error.message)
    }

    customElements.define('add-controller', AddController)
})()
