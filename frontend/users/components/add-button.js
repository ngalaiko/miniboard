import TagsService from '../services/tags.js'
import FeedsService from '../services/feeds.js'
import OperationsService from '../services/operations.js'

(async () => {
    const HTMLTemplate = document.createElement('template')
    HTMLTemplate.innerHTML = `
    <button id="add-button" type="button">Add</button>
    `

    class AddModal extends HTMLElement {
        constructor() {
            super()

            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)

            this._tagsByTitle= new Map()
        }

        connectedCallback() {
            this.shadowRoot.querySelector('#add-button')
                .addEventListener('click', _clickHandler(this, this._tagsByTitle))
        }

        set tags(tags) {
            tags.forEach(tag => this._tagsByTitle.set(tag.title, tag))
        }
    }

    const _clickHandler = (self, tagsByTitle) => {
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
                        const tag = tagsByTitle.get(item.tagTitle)
                        item.tag = tag
                            ? tag
                            : await _createTag(self, {
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
                            await _createFeed(self, {
                                url: url,
                                tagIds: tagIdsByUrl[url],
                            })
                        } catch (err) {
                            console.error(err)
                        }
                    }
                })
            })

            addModal.addEventListener('UrlAdded', (e) => _createFeed(self, {
                url: e.detail.url,
            }))
        }
    }

    const _createFeed = async (self, params) => {
        const operation = await FeedsService.create(params)

        try {
            self.dispatchEvent(new CustomEvent('FeedCreateSucceded', {
                detail: {
                    feed: await OperationsService.wait(operation.id),
                },
                bubbles: true,
            }))
        } catch (e) {
            self.dispatchEvent(new CustomEvent('FeedCreateFailed', {
                detail: {
                    params: params,
                    error: e,
                },
                bubbles: true,
            }))
        }
    }

    const _createTag = async (self, params) => {
        try {
            const tag = await TagsService.create(params)
            self.dispatchEvent(new CustomEvent('TagCreateSucceded', {
                detail: {
                    tag: tag,
                },
                bubbles: true,
            }))
            return tag
        } catch (e) {
            self.dispatchEvent(new CustomEvent('TagCreateFailed', {
                detail: {
                    params: params,
                    error: e,
                },
                bubbles: true,
            }))
            throw e
        }
    }

    customElements.define('x-add-button', AddModal)
})()
