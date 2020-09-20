import './UserArticleList/UserArticleList.js'

(async () => {
    const res = await fetch('/users/components/UserArticleController/UserArticleController.html')
    const textTemplate = await res.text()

    const HTMLTemplate = new DOMParser().parseFromString(textTemplate, 'text/html')
                            .querySelector('template')

    class UserArticleList extends HTMLElement {
        constructor() { 
             super()
        }

        connectedCallback() {
            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)

            this.render()
        }

        async render() { 
            let feedListElement = this.shadowRoot.querySelector('#user-feed-list')

            let feeds = await _loadFeeds()

            feedListElement.list = feeds
        }
    }

    const _loadArticles = async () => {
        return []
    }

    customElements.define('user-article-controller', UserArticleList)
})()
