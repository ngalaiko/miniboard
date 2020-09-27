(async () => {
    const res = await fetch('/users/components/ReaderController/Reader/Reader.html')
    const textTemplate = await res.text()

    const HTMLTemplate = new DOMParser().parseFromString(textTemplate, 'text/html')
                            .querySelector('template')

    class Reader extends HTMLElement {
        constructor() { 
             super()
        }

        connectedCallback() {
            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)
        }

        get articleData() {
            return this._articleData
        }

        set articleData(value) {
            this._articleData = value
            this.render()
        }

        render() { 
            this.shadowRoot.querySelector('.reader__title').innerText = this.articleData.title ? this.articleData.title : ''
            this.shadowRoot.querySelector('.reader__link').href = this.articleData.url ? this.articleData.url : ''
            this.shadowRoot.querySelector('.reader__content').innerHTML = _decodeContent(this.articleData.content)
        }
    }

    const _decodeContent = (encoded) => {
        if (!encoded) return ''
        return decodeURIComponent(atob(encoded).split('').map(function(c) {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
        }).join(''))
    }

    customElements.define('x-reader', Reader)
})()
