(async () => {
    const res = await fetch('/users/components/UserArticleList/UserArticle/UserArticle.html')
    const textTemplate = await res.text()

    const HTMLTemplate = new DOMParser().parseFromString(textTemplate, 'text/html')
                            .querySelector('template')

    class UserArticle extends HTMLElement {
        constructor() { 
             super()
        }

        connectedCallback() {
            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)

            this.render()
        }

        render() { 
            this.shadowRoot.querySelector('.article__user-article-container').id = this.getAttribute('id')
            this.shadowRoot.querySelector('.article__user-article-container').innerHTML = this.getAttribute('title')
        }
    }

    customElements.define('user-article', UserArticle)
})()
