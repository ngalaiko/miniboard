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
                .then(_parseXML)
            })

            addModal.addEventListener('UrlAdded', (e) => _addUrl(e.detail.url))
        })
    }

    const _parseXML = async (raw) => {
        const parser = new DOMParser()
        const dom = parser.parseFromString(raw, "text/xml")

        if (dom.documentElement.nodeName !== 'opml') {
            throw new Error('only opml files supported')
        }

        const urls = Array.from(dom.getElementsByTagName('outline'))
            .map(item => item.getAttribute('xmlUrl'))
            .filter(item => item !== null)

        for (const url of urls) {
            try {
                await _addUrl(url)
            } catch {}
        }
    }

    const _addUrl = async (url) => {
        const response = await fetch('/api/v1/sources', {
            method: 'POST',
            body: JSON.stringify({
                url: url,
            }),
        })

        const body = await response.json()

        _watchOperation(url, body.name)
    }

    const _watchOperation = async (url, name) => {
        const response = await fetch(`/api/v1/${name}`)
        const body = await response.json()

        switch (true) {
        case !body.done:
            window.setTimeout(() => _watchOperation(url, name), 1000)
            break
        case !(body.error === undefined):
            _handleError(url, body.error)
            break
        case !(body.response === undefined):
            _handleResponse(url, body.response)
            break
        }
    }

    const _handleResponse = (url, response) => {
        switch (response['@type']) {
        case 'type.googleapis.com/miniboard.feeds.v1.Feed':
            console.log(`${url} feed added`, response.title)
            break
        case 'type.googleapis.com/miniboard.articles.v1.Article':
            console.log(`${url} article added`, response.title)
            break
        default:
            console.log(`${url}`, response)
            break
        }
    }

    const _handleError = (url, error) => {
        console.error(`${url}: ${error.message}`)
    }

    customElements.define('add-controller', AddController)
})()
