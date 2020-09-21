import './UserArticle/UserArticle.js'

(async () => {
    const res = await fetch('/users/components/UserArticleController/UserArticleList/UserArticleList.html')
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
        }

        get list() {
            return this._list
        }

        set list(list) {
            this._list = list
            this.render()
        }

        render() { 
            let ulElement = this.shadowRoot.querySelector('.article-list__list')
            ulElement.innerHTML = ''

            this.list.forEach(feed => {
                let li = _createArticleElement(this, feed)
                ulElement.appendChild(li)
            })
        }
    }

    const _createArticleElement = (self, article) => {
        let userArticle = document.createElement('user-article')
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

    customElements.define('user-article-list', UserArticleList)
})()
