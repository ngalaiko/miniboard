import './Reader/Reader.js'

(async () => {
    const res = await fetch('/users/components/ReaderController/ReaderController.html')
    const textTemplate = await res.text()

    const HTMLTemplate = new DOMParser().parseFromString(textTemplate, 'text/html')
                            .querySelector('template')

    class ReaderController extends HTMLElement {
        constructor() { 
             super()
        }

        connectedCallback() {
            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)

            _loadArticleIdState(this)

            window.addEventListener('popstate', () => _loadArticleIdState(this))
        }

        get articleId() {
            return this._articleId
        }

        set articleId(value) {
            this._articleId = value
            _storeArticleIdState(value)
            this.render()
        }

        async render() { 
            const userReader = this.shadowRoot.querySelector('#x-reader')

            const articleData = this.articleId ? await _fetchArticle(this.articleId) : {}

            userReader.articleData = articleData
        }
    }

    const _loadArticleIdState = (self) => {
        const urlParams = new URLSearchParams(window.location.search.slice(1))
        self.articleId = urlParams.get('article')
    }

    const _storeArticleIdState = (articleId) => {
        const urlParams = new URLSearchParams(window.location.search.slice(1))

        if (urlParams.get('article') === articleId) return 

        urlParams.set('article', articleId)

        let refresh = window.location.protocol +
            "//" + window.location.host + window.location.pathname +
            `?${urlParams.toString()}`
        window.history.pushState({ path: refresh }, '', refresh)
    }


    const _fetchArticle = async (articleId) => {
        const response = await fetch(`/api/v1/articles/${articleId}?view=ARTICLE_VIEW_FULL`)

        const body = await response.json()

        return body
    }

    customElements.define('reader-controller', ReaderController)
})()
