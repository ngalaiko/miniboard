import './ArticleList/ArticleList.js'

(async () => {
    const res = await fetch('/users/components/ArticleController/ArticleController.html')
    const textTemplate = await res.text()

    const HTMLTemplate = new DOMParser().parseFromString(textTemplate, 'text/html')
                            .querySelector('template')

    class ArticleList extends HTMLElement {
        constructor() { 
             super()
        }

        connectedCallback() {
            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)

            _loadFeedIdState(this)

            window.addEventListener('popstate', () => _loadFeedIdState(this))
        }

        get feedId() {
            return this._feedId
        }

        set feedId(value) {
            this._feedId = value
            _storeFeedIdState(value)
            this.render()
        }

        async render() { 
            let articleListElement = this.shadowRoot.querySelector('#article-list')

            let articles = await _loadArticles(this)
            
            articleListElement.list = articles
        }
    }

    const _loadFeedIdState = (self) => {
        const urlParams = new URLSearchParams(window.location.search.slice(1))
        self.feedId = urlParams.get('feed')
    }

    const _storeFeedIdState = (feedId) => {
        const urlParams = new URLSearchParams(window.location.search.slice(1))

        if (urlParams.get('feed') === feedId) return 

        urlParams.set('feed', feedId)

        let refresh = window.location.protocol +
            "//" + window.location.host + window.location.pathname +
            `?${urlParams.toString()}`
        window.history.pushState({ path: refresh }, '', refresh)
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

    customElements.define('article-controller', ArticleList)
})()
