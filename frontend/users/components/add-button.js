import TagsService from '/services/tags.js'
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
                    const parser = new DOMParser()
                    const dom = parser.parseFromString(raw, "text/xml")

                    if (dom.documentElement.nodeName !== 'opml') {
                        throw new Error('only opml files supported')
                    }

                    const items = Array.from(dom.getElementsByTagName('outline'))
                        .filter(item => item.getAttribute('xmlUrl') === null)
                        .map(item => {
                            return {
                                tagTitle: item.getAttribute('title'),
                                urls: Array.from(item.getElementsByTagName('outline'))
                                    .filter(item => item.getAttribute('xmlUrl') !== null)
                                    .map(item => item.getAttribute('xmlUrl')),
                            }
                        })

                    const tagIdsByUrl = {}
                    for (const item of items) {
                        item.tag = await _createTag(self, {
                            title: item.tagTitle,
                        })

                        item.urls.forEach(url => {
                            if (!tagIdsByUrl[url]) tagIdsByUrl[url] = []
                            tagIdsByUrl[url].push(item.tag.id)
                        })
                    }

                    // one by one, ignoring errors
                    for (const url of Object.keys(tagIdsByUrl)) {
                        try {
                            await _createSubscription(self, {
                                url: url,
                                tagIds: tagIdsByUrl[url],
                            })
                        } catch (err) {
                            console.error(err)
                        }
                    }
                })
            })

            addModal.addEventListener('UrlAdded', (e) => _createSubscription(self, {
                url: e.detail.url,
            }))
        }
    }

    const _createSubscription = async (self, params) => {
        self.dispatchEvent(new CustomEvent('SubscriptionCreate', {
            detail: {
                params: params,
                promise: SubscriptionsService.create(params)
                    .then(operation => OperationsService.wait(operation.id)),
            },
            bubbles: true,
        }))
    }

    const _createTag = async (self, params) => {
        const tagPromise = TagsService.create(params)
        self.dispatchEvent(new CustomEvent('TagCreate', {
            detail: {
                params: params,
                promise: tagPromise,
            },
            bubbles: true,
        }))

        return await tagPromise
    }

    customElements.define('x-add-button', AddButton)
})()
