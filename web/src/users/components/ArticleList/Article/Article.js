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
            return ['title', 'feedid']
        }

        attributeChangedCallback(attribute, oldValue, newValue) {
            switch (attribute) {
            case 'title':
                this.title = newValue
                break
            case 'feedid':
                this.feedid = newValue
                break
            }
        }

        set title(value) {
            if (this.shadowRoot)
                this.shadowRoot.querySelector('.article__article-title').innerHTML = value
        }

        set feedid(value) {
            FeedService.get(value).then((feed) => {
                if (feed.iconUrl !== null)  {
                    this.shadowRoot.querySelector('.article__feed-icon').src = feed.iconUrl
                } 
                this.shadowRoot.querySelector('.article__feed-title').innerHTML = feed.title
            })
        }
    }

    customElements.define('x-article', Article)
})()
