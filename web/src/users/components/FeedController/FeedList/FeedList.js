import './Feed/Feed.js'

(async () => {
    const res = await fetch('/users/components/FeedController/FeedList/FeedList.html')
    const textTemplate = await res.text()

    const HTMLTemplate = new DOMParser().parseFromString(textTemplate, 'text/html')
                            .querySelector('template')

    class FeedList extends HTMLElement {
        constructor() { 
             super()
        }

        connectedCallback() {
            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)
        }

        get list() {
            return this._list
        }

        set list(list) {
            this._list = list
            this.render()
        }

        render() { 
            let ulElement = this.shadowRoot.querySelector('.feed-list__list')
            ulElement.innerHTML = ''

            this.list.forEach(feed => {
                let li = _createFeedElement(this, feed)
                ulElement.appendChild(li)
            })
        }
    }

    const _createFeedElement = (self, feed) => {
        let userFeed = document.createElement('x-feed')
        Object.keys(feed).forEach((key) => userFeed.setAttribute(key, feed[key]))

        let li = document.createElement('li')
        li.appendChild(userFeed)

        li.onclick = () => {
            let event = new CustomEvent('FeedSelected', {
                detail: {
                    id: feed.id
                },
                bubbles: true,
                composed: true
            })
            self.dispatchEvent(event)
        }

        return li
    }

    customElements.define('feed-list', FeedList)
})()
