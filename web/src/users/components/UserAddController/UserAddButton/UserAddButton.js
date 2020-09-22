(async () => {
    const res = await fetch('/users/components/UserAddController/UserAddButton/UserAddButton.html')
    const textTemplate = await res.text()

    const HTMLTemplate = new DOMParser().parseFromString(textTemplate, 'text/html')
                            .querySelector('template')

    class UserAddButton extends HTMLElement {
        constructor() { 
             super()
        }

        connectedCallback() {
            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)
            
            this.shadowRoot.querySelector(".button__button").addEventListener('click', e => {
                let event = new CustomEvent('AddClicked', {
                    bubbles: true,
                })
                this.dispatchEvent(event)
            })
        }
    }

    customElements.define('user-add-button', UserAddButton)
})()
