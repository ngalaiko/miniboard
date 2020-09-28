(async () => {
    const res = await fetch('/users/components/ArticleList/Article/Article.html')
    const textTemplate = await res.text()

    const HTMLTemplate = new DOMParser().parseFromString(textTemplate, 'text/html')
                            .querySelector('template')

    class Article extends HTMLElement {
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
            this.shadowRoot.querySelector('.article__article-container').id = this.getAttribute('id')
            this.shadowRoot.querySelector('.article__article-container').innerHTML = this.getAttribute('title')
        }
    }

    customElements.define('x-article', Article)
})()
