(async () => {
    const res = await fetch('/users/components/ArticleReader/ArticleReader.html')
    const textTemplate = await res.text()

    const HTMLTemplate = new DOMParser().parseFromString(textTemplate, 'text/html')
                            .querySelector('template')

    class ArticleReader extends HTMLElement {
        constructor() { 
            super()

            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)
        }

        static get observedAttributes() {
            return ['article']
        }

        attributeChangedCallback(attribute, oldValue, newValue) {
            switch (attribute) {
            case 'article':
                if (oldValue !== newValue) {
                    this.article = newValue
                }
            }
        }

        set article(value) {
            _displayArticle(this, value)
        }
    }

    const _displayArticle = async (self, articleId) => {
        const articleData = articleId && articleId !== null ? await _fetchArticle(articleId) : {}

        self.shadowRoot.querySelector('.article-reader__title').innerText = articleData.title ? articleData.title : ''
        self.shadowRoot.querySelector('.article-reader__link').href = articleData.url ? articleData.url : ''
        self.shadowRoot.querySelector('.article-reader__content').innerHTML = _decodeContent(articleData.content)
    }

    const _fetchArticle = async (articleId) => {
        const response = await fetch(`/api/v1/articles/${articleId}?view=ARTICLE_VIEW_FULL`)

        const body = await response.json()

        return body
    }

    const _decodeContent = (encoded) => {
        if (!encoded) return ''
        return decodeURIComponent(atob(encoded).split('').map(function(c) {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
        }).join(''))
    }

    customElements.define('article-reader', ArticleReader)
})()
