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

        get feedId() {
            return this._feedId
        }

        set feedId(value) {
            this._feedId = value
            this.render()
        }

        async render() { 
            let articleListElement = this.shadowRoot.querySelector('#user-article-list')

            let articles = await _loadArticles(this)
            
            articleListElement.list = articles
        }
    }

    const _loadArticles = async (self, pageToken) => {
        if (pageToken === '') return []
        if (pageToken === undefined) pageToken = ''

        let articlesUrl = `/api/v1/articles?page_size=10&page_token=${pageToken}`
        if (self.feedId !== undefined) articlesUrl += `&feed_id_eq=${self.feedId}`

        const response = await fetch(articlesUrl)

        const body = await response.json()

        return body.articles
    }

    customElements.define('user-article-controller', UserArticleList)
})()
