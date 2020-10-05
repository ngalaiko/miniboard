import './Article/Article.js'
import { All as AllFeeds } from '../../services/FeedService.js'

(async () => {
    const res = await fetch('/users/components/ArticleList/ArticleList.html')
    const textTemplate = await res.text()

    const HTMLTemplate = new DOMParser().parseFromString(textTemplate, 'text/html')
                            .querySelector('template')

    class ArticleList extends HTMLElement {
        constructor() { 
            super()

            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)
        }

        connectedCallback() {
            const list = this.shadowRoot.querySelector('.article-list__list')
            list.addEventListener('scroll', e => {
                if (_needMoreElements(list)) {
                    _loadMore(this)
                }
            })
        }

        static get observedAttributes() {
            return ['feed']
        }

        attributeChangedCallback(attribute, oldValue, newValue) {
            switch (attribute) {
            case 'feed':
                if (oldValue !== newValue) {
                    this.feed = newValue
                }
            }
        }

        set feed(value) {
            this.shadowRoot.querySelector('.article-list__list').innerHTML = ''
            this.pageToken = undefined
            _loadMore(this)
        }
    }

    const _loadMore = async (self) => {
        const feedId = self.getAttribute('feed')
        const ulElement = self.shadowRoot.querySelector('.article-list__list')
        const articles = await _loadArticles(self, feedId)
        if (articles.length === 0) return

        articles.forEach(article => {
            let li = _createArticleElement(self, article)
            ulElement.appendChild(li)
        })

        if (_needMoreElements(ulElement)) {
            _loadMore(self)
        }
    }

    const _createArticleElement = (self, article) => {
        let userArticle = document.createElement('x-article')
        Object.keys(article).forEach((key) => userArticle.setAttribute(key, article[key]))

        let li = document.createElement('li')
        li.appendChild(userArticle)

        li.onclick = () => {
            let event = new CustomEvent('ArticleSelected', {
                detail: {
                    id: article.id
                },
                bubbles: true,
                composed: true
            })
            self.dispatchEvent(event)
        }

        return li
    }

    const _needMoreElements = (elem) => {
        const { scrollTop, scrollHeight, clientHeight } = elem

        const needMore = scrollTop + clientHeight >= scrollHeight - 25
        return needMore 
    }

    const _loadArticles = async (self, feedId) => {
        if (self.pageToken === '') return []
        if (self.pageToken === undefined) self.pageToken = ''

        let articlesUrl = `/api/v1/articles?page_size=10&page_token=${self.pageToken}`
        if (feedId && feedId !== AllFeeds.id) articlesUrl += `&feed_id_eq=${feedId}`

        const response = await fetch(articlesUrl)

        const body = await response.json()

        self.pageToken = body.nextPageToken

        return body.articles
    }

    customElements.define('article-list', ArticleList)
})()
