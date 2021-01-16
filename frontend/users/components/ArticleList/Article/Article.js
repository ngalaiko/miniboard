import FeedService from '../../../services/FeedService.js'

(async () => {
    const res = await fetch('/users/components/ArticleList/Article/Article.html')
    const textTemplate = await res.text()

    const HTMLTemplate = new DOMParser().parseFromString(textTemplate, 'text/html')
                            .querySelector('template')

    class Article extends HTMLElement {
        constructor() { 
            super()

            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)
        }

        connectedCallback() {
            this.title = this.getAttribute('title')
        }


        static get observedAttributes() {
            return ['title', 'feedid', 'createtime']
        }

        attributeChangedCallback(attribute, oldValue, newValue) {
            switch (attribute) {
            case 'title':
                this.title = newValue
                break
            case 'feedid':
                this.feedid = newValue
                break
            case 'createtime':
                this.createtime = newValue
                break
            }
        }

        set createtime(value) {
            if (this.shadowRoot) {
                const date = Date.parse(value)
                const now = performance.timing.navigationStart + performance.now()
                const secondsSince = ~~((now - date) / 1000 )
                this.shadowRoot.querySelector('.article__time').innerText = _timeAgo(secondsSince)
            }
        }

        set title(value) {
            if (this.shadowRoot)
                this.shadowRoot.querySelector('.article__article-title').innerText = value
        }

        set feedid(value) {
            FeedService.get(value).then((feed) => {
                if (feed.iconUrl !== null)  {
                    this.shadowRoot.querySelector('.article__feed-icon').src = feed.iconUrl
                } 
                this.shadowRoot.querySelector('.article__feed-title').innerText = feed.title
            })
        }
    }

    const _timeAgo = (seconds) =>  {
        let minutes = ~~(seconds / 60)
        if (minutes == 0) return `${seconds}s`

        let hours = ~~(minutes / 60)
        if (hours == 0) return `${minutes}m`

        let days = ~~(hours / 24)
        if (days == 0) return `${hours}h`

        let weeks = ~~(days / 7)
        if (weeks == 0) return `${days}d`

        let years = ~~(days / 365)
        if (years == 0) return `${weeks}w`

        return `${years}y`
    }

    customElements.define('x-article', Article)
})()
