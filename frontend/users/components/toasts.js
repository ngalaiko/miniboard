(async () => {
    const HTMLTemplate = document.createElement('template')
    HTMLTemplate.innerHTML = `
    <style>
        #toasts-container {
            position: fixed;

            width: 80%;

            bottom: 1em;
            left: 1em;
        }
    </style>
    <span id="toasts-container">
    </span>
    `

    class Toasts extends HTMLElement {
        constructor() { 
            super()

            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)
        }

        async promise(message, promise, messageOnSuccess) {
            await import('./toast.js')
            
            const container = this.shadowRoot.querySelector('#toasts-container')

            const toast = document.createElement('x-toast')
            toast.setAttribute('message', message)

            if (promise) promise.then((v) => {
                if (messageOnSuccess) toast.setAttribute('message', messageOnSuccess(v))
            }).catch(e => {
                toast.setAttribute('message', e)
            }).finally(() => {
                toast.close()
            })

            container.appendChild(toast)
        }
    }

    customElements.define('x-toasts', Toasts)
})()
