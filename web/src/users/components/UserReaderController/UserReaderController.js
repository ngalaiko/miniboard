import './UserReader/UserReader.js'

(async () => {
    const res = await fetch('/users/components/UserReaderController/UserReaderController.html')
    const textTemplate = await res.text()

    const HTMLTemplate = new DOMParser().parseFromString(textTemplate, 'text/html')
                            .querySelector('template')

    class UserReaderController extends HTMLElement {
        constructor() { 
             super()
        }

        connectedCallback() {
            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)
        }

        get articleId() {
            return this._articleId
        }

        set articleId(value) {
            this._articleId = value
            this.render()
        }

        async render() { 
            const userReader = this.shadowRoot.querySelector('#user-reader')

            const articleData = await _fetchArticle(this.articleId)

            userReader.articleData = articleData
        }
    }

    const _fetchArticle = async (articleId) => {
        const response = await fetch(`/api/v1/articles/${articleId}?view=ARTICLE_VIEW_FULL`)

        const body = await response.json()

        return body
    }

    customElements.define('user-reader-controller', UserReaderController)
})()
